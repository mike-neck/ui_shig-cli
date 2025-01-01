#!/usr/bin/env bash

set -e
set -o pipefail

readonly shigreUiButton="https://leiros.cloudfree.jp/usbtn/usbtn.html"

readonly DQ='"'
readonly FS="${IFS}"
readonly self="${0}"
readonly programDir="$(cd "${self%/*}" && pwd )"
readonly rootDir="${programDir%/*}"
readonly binDir="${rootDir}/bin"
readonly cacheDir="${binDir}/caches"
readonly dataDir="${rootDir}/data"
readonly dataFile="${dataDir}/ui-shig.jsonl"

# しぐれういボタンからデータをダウンロードする

declare -a replacePattern=()
replacePattern+=('s/a:/"a":/g')
replacePattern+=('s/k:/"k":/g')
replacePattern+=('s/label:/"label":/g')
replacePattern+=('s/videoId:/"videoId":/g')
replacePattern+=('s/time:/"time":/g')
replacePattern+=('s/"new":false[,]*//g')
replacePattern+=('s/"new":true[,]*//g')
replacePattern+=('s/,[[:space:]]*/,/g')
replacePattern+=('s/[[:space:]]*$//g')
replacePattern+=('s/,}/}/g')

readonly htmlFileName="usbtn_$(date "+%Y%m").txt"
readonly htmlFile="${cacheDir}/${htmlFileName}"

if [[ ! -d "${cacheDir}" ]]; then
  mkdir -p "${cacheDir}"
  echo "init.sh[SUCCESS]: created cache directory ${cacheDir}"
fi
if [[ ! -f "${htmlFile}" ]]; then
  curl --request GET \
      --silent --location \
      --url "${shigreUiButton}" |
    grep '"id"' |
    grep '"src"' |
    grep -v '//,' > "${htmlFile}"
  echo "init.sh[SUCCESS]: downloaded しぐれういボタン ${htmlFile}"
fi

if [[ ! -f "${htmlFile}" ]]; then
  echo "しぐれういボタンのダウンロードに失敗しました。" > /dev/stderr
  echo "多分、次のどれかが原因だと思う。" > /dev/stderr
  echo "  - ディスクがいっぱいになっている" > /dev/stderr
  echo "  - しぐれういボタンのURLが変わった" > /dev/stderr
  echo "  - しぐれういボタンが落ちてる" > /dev/stderr
  echo "  - ネットワークが繋がってない" > /dev/stderr
  echo "  - このスクリプトがバグってる" > /dev/stderr
  echo "しぐれういボタンが変更された可能性もあるので、報告してください" > /dev/stderr
  echo "" > /dev/stderr
  exit 3
fi

if [[ ! -d "${dataDir}" ]]; then
  mkdir -p "${dataDir}"
  echo "init.sh[SUCCESS]: created data directory ${dataDir}"
fi
if [[ -f "${dataFile}" ]]; then
  rm -rf "${dataFile}"
  echo "init.sh[SUCCESS]: removed old file ${dataFile}"
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

echo "init.sh[SUCCESS]: created file ${dataFile}"

# イラストレーターのしぐれういと申します

readonly illustratorPath="$(jq --raw-output 'select(.id == "illust") | .src | gsub("\\./";"")' "${dataFile}")"
if [[ -z "${illustratorPath}" ]]; then
  echo "「イラストレーターのしぐれういと申します」データが存在しないです" > /dev/stderr
  echo "しぐれういボタンが変更された可能性もあるので、報告してください" > /dev/stderr
  echo "" > /dev/stderr
  exit 3
fi

readonly illustratorFileName="${illustratorPath##*/}"
readonly illustratorFile="${dataDir}/${illustratorFileName}"

if [[ -f "${illustratorFile}" ]]; then
  rm "${illustratorFile}"
  echo "init.sh[SUCCESS]: removed old file ${illustratorFile}"
fi

curl --request GET \
     --silent --location \
     --url "${shigreUiButton%/*}/${illustratorPath}"\
     --output "${illustratorFile}"

if [[ -f "${illustratorFile}" ]]; then
  echo "init.sh[SUCCESS]: downloaded 「イラストレーターのしぐれういと申します」 voice ${illustratorFile}"
else
  echo "「イラストレーターのしぐれういと申します」のダウンロードに失敗しました" > /dev/stderr
  echo "しぐれういボタンが変更された可能性もあるので、報告してください" > /dev/stderr
  echo "" > /dev/stderr
  exit 3
fi
