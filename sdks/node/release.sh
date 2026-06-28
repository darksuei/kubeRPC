#!/usr/bin/env bash
set -euo pipefail

usage() {
  echo ""
  echo "  Usage: $0 -t <npm-token> -v <version> [--otp <code>]"
  echo ""
  echo "  -t, --token   npm access token (or set NPM_TOKEN env var)"
  echo "  -v, --version Version to release (e.g. 2.1.0)"
  echo "  --otp         One-time password (required if 2FA is enabled)"
  echo ""
  exit 1
}

TOKEN="${NPM_TOKEN:-}"
VERSION=""
OTP=""

while [[ $# -gt 0 ]]; do
  case $1 in
    -t|--token)   TOKEN="$2"; shift 2 ;;
    -v|--version) VERSION="$2"; shift 2 ;;
    --otp)        OTP="$2"; shift 2 ;;
    *) usage ;;
  esac
done

[[ -z "$TOKEN" || -z "$VERSION" ]] && usage

SDK_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$SDK_DIR"

echo "Updating package.json to $VERSION..."
npm version "$VERSION" --no-git-tag-version

echo "Building SDK..."
yarn build

echo "Authenticating with npm..."
npm config set //registry.npmjs.org/:_authToken "$TOKEN"

PUBLISH_ARGS="--access public"
[[ -n "$OTP" ]] && PUBLISH_ARGS="$PUBLISH_ARGS --otp=$OTP"

echo "Publishing kuberpc-node@$VERSION..."
npm publish $PUBLISH_ARGS

echo "Done. Remember to commit the version bump in package.json."
