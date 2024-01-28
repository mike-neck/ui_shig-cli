#!/usr/bin/env bash

readonly self="${0}"
readonly scriptDir="$(cd "${self%/*}" && pwd)"
readonly rootDir="${scriptDir%/*}"
readonly baseName="${rootDir##*/}"

readonly releaseId="${RELEASE_ID:-""}"
if [[ -z "${releaseId}" ]]; then
  echo "release-id is not available" >> /dev/stderr
  exit 2
fi

readonly osName="${1:-"$(go env GOOS)"}"
readonly archName="${2:-"$(go env GOARCH)"}"
readonly versionName="$(git describe --tags --abbrev 2>/dev/null || echo "v0.0.0")"
readonly fileName="${baseName}-${osName}-${archName}-${versionName}.zip"

readonly archiveFilePath="${rootDir}/bin/${osName}/${archName}/${fileName}"
if [[ -f "${archiveFilePath}" ]]; then
  echo "release asset file does not exist: ${archiveFilePath}" >> /dev/stderr
  # arm64 のビルド諦めた
  exit 0
fi

readonly apiPath="$(gh repo view --jq '"/repos/\(.owner.login)/\(.name)/releases"' --json name --json owner | tr -d '"' | tr -d '\n')"
if [[ -z "${apiPath}" ]]; then
  echo "no repository found" >> /dev/stderr
  exit 2
fi

gh api \
  --method POST \
  -H "Accept: application/vnd.github+json" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  "/repos/OWNER/REPO/releases/${releaseId}/assets?name=${fileName}" \
  -f "@${archiveFilePath}"
