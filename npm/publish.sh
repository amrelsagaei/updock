#!/usr/bin/env bash
# Copyright (c) 2026 Amr Elsagaei. All rights reserved.
#
# Publishes the updock npm packages for a release.
# Usage: ./npm/publish.sh <version>   (version without the leading 'v')
#
# Expects GoReleaser to have produced ./dist with per-platform binaries.
# Publishes one @updock/<os>-<arch> package per platform, then the main
# `updock` package that depends on them via optionalDependencies.
set -euo pipefail

VERSION="${1:?usage: publish.sh <version>}"
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DIST="${ROOT}/dist"
WORK="$(mktemp -d)"
trap 'rm -rf "${WORK}"' EXIT

# platform key -> "npm_os npm_cpu goreleaser_dir binary_name"
declare -A PLATFORMS=(
  ["linux-x64"]="linux x64 updock_linux_amd64_v1 updock"
  ["linux-arm64"]="linux arm64 updock_linux_arm64_v8.0 updock"
  ["darwin-x64"]="darwin x64 updock_darwin_amd64_v1 updock"
  ["darwin-arm64"]="darwin arm64 updock_darwin_arm64_v8.0 updock"
  ["win32-x64"]="win32 x64 updock_windows_amd64_v1 updock.exe"
)

publish_platform() {
  local key="$1"
  read -r npm_os npm_cpu gr_dir bin_name <<<"${PLATFORMS[$key]}"

  local pkgdir="${WORK}/${key}"
  mkdir -p "${pkgdir}/bin"

  local src
  src="$(find "${DIST}" -type f -path "*${gr_dir}*/${bin_name}" | head -1)"
  if [[ -z "${src}" ]]; then
    echo "error: binary for ${key} not found under ${DIST}/${gr_dir}" >&2
    return 1
  fi
  cp "${src}" "${pkgdir}/bin/${bin_name}"
  chmod +x "${pkgdir}/bin/${bin_name}"

  cat >"${pkgdir}/package.json" <<JSON
{
  "name": "@updock/${key}",
  "version": "${VERSION}",
  "description": "updock binary for ${npm_os} ${npm_cpu}",
  "license": "MIT",
  "os": ["${npm_os}"],
  "cpu": ["${npm_cpu}"],
  "files": ["bin/"]
}
JSON

  echo "Publishing @updock/${key}@${VERSION}"
  # --provenance ties the package to this source build via GitHub OIDC.
  (cd "${pkgdir}" && npm publish --access public --provenance)
}

# 1. Publish all platform packages first so the main package's deps resolve.
for key in "${!PLATFORMS[@]}"; do
  publish_platform "${key}"
done

# 2. Publish the main package, pinning the platform dep versions to this release.
maindir="${WORK}/main"
mkdir -p "${maindir}/bin"
cp "${ROOT}/npm/bin/updock.js" "${maindir}/bin/updock.js"

node -e "
const pkg = require('${ROOT}/npm/package.json');
pkg.version = '${VERSION}';
for (const dep of Object.keys(pkg.optionalDependencies)) {
  pkg.optionalDependencies[dep] = '${VERSION}';
}
require('fs').writeFileSync('${maindir}/package.json', JSON.stringify(pkg, null, 2));
"

echo "Publishing updock@${VERSION}"
(cd "${maindir}" && npm publish --access public --provenance)

echo "All npm packages published for ${VERSION}."
