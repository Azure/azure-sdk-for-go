#!/usr/bin/env python

# --------------------------------------------------------------------------------------------
# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License. See License.txt in the project root for license information.
# --------------------------------------------------------------------------------------------
from typing import Dict, List, Optional
from pathlib import Path
import subprocess
from datetime import datetime
from subprocess import check_call, check_output
import argparse
import logging
import json
import re
import glob
import urllib.request


def get_latest_typespec_go_package_info():
    """Get the latest version and dependencies of @azure-tools/typespec-go from npm registry."""
    try:
        logging.info("Fetching latest @azure-tools/typespec-go info from npm registry")
        
        # Get package info from npm registry
        url = "https://registry.npmjs.org/@azure-tools/typespec-go/latest"
        with urllib.request.urlopen(url) as response:
            package_info = json.loads(response.read().decode())
        
        version = package_info.get("version")
        dev_dependencies = package_info.get("devDependencies", {})
        
        logging.info(f"Latest @azure-tools/typespec-go version: {version}")
        
        return {
            "version": version,
            "devDependencies": dev_dependencies
        }
    except Exception as e:
        logging.error(f"Failed to fetch latest typespec-go info: {e}")
        raise


def get_packed_tgz_deps(typespec_go_tgz: Path) -> dict:
    """Read the packed .tgz's package.json and return a merged map of its
    peerDependencies and dependencies.

    When pnpm packs the emitter, `workspace:` specifiers are rewritten to
    concrete published version ranges (for example `^0.69.1`). The emitter's
    generation-time dependencies live under `peerDependencies` in the tgz, so
    these ranges are the correct versions to pin the SDK's emitter-package.json
    `devDependencies` to. Reading the source package.json instead would yield
    `workspace:` specifiers, which are meaningless outside this monorepo."""
    import tarfile

    with tarfile.open(typespec_go_tgz, "r:gz") as tar:
        member = tar.extractfile("package/package.json")
        if member is None:
            raise FileNotFoundError(f"package/package.json not found in {typespec_go_tgz}")
        packed = json.load(member)

    deps = {}
    deps.update(packed.get("peerDependencies", {}))
    deps.update(packed.get("dependencies", {}))
    return deps


def find_typespec_go_tgz(typespec_go_root: str) -> Path:
    """Locate the packed @azure-tools/typespec-go .tgz in the emitter root."""
    for item in Path(typespec_go_root).iterdir():
        if "typespec-go" in item.name and item.name.endswith(".tgz"):
            return item
    logging.error("Cannot find .tgz for typespec-go")
    raise FileNotFoundError("Cannot find .tgz for typespec-go")


def update_dev_dependencies(emitter_package: dict, source_deps: dict):
    """Update devDependencies in emitter_package with versions from source_deps."""
    if "devDependencies" not in emitter_package:
        return
    for package_name in emitter_package["devDependencies"].keys():
        if package_name in source_deps:
            emitter_package["devDependencies"][package_name] = source_deps[package_name]
            logging.info(f"Updated {package_name} to version {source_deps[package_name]}")
        else:
            logging.info(f"Package {package_name} not found in dependencies, keeping existing version")


def update_emitter_package(sdk_root: str, typespec_go_root: str, use_dev_package: bool):
    # Load existing emitter-package.json
    emitter_package_path = Path(sdk_root) / "eng/emitter-package.json"
    if not emitter_package_path.exists():
        logging.error(f"emitter-package.json not found at {emitter_package_path}")
        raise FileNotFoundError(f"emitter-package.json not found at {emitter_package_path}")
    logging.info("Loading existing emitter-package.json")
    with open(emitter_package_path, "r", encoding="utf-8") as f:
        emitter_package = json.load(f)
    
    if use_dev_package:
        logging.info("Using dev package mode")

        # Find the locally packed typespec-go .tgz.
        typespec_go_tgz = find_typespec_go_tgz(typespec_go_root)

        # Read the emitter's generation-time dependencies from the packed tgz,
        # where pnpm has rewritten `workspace:` specifiers to concrete ranges.
        logging.info("Reading packed .tgz to get dependency versions")
        dev_deps = get_packed_tgz_deps(typespec_go_tgz)

        # Update devDependencies in emitter_package
        update_dev_dependencies(emitter_package, dev_deps)

        # Update emitter-package.json to use the dev package path
        emitter_package["dependencies"]["@azure-tools/typespec-go"] = typespec_go_tgz.absolute().as_posix()
    else:
        logging.info("Using released package mode")

        # Find the package.json in recent released typespec-go
        package_info = get_latest_typespec_go_package_info()

        # Update emitter-package.json to use the released package version
        if "dependencies" not in emitter_package:
            emitter_package["dependencies"] = {}
        emitter_package["dependencies"]["@azure-tools/typespec-go"] = package_info["version"]
        logging.info(f"Updated @azure-tools/typespec-go to version {package_info['version']}")

        # Update devDependencies in emitter_package
        dev_deps = package_info["devDependencies"]
        update_dev_dependencies(emitter_package, dev_deps)
    
    # Print the complete emitter_package before writing
    logging.info("Complete emitter-package.json content:")
    logging.info(json.dumps(emitter_package, indent=2))
    
    # Write the updated emitter-package.json
    with open(emitter_package_path, "w", encoding="utf-8") as f:
        json.dump(emitter_package, f, indent=2)
    
    # Update emitter-package-lock.json
    logging.info("Update emitter-package-lock.json")
    try:
        check_call(["tsp-client", "generate-lock-file"], cwd=sdk_root)
    except Exception as e:
        logging.error("Failed to update emitter-package-lock.json")
        logging.error(e)
        raise

