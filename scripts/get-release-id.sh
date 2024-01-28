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

readonly apiPath="$(gh repo view --jq '"/repos/\(.owner.login)/\(.name)/releases/tags"' --json name --json owner | tr -d '"' | tr -d '\n')"
if [[ -z "${apiPath}" ]]; then
  echo "no repository found" >> /dev/stderr
  exit 2
fi

readonly releaseId="$(
gh api  \
    -H "Accept: application/vnd.github+json" \
    -H "X-GitHub-Api-Version: 2022-11-28" \
    "${apiPath}/${currentTag}" |
jq --raw-output 'select(.assets | length == 0) | .id' |
tr -d '\n'
)"
if [[ -z "${releaseId}" ]]; then
  echo "no release-id found for tag ${currentTag}" >> /dev/stderr
  exit 3
fi

echo "${releaseId}"
