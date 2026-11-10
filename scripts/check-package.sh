#!/usr/bin/env bash
# Checks that every macOS and Linux build in build/ ships the Agari engine
# of its own platform where the app looks for it: a package without one
# builds and opens, and only fails when a hand is scored.
set -euo pipefail

cd "$(dirname "$0")/.."
checked=0
fail() {
  echo "check-package: $*" >&2
  exit 1
}

for app in build/darwin-*/*.app; do
  [ -d "$app" ] || continue
  engine="$app/Contents/Resources/agari/agari"
  [ -x "$engine" ] || fail "$app has no executable agari/agari among its resources"
  case "$app" in
    build/darwin-universal/*) want="arm64 x86_64" ;;
    build/darwin-arm64/*) want="arm64" ;;
    *) want="x86_64" ;;
  esac
  for arch in $want; do
    lipo "$engine" -verify_arch "$arch" || fail "$engine is not built for $arch"
  done
  echo "$app: Agari for $want"
  checked=$((checked + 1))
done

for dir in build/linux-*; do
  [ -d "$dir" ] || continue
  engine="$dir/agari/agari"
  [ -x "$engine" ] || fail "$dir has no executable agari/agari"
  case "$dir" in
    *-arm64) want="aarch64" ;;
    *) want="x86-64" ;;
  esac
  # Each listing is taken whole before it is searched: grep -q stops at its
  # match, and under pipefail the lister's broken pipe would then fail the
  # very check that passed.
  kind=$(file -b "$engine")
  grep -q "ELF 64-bit.*$want" <<< "$kind" || fail "$engine is not a Linux $want executable: $kind"
  for archive in "$dir"/*.tar.gz; do
    listing=$(tar -tzf "$archive")
    grep -q '\(^\|/\)agari/agari$' <<< "$listing" || fail "$archive does not hold agari/agari"
  done
  for deb in "$dir"/*.deb; do
    listing=$(dpkg-deb -c "$deb")
    grep -q '/agari/agari$' <<< "$listing" || fail "$deb does not hold agari/agari"
  done
  echo "$dir: Agari for $want, in its archive and package"
  checked=$((checked + 1))
done

[ "$checked" -gt 0 ] || fail "no macOS or Linux build in build/"