def get_latest_commit_id() -> str:
    return (
        check_output(
            "git ls-remote https://github.com/Azure/azure-rest-api-specs.git HEAD | awk '{ print $1}'", shell=True
        )
        .decode("utf-8")
        .split("\n")[0]
        .strip()
    )


def get_typespec_go_commit_hash(typespec_go_root: str) -> str:
    """Get the current commit hash of the typespec-go repository."""
    try:
        return (
            check_output(
                "git rev-parse HEAD", shell=True, cwd=typespec_go_root
            )
            .decode("utf-8")
            .strip()
        )
    except Exception as e:
        logging.warning(f"Failed to get typespec-go commit hash: {e}")
        return "unknown"


def update_commit_id(file: Path, commit_id: str):
    with open(file, "r") as f:
        content = f.readlines()
    for idx in range(len(content)):
        if "commit:" in content[idx]:
            content[idx] = f"commit: {commit_id}\n"
            break
    with open(file, "w") as f:
        f.writelines(content)


def get_api_version_from_metadata(package_folder: Path) -> Optional[str]:
    """Extract API version from metadata.json file if it exists."""
    # Construct the metadata.json path based on the package folder structure
    # {package_folder}/testdata/_metadata.json
    metadata_path = package_folder / "testdata" / "_metadata.json"
    
    if metadata_path.exists():
        try:
            with open(metadata_path, "r") as f:
                metadata = json.load(f)
                api_version = metadata.get("apiVersion")
                if api_version:
                    logging.info(f"Found API version {api_version} in metadata.json for {package_folder.name}")
                    return api_version
        except (json.JSONDecodeError, FileNotFoundError) as e:
            logging.warning(f"Failed to read metadata.json for {package_folder.name}: {e}")
    
    return None


def get_api_version_from_client_files(package_folder: Path) -> Optional[str]:
    """Extract API version from client Go files by searching for 'Generated from API version' comment."""
    # Look for *_client.go files in the package folder
    client_files_pattern = str(package_folder / "*_client.go")
    client_files = glob.glob(client_files_pattern)
    
    api_version_pattern = re.compile(r"Generated from API version\s+([^\s,]+)")
    
    for client_file in client_files:
        try:
            with open(client_file, "r", encoding="utf-8") as f:
                content = f.read()
                match = api_version_pattern.search(content)
                if match:
                    api_version = match.group(1)
                    logging.info(f"Found API version {api_version} in {Path(client_file).name} for {package_folder.name}")
                    return api_version
        except (FileNotFoundError, UnicodeDecodeError) as e:
            logging.warning(f"Failed to read client file {client_file}: {e}")
    
    return None


def get_api_version(package_folder: Path) -> Optional[str]:
    """Get API version for a package, first trying metadata.json, then client files."""
    # First, try to get from metadata.json
    api_version = get_api_version_from_metadata(package_folder)
    
    if api_version:
        return api_version
    
    # If not found in metadata, try client files
    api_version = get_api_version_from_client_files(package_folder)
    
    if not api_version:
        logging.warning(f"Could not find API version for {package_folder.name}")
    
    return api_version

