#!/usr/bin/env python

# --------------------------------------------------------------------------------------------
# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License. See License.txt in the project root for license information.
# --------------------------------------------------------------------------------------------
"""Query an Azure SDK Release Plan from Azure DevOps and print its key fields.

A "Release Plan" is a work item in the `Release` project of the
`dev.azure.com/azure-sdk` org. The auth-gated dashboard URL
`https://azsdk-releaseplan-dashboard-*.azurewebsites.net/?releaseplan=<id>`
returns HTTP 401 on plain web fetch, but the same `<id>` is the work item id and
is readable via the ADO REST API with a corp-account token.

Look up a plan by:
  --id <n>            release plan / work item id (from the dashboard URL)
  --java-package <pkg> Java package name, e.g. azure-resourcemanager-foo
  --package <pkg>     any-language package name (searches all *PackageName fields)

Usage:
  python query_release_plan.py --id 35135
  python query_release_plan.py --java-package azure-resourcemanager-containerservicepreparedimgspec
  python query_release_plan.py --package azure-mgmt-foo --json

Auth: uses `az account get-access-token` for the Azure DevOps resource. You must
be signed in with the corp account. If the default `az` login is a personal
tenant you'll get TF400813; run first:
  az account set --subscription "Azure SDK Developer Playground"
"""
import argparse
import json
import subprocess
import sys

ORG = "https://dev.azure.com/azure-sdk"
PROJECT = "Release"
ADO_RESOURCE = "499b84ac-1321-427f-aa17-267ca6975798"  # Azure DevOps

# Per-language custom fields worth printing, grouped for a compact status table.
LANGS = [".NET", "Java", "Go", "Python", "JavaScript"]
LANG_FIELD = {  # display language -> field-name suffix used by ADO
    ".NET": "Dotnet", "Java": "Java", "Go": "Go",
    "Python": "Python", "JavaScript": "JavaScript",
}
PKG_FIELDS = [
    "Custom.DotnetPackageName", "Custom.JavaPackageName", "Custom.GoPackageName",
    "Custom.PythonPackageName", "Custom.JavaScriptPackageName",
]


def get_token():
    import shutil
    az = shutil.which("az") or shutil.which("az.cmd") or "az"
    out = subprocess.run(
        [az, "account", "get-access-token", "--resource", ADO_RESOURCE,
         "--query", "accessToken", "-o", "tsv"],
        capture_output=True, text=True, shell=(sys.platform == "win32"))
    tok = out.stdout.strip()
    if not tok:
        sys.exit("Failed to get ADO token. Run: az account set --subscription "
                 "\"Azure SDK Developer Playground\" (corp account required)\n"
                 + out.stderr)
    return tok


def rest(token, url, method="GET", body=None):
    import urllib.request
    import urllib.error
    req = urllib.request.Request(url, method=method)
    req.add_header("Authorization", f"Bearer {token}")
    req.add_header("Content-Type", "application/json")
    data = json.dumps(body).encode() if body is not None else None
    try:
        with urllib.request.urlopen(req, data=data) as r:
            return json.loads(r.read().decode())
    except urllib.error.HTTPError as e:
        detail = e.read().decode(errors="replace")[:300]
        sys.exit(f"ADO REST call failed: HTTP {e.code} for {url}\n{detail}")


def find_id_by_package(token, pkg, java_only):
    fields = ["Custom.JavaPackageName"] if java_only else PKG_FIELDS
    clauses = " OR ".join(f"[{f}]='{pkg}'" for f in fields)
    wiql = (
        "SELECT [System.Id] FROM WorkItems "
        f"WHERE [System.TeamProject]='{PROJECT}' "
        "AND [System.WorkItemType]='Release Plan' "
        f"AND ({clauses}) ORDER BY [System.Id] DESC"
    )
    res = rest(token, f"{ORG}/{PROJECT}/_apis/wit/wiql?api-version=7.1",
               method="POST", body={"query": wiql})
    items = res.get("workItems", [])
    return [it["id"] for it in items]


def wiql_ids(token, where):
    wiql = (
        "SELECT [System.Id] FROM WorkItems "
        f"WHERE [System.TeamProject]='{PROJECT}' "
        "AND [System.WorkItemType]='Release Plan' "
        f"AND ({where}) ORDER BY [System.Id] DESC"
    )
    res = rest(token, f"{ORG}/{PROJECT}/_apis/wit/wiql?api-version=7.1",
               method="POST", body={"query": wiql})
    return [it["id"] for it in res.get("workItems", [])]


