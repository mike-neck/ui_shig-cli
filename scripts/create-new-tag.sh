#!/usr/bin/env bash

readonly self="${0}"
readonly programDir="$(cd "${self%/*}" && pwd )"
readonly selectNewTag="${programDir}/select-new-tag.sh"
if [[ ! -x "${selectNewTag}" ]]; then
  echo "select-new-tag command not found: ${selectNewTag}" >> /dev/stderr
  exit 1
fi

readonly tagName="${1:-"$("${selectNewTag}")"}"
if [[ -z "${tagName}" ]]; then
  echo "no tag given" >> /dev/stderr
  exit 1
fi

git tag --annotate "${tagName}" \
    --message="${tagName}"

git push origin "${tagName}"