def regenerate_sdk(use_latest_spec: bool, service_filter: str, sdk_root: str, typespec_go_root: str) -> Dict[str, List[str]]:
    result = {
        "succeed_to_regenerate": [], 
        "fail_to_regenerate": [], 
        "not_found_api_version": [], 
        "time_to_regenerate": str(datetime.now()),
        "typespec_go_commit_hash": get_typespec_go_commit_hash(typespec_go_root)
    }
    # get all tsp-location.yaml
    commit_id = get_latest_commit_id()
    sdk_resourcemanager_path = Path(sdk_root) / "sdk" / "resourcemanager"
    for item in sdk_resourcemanager_path.rglob("tsp-location.yaml"):
        package_folder = item.parent
        if len(service_filter) > 0 and re.match(service_filter, package_folder.name) is None:
            continue
        logging.info(f"Regenerating {package_folder.name}...")
        if use_latest_spec:
            logging.info("Using latest spec")
            update_commit_id(item, commit_id)
        try:
            # Get API version for this package
            api_version = get_api_version(package_folder)
            
            # Build the tsp-client command with optional API version
            tsp_command = "tsp-client update"
            if api_version:
                tsp_command += f" --emitter-options api-version={api_version}"
                logging.info(f"Using API version {api_version} for {package_folder.name}")
            else:
                logging.info(f"No API version specified for {package_folder.name}, using default behavior")
                result["not_found_api_version"].append(package_folder.name)

            # Use subprocess.run for better control over output
            proc_result = subprocess.run(
                tsp_command, 
                shell=True, 
                cwd=str(package_folder),
                capture_output=True,
                text=True,
                check=True
            )
            
            # Log the output for progress tracking
            if proc_result.stdout:
                logging.info(f"Output for {package_folder.name}:")
                for line in proc_result.stdout.split('\n'):
                    if line.strip():
                        logging.info(f"  {line}")
            
            if proc_result.stderr:
                logging.warning(f"Stderr for {package_folder.name}:")
                for line in proc_result.stderr.split('\n'):
                    if line.strip():
                        logging.warning(f"  {line}")
                        
            # Check for errors in output
            output_lines = proc_result.stdout.split('\n') if proc_result.stdout else []
            errors = [line for line in output_lines if "- error " in line.lower()]
            if errors:
                raise Exception("Errors found in output:\n" + "\n".join(errors))
                
        except subprocess.CalledProcessError as e:
            logging.error(f"Failed to regenerate {package_folder.name}")
            logging.error(f"Command failed with exit code {e.returncode}")
            if e.stdout:
                logging.error(f"Stdout:\n{e.stdout}")
            if e.stderr:
                logging.error(f"Stderr:\n{e.stderr}")
            result["fail_to_regenerate"].append(package_folder.name)
        except Exception as e:
            logging.error(f"Failed to regenerate {package_folder.name}")
            logging.error(f"Error: {str(e)}")
            result["fail_to_regenerate"].append(package_folder.name)
        else:
            logging.info(f"Successfully regenerated {package_folder.name}")
            result["succeed_to_regenerate"].append(package_folder.name)
            
    result["succeed_to_regenerate"].sort()
    result["fail_to_regenerate"].sort()
    result["not_found_api_version"].sort()
    return result


def main(sdk_root: str, typespec_go_root: str, use_latest_spec: bool, service_filter: str, use_dev_package: bool):
    # Configure logging for better pipeline visibility
    logging.basicConfig(
        level=logging.INFO,
        format='%(asctime)s - %(levelname)s - %(message)s',
        handlers=[
            logging.StreamHandler(),  # Console output for pipeline
        ]
    )

    # Branch management, committing, and PR creation are handled by the GitHub
    # Actions workflow that invokes this script. This script is responsible only
    # for updating eng/emitter-package.json and regenerating the SDKs, leaving the
    # resulting changes in the working tree for the workflow to commit.
    update_emitter_package(sdk_root, typespec_go_root, use_dev_package)
    result = regenerate_sdk(use_latest_spec, service_filter, sdk_root, typespec_go_root)
    with open("regenerate-sdk-result.json", "w") as f:
        json.dump(result, f, indent=2)


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="SDK regeneration")

    parser.add_argument(
        "--sdk-root",
        help="SDK repo root folder",
        type=str,
    )

    parser.add_argument(
        "--typespec-go-root",
        help="typespec-go emitter root folder (required only for --use-dev-package)",
        type=str,
        default="",
    )

    parser.add_argument(
        "--use-latest-spec",
        help="Whether to use the latest spec",
        type=lambda x: x.lower() == 'true',
        default=False,
    )

    parser.add_argument(
        "--service-filter",
        help="An regex filter to specify which service to regenerate. If not specified, all services will be regenerated.",
        type=str,
    )

    parser.add_argument(
        "--use-dev-package",
        help="Whether to use dev package or released package",
        type=lambda x: x.lower() == 'true',
        default=False,
    )

    args = parser.parse_args()

    main(args.sdk_root, args.typespec_go_root, args.use_latest_spec, args.service_filter, args.use_dev_package)