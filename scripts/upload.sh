#!/usr/bin/env bash

readonly self="${0}"
readonly scriptDir="$(cd "${self%/*}" && pwd)"
readonly rootDir="${scriptDir%/*}"
readonly baseName="$(basename "${PWD%-*}")"

readonly osName="${1:-"$(go env GOOS)"}"
readonly archName="${2:-"$(go env GOARCH)"}"
readonly versionName="$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")"
readonly fileName="${baseName}-${osName}-${archName}-${versionName}.zip"

readonly archiveFilePath="${rootDir}/bin/${osName}/${archName}/${fileName}"
if [[ ! -f "${archiveFilePath}" ]]; then
  echo "SKIP release asset file does not exist: ${archiveFilePath}" >> /dev/stderr
  # arm64 のビルド諦めた
  exit 0
fi

readonly apiPath="$(gh repo view --jq '"\(.owner.login)/\(.name)"' --json name --json owner | tr -d '"' | tr -d '\n')"
if [[ -z "${apiPath}" ]]; then
  echo "no repository found" >> /dev/stderr
  exit 2
fi

echo "repo: ${apiPath} tag: ${versionName}"

gh release upload \
  --repo "${apiPath}" \
  "${versionName}" \
  "${archiveFilePath}"