def resolve_release_plan_id(token, n):
    """Map a dashboard `?releaseplan=<n>` number to the ADO work item id.

    The dashboard number is the `Custom.ReleasePlanID` field, NOT the work item
    id. For newer plans the two coincide (e.g. 35135), but for older ones they
    differ (e.g. releaseplan 2101 -> work item 32864). So resolve via WIQL on
    Custom.ReleasePlanID first; only if that finds nothing, fall back to treating
    <n> as a raw work item id.
    """
    ids = wiql_ids(token, f"[Custom.ReleasePlanID]='{n}'")
    if ids:
        return ids, "release-plan-id"
    return [n], "work-item-id (no plan has ReleasePlanID={}; treating as raw id)".format(n)


def show(token, wid, as_json):
    wi = rest(token, f"{ORG}/{PROJECT}/_apis/wit/workitems/{wid}?api-version=7.1")
    f = wi["fields"]
    if as_json:
        print(json.dumps(f, indent=2))
        return
    g = f.get
    rp_id = g("Custom.ReleasePlanID", wid)
    print(f"Release Plan {rp_id} (work item {wid}): {g('System.Title', '')}")
    print(f"  dashboard: https://azsdk-releaseplan-dashboard-*"
          f".azurewebsites.net/?releaseplan={rp_id} (401 on web fetch)")
    print(f"  state: {g('System.State', '')}")
    print(f"  plan type: {g('Custom.ReleasePlanType', '')}")
    print(f"  release month: {g('Custom.SDKReleasemonth', '')}")
    print(f"  SDK type: {g('Custom.SDKtypetobereleased', '')}")
    print(f"  languages: {g('Custom.SDKLanguages', '')}")
    print(f"  spec approval: {g('Custom.APISpecApprovalStatus', '')}")
    print(f"  submitted by: {g('Custom.ReleasePlanSubmittedby', '')}")
    ns = g("Custom.NamespaceApprovalIssue")
    if ns:
        print(f"  namespace issue: {ns}")
    print("\n  language       release-status    version         PR")
    for lang in LANGS:
        s = LANG_FIELD[lang]
        status = g(f"Custom.ReleaseStatusFor{s}", "-")
        ver = g(f"Custom.ReleasedVersionFor{s}", "-")
        pr = g(f"Custom.SDKPullRequestFor{s}", "-")
        pkg = g(f"Custom.{s}PackageName", "")
        print(f"  {lang:<14} {status:<17} {ver:<15} {pr}")
        if pkg:
            print(f"    pkg: {pkg}")


def main():
    ap = argparse.ArgumentParser()
    grp = ap.add_mutually_exclusive_group(required=True)
    grp.add_argument("--id", type=int,
                     help="dashboard release-plan number (?releaseplan=<n>); "
                          "resolved via Custom.ReleasePlanID")
    grp.add_argument("--work-item", type=int,
                     help="raw ADO work item id (skip ReleasePlanID resolution)")
    grp.add_argument("--java-package", help="Java package name to search by")
    grp.add_argument("--package", help="any-language package name to search by")
    ap.add_argument("--json", action="store_true", help="dump raw fields JSON")
    args = ap.parse_args()

    token = get_token()

    if args.work_item:
        show(token, args.work_item, args.json)
        return

    if args.id:
        ids, how = resolve_release_plan_id(token, args.id)
        if how != "release-plan-id":
            print(f"note: {how}")
        if len(ids) > 1:
            print(f"Multiple plans have ReleasePlanID={args.id}: {ids}. "
                  f"Showing newest ({ids[0]}).\n")
        show(token, ids[0], args.json)
        return

    pkg = args.java_package or args.package
    ids = find_id_by_package(token, pkg, java_only=bool(args.java_package))
    if not ids:
        # fall back to a fuzzy CONTAINS match on the package fields
        fields = ["Custom.JavaPackageName"] if args.java_package else PKG_FIELDS
        where = " OR ".join(f"[{f}] CONTAINS '{pkg}'" for f in fields)
        ids = wiql_ids(token, where)
        if ids:
            print(f"No exact package match for '{pkg}'; showing closest CONTAINS "
                  f"match(es): {ids}\n")
    if not ids:
        print(f"No Release Plan found with package name '{pkg}'.")
        sys.exit(1)
    if len(ids) > 1:
        print(f"Multiple Release Plans match '{pkg}': {ids}. Showing newest "
              f"({ids[0]}); pass --work-item to pick another.\n")
    show(token, ids[0], args.json)


if __name__ == "__main__":
    main()
