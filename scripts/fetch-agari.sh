#!/usr/bin/env bash
# Fetches the Agari scoring engine for macOS and Linux into the platform
# directories of resources/, where mygo build and go test find it:
#
#   scripts/fetch-agari.sh                      # the machine's own platform
#   scripts/fetch-agari.sh darwin-arm64 darwin-amd64
#   scripts/fetch-agari.sh linux-amd64 linux-arm64
#
# The archives come from Agari's release and are checked against the
# SHA-256 sums pinned here before anything is unpacked. The Windows engine
# is committed in resources/windows-amd64. To move to another Agari
# version, change VERSION and the sums (sha256.sum of the release), and
# replace the Windows executable.
set -euo pipefail

VERSION=0.27.0
BASE="https://github.com/agari-industries/agari/releases/download/v$VERSION"

target() {
  case "$1" in
    darwin-arm64) echo "aarch64-apple-darwin e0002b22e753806676a05e327544e8a87ecd8d5a52436c6b048105a1c168bd28" ;;
    darwin-amd64) echo "x86_64-apple-darwin ac9680ea74aed97d13834f46d98c343000cb9ddd87c758b4d1d85082246ff2af" ;;
    linux-arm64) echo "aarch64-unknown-linux-gnu 053d7bcc26322bf5662f891136e5f1a80f4d50495c210bb323ff2c9ad17a2a6c" ;;
    linux-amd64) echo "x86_64-unknown-linux-gnu d7392f127b0c0cbf526e8e48b014e07d49cc56e28aa3f530d5e29423d308b9b5" ;;
    *) echo "fetch-agari: no Agari build for $1 (darwin-arm64, darwin-amd64, linux-arm64, linux-amd64)" >&2; return 1 ;;
  esac
}

sha256() {
  if command -v sha256sum > /dev/null; then sha256sum "$1" | cut -d' ' -f1; else shasum -a 256 "$1" | cut -d' ' -f1; fi
}

if [ $# -eq 0 ]; then
  os=$(uname -s | tr '[:upper:]' '[:lower:]')
  case "$(uname -m)" in
    arm64 | aarch64) arch=arm64 ;;
    *) arch=amd64 ;;
  esac
  set -- "$os-$arch"
fi

root=$(cd "$(dirname "$0")/.." && pwd)
work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT

for platform in "$@"; do
  spec=$(target "$platform")
  read -r triple want <<< "$spec"
  out="$root/resources/$platform/agari/agari"
  archive="$work/agari-$triple.tar.xz"
  echo "Agari $VERSION for $platform"
  curl --proto '=https' --tlsv1.2 -fsSL --retry 3 -o "$archive" "$BASE/agari-$triple.tar.xz"
  got=$(sha256 "$archive")
  if [ "$got" != "$want" ]; then
    echo "fetch-agari: $archive has SHA-256 $got, expected $want" >&2
    exit 1
  fi
  mkdir -p "$work/$platform" "$(dirname "$out")"
  tar -xJf "$archive" -C "$work/$platform"
  engine=$(find "$work/$platform" -type f -name agari | head -n 1)
  if [ -z "$engine" ]; then
    echo "fetch-agari: no agari executable in $archive" >&2
    exit 1
  fi
  install -m 0755 "$engine" "$out"
done
