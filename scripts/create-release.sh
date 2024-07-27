#!/usr/bin/env bash

set -e

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
