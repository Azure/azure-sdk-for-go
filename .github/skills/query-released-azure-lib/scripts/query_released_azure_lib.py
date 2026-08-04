#!/usr/bin/env python

# --------------------------------------------------------------------------------------------
# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License. See License.txt in the project root for license information.
# --------------------------------------------------------------------------------------------
"""Check which Azure SDK languages have released a given management (or client)
library, by reading the per-language release inventory CSVs published at
`https://raw.githubusercontent.com/Azure/azure-sdk/main/_data/releases/latest/`.

These are the same CSVs that drive https://azure.github.io/azure-sdk/releases/latest
(whose `?search=` box is client-side JS, useless to a static fetch). A row carries
`VersionGA` and/or `VersionPreview`; once a language publishes a package it appears
here with a version, so absence / empty-version = not released.

The hard part is that every language names the same service differently:
  .NET    Azure.ResourceManager.<Svc>[.<Sub>] (PascalCase, dot-separated; 3-4 segs)
  Java    azure-resourcemanager-<svc>        (lower, dash)
  Python  azure-mgmt-<svc>                   (lower, dash, no dots)
  JS      @azure/arm-<svc>                   (lower, dash)
  Go      sdk/resourcemanager/<mid>/arm<suffix> (lower path)
So we reduce each Package to a normalized service token (strip the language prefix,
drop non-alphanumerics, lowercase) and compare tokens.

Usage:
  # by shared service token (normalized match across all languages)
  python query_released_azure_lib.py --service commvaultcontentstore

  # give exact per-language package names when the token differs by language
  # (e.g. from the query-release-plan skill's Custom.<Lang>PackageName fields)
  python query_released_azure_lib.py --dotnet Azure.ResourceManager.Commvault \
      --java azure-resourcemanager-commvault --python azure-mgmt-commvault \
      --js @azure/arm-commvault --go armcommvault

  # emphasise .NET only
  python query_released_azure_lib.py --service networkcloud --lang dotnet
"""
import argparse
import csv
import io
import re
import sys
import urllib.request

BASE = "https://raw.githubusercontent.com/Azure/azure-sdk/main/_data/releases/latest"
LANGS = ["dotnet", "java", "python", "js", "go"]
# Language package prefixes stripped (case-insensitive) before normalizing.
# Go is handled separately (its path is split), so no go/arm prefix here.
PREFIXES = [
    "azure.resourcemanager.",
    "azure-resourcemanager-",
    "azure-mgmt-",
    "@azure/arm-",
]


def norm(s):
    return re.sub(r"[^a-z0-9]", "", (s or "").lower())


def strip_prefixes(s):
    s = (s or "").strip().lower()
    for p in PREFIXES:
        if s.startswith(p):
            return s[len(p):]
    return s


def go_tokens(package):
    """Candidate tokens for a Go package `sdk/resourcemanager/<mid>/arm<suffix>`.

    Go may namespace a sub-service as `compute/armbulkactions` (mid=compute,
    suffix=bulkactions) or flatten it as `computebulkactions/armcomputebulkactions`.
    So emit mid, suffix, and (when they differ) mid+suffix, to match either the
    parent token or the combined service token used by other languages.
    """
    toks, parts = set(), [x for x in (package or "").split("/") if x]
    if "resourcemanager" in parts:
        i = parts.index("resourcemanager")
        mid = parts[i + 1] if i + 1 < len(parts) else ""
        last = parts[i + 2] if i + 2 < len(parts) else ""
        suf = last[3:] if last.lower().startswith("arm") else last
        if mid:
            toks.add(norm(mid))
        if suf:
            toks.add(norm(suf))
        if mid and suf and norm(mid) != norm(suf):
            toks.add(norm(mid + suf))
    if not toks:
        toks.add(norm(package))
    return {t for t in toks if t}


def service_token(package, lang):
    """Normalized service token(s) from a language Package string."""
    if lang == "go":
        return go_tokens(package)
    return {norm(strip_prefixes(package))}


