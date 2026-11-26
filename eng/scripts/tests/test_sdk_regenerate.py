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
    def test_update_emitter_package_dev_mode_packs_workspace_deps(self):
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

            tcgc_dir = packages_root / "typespec-client-generator-core"
            tcgc_dir.mkdir()
            (tcgc_dir / "package.json").write_text(
                json.dumps(
                    {
                        "name": "@azure-tools/typespec-client-generator-core",
                        "version": "0.74.0-dev.1",
                        "scripts": {"build": "tsc"},
                    }
                ),
                encoding="utf-8",
            )
            tcgc_tgz = tcgc_dir / "azure-tools-typespec-client-generator-core-0.74.0-dev.1.tgz"
            tcgc_tgz.write_bytes(b"tgz")

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
            portal_tgz = portal_dir / "azure-tools-typespec-azure-portal-core-0.74.0-dev.1.tgz"
            portal_tgz.write_bytes(b"tgz")

            build_calls = []

            def fake_check_call(cmd, cwd=None, **kwargs):
                if Path(cmd[0]).name.startswith("pnpm") and cmd[1:] == ["build"]:
                    build_calls.append(Path(cwd))

            def fake_check_output(cmd, cwd=None, text=None, **kwargs):
                if Path(cwd) == tcgc_dir:
                    return f"{tcgc_tgz.name}\n"
                if Path(cwd) == portal_dir:
                    return f"{portal_tgz.name}\n"
                raise AssertionError(f"Unexpected check_output cwd: {cwd}")

            with (
                patch("sdk_regenerate.shutil.which", return_value="pnpm"),
                patch("sdk_regenerate.check_call", side_effect=fake_check_call),
                patch("sdk_regenerate.check_output", side_effect=fake_check_output),
                patch("sdk_regenerate.get_registry_latest_version", return_value="0.14.0"),
            ):
                sdk_regenerate.update_emitter_package(
                    str(sdk_root), str(typespec_go_root), use_dev_package=True
                )

            updated = json.loads(emitter_package_path.read_text(encoding="utf-8"))
            self.assertEqual(build_calls, [tcgc_dir])
            self.assertEqual(
                updated["dependencies"]["@azure-tools/typespec-go"],
                (typespec_go_root / "azure-tools-typespec-go-0.19.0.tgz").absolute().as_posix(),
            )
            self.assertEqual(
                updated["devDependencies"]["@azure-tools/typespec-client-generator-core"],
                tcgc_tgz.absolute().as_posix(),
            )
            self.assertEqual(
                updated["devDependencies"]["@azure-tools/typespec-azure-portal-core"],
                portal_tgz.absolute().as_posix(),
            )
            self.assertEqual(updated["devDependencies"]["@azure-tools/typespec-liftr-base"], "0.14.0")
            self.assertEqual(updated["devDependencies"]["@typespec/compiler"], "1.18.0")
            self.assertEqual(
                updated["overrides"],
                {
                    "@azure-tools/typespec-client-generator-core": tcgc_tgz.absolute().as_posix(),
                    "@azure-tools/typespec-azure-portal-core": portal_tgz.absolute().as_posix(),
                },
            )

    def test_update_emitter_package_released_mode_skips_overrides(self):
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

            with (
                patch(
                    "sdk_regenerate.get_latest_typespec_go_package_info",
                    return_value={
                        "version": "0.19.0",
                        "devDependencies": {"@typespec/compiler": "1.18.0"},
                    },
                ),
                patch("sdk_regenerate.check_call"),
            ):
                sdk_regenerate.update_emitter_package(str(sdk_root), "", use_dev_package=False)

            updated = json.loads(emitter_package_path.read_text(encoding="utf-8"))
            self.assertNotIn("overrides", updated)
            self.assertEqual(updated["dependencies"]["@azure-tools/typespec-go"], "0.19.0")
            self.assertEqual(updated["devDependencies"]["@typespec/compiler"], "1.18.0")


if __name__ == "__main__":
    unittest.main()
