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

GOOS="${myOS}" \
  GOARCH="${myARCH}" \
  go build -o "${destinationDir}/${binaryName}" "${PWD}"/*.go

[[ -f "${destinationDir}/${binaryName}" ]] || (echo "失敗した…" > /dev/stderr && exit 3)
