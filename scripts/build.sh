#!/usr/bin/env bash

readonly myOS="${1:-"$(go env GOOS)"}"
if [[ -z "${myOS}" ]]; then
  echo "GOOS が存在しないが？" > /dev/stderr
  exit 1
fi

readonly myARCH="${2:-"$(go env GOARCH)"}"
if [[ -z "${myARCH}" ]]; then
  echo "GOARCH が存在しないが？" > /dev/stderr
  exit 1
fi

readonly destinationDir="${3:-"${PWD}/bin"}"
if [[ -z "${destinationDir}" ]]; then
  echo "どこに出力すればいいのかわからないが？" > /dev/stderr
  exit 1
fi

readonly binaryName="${4:-ui_shig}"
if [[ -z "${binaryName}" ]]; then
  echo "バイナリーの名前がないが？" > /dev/stderr
  exit 1
fi

readonly commitHash="$(git rev-parse HEAD 2>/dev/null || echo "000000")"
if [[ -z "${commitHash}" ]]; then
  echo "コミットがないが？" > /dev/stderr
  exit 2
fi

readonly version="$(git describe --tags --abbrev 2>/dev/null || echo "v0.0.0")"
if [[ -z "${version}" ]]; then
  echo "バージョンがないが？" > /dev/stderr
  exit 2
fi

readonly currentDateTime="$(date '+%Y-%m-%dT%H:%M:%S%Z')"
if [[ -z "${currentDateTime}" ]]; then
  echo "日付がないが？" > /dev/stderr
  exit 2
fi

GOOS="${myOS}" \
  GOARCH="${myARCH}" \
  go build \
      -ldflags "-X main.UiShigVersion=${version} -X main.UiShigCommit=${commitHash} -X main.UiShigBuildDate=${currentDateTime}" \
      -o "${destinationDir}/${myOS}/${myARCH}/${binaryName}" "${PWD}"/*.go

[[ -f "${destinationDir}/${myOS}/${myARCH}/${binaryName}" ]] || (echo "失敗した…" > /dev/stderr && exit 3)