def row_tokens(row, lang):
    """All normalized identity tokens for a row: package-derived + ServiceName."""
    toks = set(service_token(row.get("Package"), lang))
    sn = norm(row.get("ServiceName"))
    if sn and sn != "unknown":
        toks.add(sn)
    return {t for t in toks if t}


def fetch_csv(lang):
    with urllib.request.urlopen(f"{BASE}/{lang}-packages.csv") as r:
        text = r.read().decode("utf-8-sig")
    return list(csv.DictReader(io.StringIO(text)))


def status_of(row):
    ga = (row.get("VersionGA") or "").strip()
    prev = (row.get("VersionPreview") or "").strip()
    if ga and prev:
        return f"GA {ga} (+ preview {prev})"
    if ga:
        return f"GA {ga}"
    if prev:
        return f"preview {prev}"
    return "listed, no version"


def find_rows(rows, lang, query_token):
    """Token/fuzzy match for the --service path. Returns (rows, kind)."""
    exact = [r for r in rows if query_token in row_tokens(r, lang)]
    if exact:
        return exact, "token"
    # Fuzzy: only when the query is a proper substring of a *package* token
    # (a more specific variant), never on ServiceName and never the reverse
    # direction (which lets short names like "store" match "...contentstore").
    def pkg_superset(r):
        if not query_token or len(query_token) < 4:
            return False
        return any(query_token in t and len(t) > len(query_token)
                   for t in service_token(r.get("Package"), lang))

    fuzzy = [r for r in rows if pkg_superset(r)]
    return (fuzzy, "fuzzy") if fuzzy else ([], None)


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--service", help="shared service token, e.g. commvaultcontentstore")
    for lang in LANGS:
        ap.add_argument(f"--{lang}", help=f"exact {lang} package name to match")
    ap.add_argument("--lang", choices=LANGS, action="append",
                    help="limit output to these languages (repeatable)")
    args = ap.parse_args()

    overrides = {lang: getattr(args, lang) for lang in LANGS}
    if not args.service and not any(overrides.values()):
        ap.error("provide --service and/or per-language package name(s)")

    query_token = norm(strip_prefixes(args.service)) if args.service else None
    langs = args.lang or LANGS

    label = args.service or next(v for v in overrides.values() if v)
    print(f"Release status for '{label}' (source: {BASE})\n")
    print(f"  {'language':<9} {'status':<26} {'GA date':<12} {'preview date':<12} package")
    released_langs, missing_langs = [], []
    for lang in langs:
        try:
            rows = fetch_csv(lang)
        except Exception as e:  # noqa: BLE001 - report and continue per language
            print(f"  {lang:<9} ERROR fetching CSV: {e}")
            continue
        exact_pkg = overrides.get(lang)
        if exact_pkg:
            want = norm(exact_pkg)
            matched = [r for r in rows if norm(r.get("Package")) == want]
            kind = "exact-package" if matched else None
        else:
            matched, kind = find_rows(rows, lang, query_token)
        if not matched:
            print(f"  {lang:<9} {'NOT FOUND':<26} {'':<12} {'':<12} -")
            missing_langs.append(lang)
            continue
        lang_released = False
        for row in matched:
            st = status_of(row)
            gad = (row.get("LatestGADate") or row.get("FirstGADate") or "").strip()
            prd = (row.get("FirstPreviewDate") or "").strip()
            pkg = row.get("Package")
            flag = "" if kind in ("exact-package", "token") else f" [{kind} match]"
            print(f"  {lang:<9} {st:<26} {gad:<12} {prd:<12} {pkg}{flag}")
            if st != "listed, no version":
                lang_released = True
        (released_langs if lang_released else missing_langs).append(lang)

    print()
    print(f"  released:     {', '.join(released_langs) or '(none)'}")
    print(f"  not released: {', '.join(missing_langs) or '(none)'}")
    dotnet_state = "released" if "dotnet" in released_langs else (
        "NOT released" if "dotnet" in langs else "not checked")
    print(f"  .NET: {dotnet_state}")


if __name__ == "__main__":
    sys.exit(main())
