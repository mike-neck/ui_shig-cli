#!/usr/bin/env bash

set -e
set -o pipefail

readonly shigreUiButton="http://cbtm.html.xdomain.jp/usbtn/usbtn.html"

readonly DQ='"'
readonly FS="${IFS}"
readonly self="${0}"
readonly programDir="$(cd "${self%/*}" && pwd )"
readonly rootDir="${programDir%/*}"
readonly binDir="${rootDir}/bin"
readonly cacheDir="${binDir}/caches"
readonly dataDir="${rootDir}/data"
readonly dataFile="${dataDir}/ui-shig.jsonl"

declare -a replacePattern=()
replacePattern+=('s/a:/"a":/g')
replacePattern+=('s/k:/"k":/g')
replacePattern+=('s/label:/"label":/g')
replacePattern+=('s/videoId:/"videoId":/g')
replacePattern+=('s/time:/"time":/g')
replacePattern+=('s/new:false[,]*//g')
replacePattern+=('s/,[[:space:]]*/,/g')

readonly htmlFileName="usbtn_$(date "+%Y%m").txt"
readonly htmlFile="${cacheDir}/${htmlFileName}"

if [[ ! -d "${cacheDir}" ]]; then
  mkdir -p "${cacheDir}"
  echo "init.sh: created cache directory ${cacheDir}"
fi
if [[ ! -f "${htmlFile}" ]]; then
  curl --request GET \
      --silent --location \
      --url "${shigreUiButton}" |
    grep '"id"' |
    grep '"src"' |
    grep 'a:' |
    grep 'k:' > "${htmlFile}"
  echo "init.sh: download しぐれういボタン ${htmlFile}"
fi

if [[ ! -f "${htmlFile}" ]]; then
  echo "しぐれういボタンのダウンロードに失敗しました。" > /dev/stderr
  echo "多分、次のどれかが原因だと思う。" > /dev/stderr
  echo "  - ディスクがいっぱいになっている" > /dev/stderr
  echo "  - しぐれういボタンのURLが変わった" > /dev/stderr
  echo "  - しぐれういボタンが落ちてる" > /dev/stderr
  echo "  - ネットワークが繋がってない" > /dev/stderr
  echo "  - このスクリプトがバグってる" > /dev/stderr
  exit 3
fi

if [[ ! -d "${dataDir}" ]]; then
  mkdir -p "${dataDir}"
  echo "init.sh: created data directory ${dataDir}"
fi
if [[ -f "${dataFile}" ]]; then
  rm -rf "${dataFile}"
  echo "init.sh: removed old file ${dataFile}"
fi

declare line text pattern
while read -r line; do
  text="{${line#*\{}"
  for pattern in "${replacePattern[@]}"; do
    text="$(sed "${pattern}" <<< "${text}")"
  done
  if [[ -n "${text}" ]]; then
    echo "${text}" >> "${dataFile}"
  fi
done < "${htmlFile}"

echo "init.sh: created file ${dataFile}"
