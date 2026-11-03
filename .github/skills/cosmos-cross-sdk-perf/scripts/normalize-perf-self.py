#!/usr/bin/env python3
# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License.

"""Normalize rounded perf Self rows to an exact 100 percent composition."""

from __future__ import annotations

import argparse
import csv
import pathlib


def normalize(path: pathlib.Path) -> None:
    with path.open(newline="", encoding="utf-8") as handle:
        rows = list(csv.DictReader(handle))
    if not rows:
        raise ValueError("perf Self CSV has no rows")
    values = [float(row["self_percent"]) for row in rows]
    total = sum(values)
    if total <= 0:
        raise ValueError("perf Self total must be positive")
    normalized = [value / total * 100.0 for value in values]
    normalized[-1] += 100.0 - sum(normalized)
    with path.open("w", newline="", encoding="utf-8") as handle:
        writer = csv.DictWriter(handle, fieldnames=("symbol", "dso", "self_percent"))
        writer.writeheader()
        for row, value in zip(rows, normalized):
            writer.writerow(
                {
                    "symbol": row["symbol"],
                    "dso": row["dso"],
                    "self_percent": f"{value:.12f}",
                }
            )


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("path", type=pathlib.Path)
    arguments = parser.parse_args()
    normalize(arguments.path)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
