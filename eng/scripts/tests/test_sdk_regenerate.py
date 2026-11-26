#!/usr/bin/env python

# --------------------------------------------------------------------------------------------
# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License. See License.txt in the project root for license information.
# --------------------------------------------------------------------------------------------

import io
import json
import sys
import tarfile
import tempfile
import unittest
from pathlib import Path
from unittest.mock import patch

SCRIPTS_DIR = Path(__file__).resolve().parent.parent
if str(SCRIPTS_DIR) not in sys.path:
    sys.path.insert(0, str(SCRIPTS_DIR))

import sdk_regenerate  # noqa: E402


def _create_tgz(path: Path, package_json: dict) -> None:
    payload = json.dumps(package_json).encode("utf-8")
    with tarfile.open(path, "w:gz") as tar:
        info = tarfile.TarInfo("package/package.json")
        info.size = len(payload)
        tar.addfile(info, io.BytesIO(payload))


class TestSdkRegenerate(unittest.TestCase):
    def test_update_emitter_package_dev_mode(self):
        with tempfile.TemporaryDirectory() as tmpdir:
            root = Path(tmpdir)
            sdk_root = root / "azure-sdk-for-go"
            (sdk_root / "eng").mkdir(parents=True)

            emitter_package_path = sdk_root / "eng" / "emitter-package.json"
            emitter_package_path.write_text(
                json.dumps(
                    {
                        "main": "dist/src/index.js",
                        "dependencies": {
                            "@azure-tools/typespec-go": "0.18.0",
                        },
                        "devDependencies": {
                            "@azure-tools/typespec-client-generator-core": "0.73.0",
                            "@azure-tools/typespec-azure-portal-core": "0.73.0",
                            "@azure-tools/typespec-liftr-base": "0.13.0",
                            "@typespec/compiler": "1.17.0",
                        },
                    }
                ),
                encoding="utf-8",
            )

            packages_root = root / "typespec-azure" / "packages"
            typespec_go_root = packages_root / "typespec-go"
            typespec_go_root.mkdir(parents=True)
            _create_tgz(
                typespec_go_root / "azure-tools-typespec-go-0.19.0.tgz",
                {
                    "name": "@azure-tools/typespec-go",
                    "version": "0.19.0",
                    "peerDependencies": {
                        "@azure-tools/typespec-client-generator-core": "^0.74.0",
                        "@typespec/compiler": "^1.18.0",
                    },
                },
            )

            portal_dir = packages_root / "typespec-azure-portal-core"
            portal_dir.mkdir()
            (portal_dir / "package.json").write_text(
                json.dumps(
                    {
                        "name": "@azure-tools/typespec-azure-portal-core",
                        "version": "0.74.0-dev.1",
                    }
                ),
                encoding="utf-8",
            )

            calls = []
            with (
                patch("sdk_regenerate.check_call", side_effect=lambda cmd, cwd=None: calls.append((cmd, cwd))),
                patch("sdk_regenerate.get_registry_latest_version", return_value="0.14.0"),
            ):
                sdk_regenerate.update_emitter_package(
                    str(sdk_root), str(typespec_go_root), use_dev_package=True
                )

            updated = json.loads(emitter_package_path.read_text(encoding="utf-8"))
            self.assertEqual(
                updated["dependencies"]["@azure-tools/typespec-go"],
                (typespec_go_root / "azure-tools-typespec-go-0.19.0.tgz").absolute().as_posix(),
            )
            self.assertEqual(updated["devDependencies"]["@azure-tools/typespec-client-generator-core"], "0.74.0")
            self.assertEqual(updated["devDependencies"]["@azure-tools/typespec-azure-portal-core"], "0.74.0-dev.1")
            self.assertEqual(updated["devDependencies"]["@azure-tools/typespec-liftr-base"], "0.14.0")
            self.assertEqual(updated["devDependencies"]["@typespec/compiler"], "1.18.0")
            self.assertNotIn("overrides", updated)
            self.assertEqual(calls, [(["tsp-client", "generate-lock-file"], str(sdk_root))])

    def test_update_emitter_package_released_mode(self):
        with tempfile.TemporaryDirectory() as tmpdir:
            sdk_root = Path(tmpdir) / "azure-sdk-for-go"
            (sdk_root / "eng").mkdir(parents=True)
            emitter_package_path = sdk_root / "eng" / "emitter-package.json"
            emitter_package_path.write_text(
                json.dumps(
                    {
                        "main": "dist/src/index.js",
                        "dependencies": {"@azure-tools/typespec-go": "0.18.0"},
                        "devDependencies": {"@typespec/compiler": "1.17.0"},
                    }
                ),
                encoding="utf-8",
            )

            calls = []
            with (
                patch(
                    "sdk_regenerate.get_latest_typespec_go_package_info",
                    return_value={
                        "version": "0.19.0",
                        "devDependencies": {"@typespec/compiler": "1.18.0"},
                    },
                ),
                patch("sdk_regenerate.check_call", side_effect=lambda cmd, cwd=None: calls.append((cmd, cwd))),
            ):
                sdk_regenerate.update_emitter_package(str(sdk_root), "", use_dev_package=False)

            updated = json.loads(emitter_package_path.read_text(encoding="utf-8"))
            self.assertNotIn("overrides", updated)
            self.assertEqual(updated["dependencies"]["@azure-tools/typespec-go"], "0.19.0")
            self.assertEqual(updated["devDependencies"]["@typespec/compiler"], "1.18.0")
            self.assertEqual(calls, [(["tsp-client", "generate-lock-file"], str(sdk_root))])

    def test_update_emitter_package_next_version_runs_bump_deps_before_lock_file(self):
        with tempfile.TemporaryDirectory() as tmpdir:
            sdk_root = Path(tmpdir) / "azure-sdk-for-go"
            (sdk_root / "eng").mkdir(parents=True)
            emitter_package_path = sdk_root / "eng" / "emitter-package.json"
            emitter_package_path.write_text(
                json.dumps(
                    {
                        "main": "dist/src/index.js",
                        "dependencies": {"@azure-tools/typespec-go": "0.18.0"},
                        "devDependencies": {"@typespec/compiler": "1.17.0"},
                    }
                ),
                encoding="utf-8",
            )

            calls = []
            with (
                patch(
                    "sdk_regenerate.get_latest_typespec_go_package_info",
                    return_value={
                        "version": "0.19.0",
                        "devDependencies": {"@typespec/compiler": "1.18.0"},
                    },
                ),
                patch("sdk_regenerate.shutil.which", return_value="npx"),
                patch("sdk_regenerate.check_call", side_effect=lambda cmd, cwd=None: calls.append((cmd, cwd))),
            ):
                sdk_regenerate.update_emitter_package(
                    str(sdk_root), "", use_dev_package=False, use_next_version=True
                )

            self.assertEqual(
                calls,
                [
                    (
                        [
                            "npx",
                            "--yes",
                            "@azure-tools/typespec-bump-deps",
                            "eng/emitter-package.json",
                            "--add-npm-overrides",
                        ],
                        str(sdk_root),
                    ),
                    (["tsp-client", "generate-lock-file"], str(sdk_root)),
                ],
            )


if __name__ == "__main__":
    unittest.main()
