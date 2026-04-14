import json
import requests
import datetime
import hashlib
import os
import tempfile
import zipfile
from typing import List, Dict, Optional

# Configuration
GITHUB_RELEASES_URL = "https://api.github.com/repos/CDFER/JLCPCB-Kicad-Library/releases"
DEFAULT_STATUS = "stable"
DEFAULT_KICAD_VERSION = "8.0"


def get_latest_releases(releases_url: str) -> List[Dict]:
    """Fetch latest releases from GitHub repository."""
    headers = {"Accept": "application/vnd.github.v3+json", "User-Agent": "CDFER-Kicad-Repo-Updater"}

    try:
        response = requests.get(releases_url, headers=headers)
        response.raise_for_status()
        return response.json()
    except requests.exceptions.RequestException as e:
        raise RuntimeError(f"Failed to fetch releases: {e}")


def calculate_sha256(file_path: str) -> str:
    """Calculate SHA256 hash of a file."""
    sha256 = hashlib.sha256()
    with open(file_path, "rb") as f:
        while chunk := f.read(4096):
            sha256.update(chunk)
    return sha256.hexdigest()


def get_zip_contents_size(zip_path: str) -> int:
    """Calculate total size of all files in a ZIP archive."""
    total_size = 0
    with zipfile.ZipFile(zip_path, "r") as z:
        for file_info in z.infolist():
            total_size += file_info.file_size
    return total_size


def extract_metadata(zip_path: str, version: str) -> tuple:
    """Extract metadata from zip file's metadata.json."""
    try:
        with zipfile.ZipFile(zip_path, "r") as z:
            with z.open("metadata.json") as meta_file:
                metadata = json.load(meta_file)

                for ver in metadata.get("versions", []):
                    if ver.get("version") == version:
                        return (ver.get("status", DEFAULT_STATUS), ver.get("kicad_version", DEFAULT_KICAD_VERSION))
    except (KeyError, json.JSONDecodeError, zipfile.BadZipFile):
        pass

    return DEFAULT_STATUS, DEFAULT_KICAD_VERSION


def process_asset(asset_url: str, version: str) -> Optional[tuple]:
    """Download and process asset to get required metadata."""
    tmp_path = None
    try:
        with tempfile.NamedTemporaryFile(delete=False) as tmp_file:
            response = requests.get(asset_url, stream=True)
            response.raise_for_status()

            for chunk in response.iter_content(chunk_size=8192):
                tmp_file.write(chunk)
            tmp_path = tmp_file.name

        sha256 = calculate_sha256(tmp_path)
        download_size = os.path.getsize(tmp_path)
        install_size = get_zip_contents_size(tmp_path)
        status, kicad_version = extract_metadata(tmp_path, version)

        return sha256, download_size, install_size, status, kicad_version
    except Exception as e:
        print(f"Failed to process asset {asset_url}: {e}")
        return None
    finally:
        if tmp_path and os.path.exists(tmp_path):
            os.remove(tmp_path)


def validate_packages_structure(packages: Dict) -> None:
    """Validate the structure of the packages.json data."""
    required_keys = ["packages"]
    if not all(key in packages for key in required_keys):
        raise KeyError(f"Missing required keys in packages.json: {required_keys}")


def update_packages_json(packages_json_file: str, releases: List[Dict]) -> None:
    """Update packages.json with new release information."""
    try:
        with open(packages_json_file, "r") as f:
            packages = json.load(f)
    except (FileNotFoundError, json.JSONDecodeError) as e:
        raise RuntimeError(f"Error loading {packages_json_file}: {e}")

    validate_packages_structure(packages)

    existing_versions = {v["version"]: v for package in packages["packages"] for v in package.get("versions", [])}

    for release in releases:
        version = release.get("tag_name")
        assets = release.get("assets", [])

        asset = next(
            (
                a
                for a in assets
                if a.get("name", "").startswith("JLCPCB-KiCad-Library-") and a.get("name", "").endswith(".zip")
            ),
            None,
        )

        if not asset or not version or version in existing_versions:
            continue

        download_url = asset.get("browser_download_url")
        result = process_asset(download_url, version)

        if not result:
            continue

        sha256, d_size, i_size, status, kicad_ver = result

        new_version = {
            "version": version,
            "status": status,
            "kicad_version": kicad_ver,
            "download_sha256": sha256,
            "download_size": d_size,
            "install_size": i_size,
            "download_url": download_url,
        }

        packages["packages"][0]["versions"].insert(0, new_version)

    try:
        with open(packages_json_file, "w") as f:
            json.dump(packages, f, indent=4)
    except IOError as e:
        raise RuntimeError(f"Error writing to {packages_json_file}: {e}")


def update_repository_json(repository_json_file: str) -> None:
    """Update repository.json with current timestamp."""
    try:
        with open(repository_json_file, "r") as f:
            repository = json.load(f)
    except (FileNotFoundError, json.JSONDecodeError) as e:
        raise RuntimeError(f"Error loading {repository_json_file}: {e}")

    current_time = datetime.datetime.now(datetime.timezone.utc)
    repository["packages"]["update_time_utc"] = current_time.strftime("%Y-%m-%d %H:%M:%S")
    repository["packages"]["update_timestamp"] = int(current_time.timestamp())

    try:
        with open(repository_json_file, "w") as f:
            json.dump(repository, f, indent=4)
    except IOError as e:
        raise RuntimeError(f"Error writing to {repository_json_file}: {e}")


def main() -> None:
    """Main function to update repository metadata."""
    try:
        releases = get_latest_releases(GITHUB_RELEASES_URL)
        update_packages_json("packages.json", releases)
        update_repository_json("repository.json")
        print("Successfully updated repository metadata")
    except Exception as e:
        print(f"Error occurred: {e}")
        raise


if __name__ == "__main__":
    main()
