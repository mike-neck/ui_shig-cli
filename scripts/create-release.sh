#!/usr/bin/env bash

set -e

readonly currentBranch="$(git rev-parse --abbrev-ref HEAD)"
if [[ -z "${currentBranch}" ]]; then
  echo "repository is broken because a current branch is not available." >> /dev/stderr
  exit 1
elif [[ ! "${currentBranch}" == "main" ]]; then
  echo "current branch is not [main], but ${currentBranch}" >> /dev/stderr
  exit 1
fi

readonly currentTag="$(git describe --tags --abbrev=0)"
if [[ -z "${currentTag}" ]]; then
  echo "no tag given" >> /dev/stderr
  exit 2
fi

readonly apiPath="$(\
gh repo view \
   --jq '"repos/\(.owner.login)/\(.name)/releases"' \
   --json name --json owner |
tr -d '"' |
tr -d '\n'
)"
if [[ -z "${apiPath}" ]]; then
  echo "no repository found" >> /dev/stderr
  exit 2
fi

gh release create \
  "${currentTag}" \
  --generate-notes \
  --title "Release of ${currentTag}" \
  --verify-tag
