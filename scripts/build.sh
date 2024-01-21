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

function useCgo() {
  if [[ "${myOS}" == "linux" ]]; then
    echo "1"
  else
    echo "0"
  fi
}

# CGO_ENABLED=1 が必要
# GOARCH を設定すると CGO_ENABLED=0 に設定されてしまうので明示的に指定する
GOOS="${myOS}" \
  GOARCH="${myARCH}" \
  CGO_ENABLED="$(useCgo)" \
  go build \
      -ldflags "-X main.UiShigVersion=${version} -X main.UiShigCommit=${commitHash} -X main.UiShigBuildDate=${currentDateTime}" \
      -o "${destinationDir}/${myOS}/${myARCH}/${binaryName}" "${PWD}"/*.go

if [[ -f "${destinationDir}/${myOS}/${myARCH}/${binaryName}" ]] ; then
  echo "build success ${destinationDir}/${myOS}/${myARCH}/${binaryName}"
  exit 0
else
  echo "失敗した…" > /dev/stderr
  exit 3
fi
