#!/usr/bin/env bash
set -euo pipefail

usage() {
  echo ""
  echo "  Usage: $0 -u <ghcr-username> -t <ghcr-token> -v <version>"
  echo ""
  echo "  -u, --username   GitHub username (GHCR owner)"
  echo "  -t, --token      GitHub personal access token (needs write:packages scope)"
  echo "  -v, --version    Version to release (e.g. 2.1.0)"
  echo ""
  exit 1
}

GHCR_USER=""
GHCR_TOKEN=""
VERSION=""

while [[ $# -gt 0 ]]; do
  case $1 in
    -u|--username) GHCR_USER="$2"; shift 2 ;;
    -t|--token)    GHCR_TOKEN="$2"; shift 2 ;;
    -v|--version)  VERSION="$2"; shift 2 ;;
    *) usage ;;
  esac
done

[[ -z "$GHCR_USER" || -z "$GHCR_TOKEN" || -z "$VERSION" ]] && usage

CHART_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$CHART_DIR"

echo "Updating versions to $VERSION..."
sed -i.bak "s/^version:.*/version: ${VERSION}/" Chart.yaml && rm -f Chart.yaml.bak
sed -i.bak "s/^appVersion:.*/appVersion: \"${VERSION}\"/" Chart.yaml && rm -f Chart.yaml.bak
sed -i.bak "s/tag: \"[^\"]*\"/tag: \"${VERSION}\"/" values.yaml && rm -f values.yaml.bak

echo "Logging in to ghcr.io..."
echo "$GHCR_TOKEN" | helm registry login ghcr.io -u "$GHCR_USER" --password-stdin

echo "Packaging chart..."
PACKAGE=$(helm package . | awk '{print $NF}')

echo "Pushing $PACKAGE to oci://ghcr.io/$GHCR_USER/charts..."
helm push "$PACKAGE" "oci://ghcr.io/$GHCR_USER/charts"

rm -f "$PACKAGE"
echo "Done. Remember to commit the version bump in Chart.yaml and values.yaml."
