#!/usr/bin/env bash

# 現在のバージョンを取得します。タグが存在しない場合はデフォルトで"v0.0.0"を返します。
readonly currentVersion="$(git describe --tags --abbrev 2>/dev/null || echo "v0.0.0")"

# 次のバージョン番号を算出します。
# メジャー、マイナー、またはパッチバージョンのどのバージョンを上げるかを引数で指定します。
# デフォルトではパッチバージョンが上がります。
function getNextVersion() {
  local currentMajorVersion
  # shellcheck disable=SC2001
  currentMajorVersion="$(sed 's/v//g' <<< "${currentVersion%%.*}")"
  local tmp="${currentVersion%.*}"
  local currentMinorVersion="${tmp#*.}"
  local currentPatchVersion="${currentVersion##*.}"
  local whichNextVersion="${1:-patch}"
  if [[ "${whichNextVersion}" == "major" ]]; then
    echo "v$(( currentMajorVersion + 1 )).0.0"
  elif [[ "${whichNextVersion}" == "minor" ]]; then
    echo "v${currentMajorVersion}.$(( currentMinorVersion + 1 )).0"
  elif [[ "${whichNextVersion}" == "patch" ]]; then
    echo "v${currentMajorVersion}.${currentMinorVersion}.$(( currentPatchVersion + 1 ))"
  fi
}

# メジャー、マイナー、パッチバージョンそれぞれを増やした場合の次のバージョンを表示します。
function printNextVersions() {
    local updateType
    for updateType in {major,minor,patch}; do
      echo "${updateType}: $(getNextVersion "${updateType}")"
    done
}

# メジャー、マイナー、パッチそれぞれの次のバージョンの中から一つを選択してそのバージョン番号を返します。
function selectNextVersion() {
  peco --select-1 --prompt "Next of ${currentVersion} > " < <(printNextVersions) |
    cut -d ' ' -f2 |
    tr -d '\n'
}

# 引数で与えられたメジャー、マイナー、またはパッチバージョンを上げた場合の次のバージョンを返します。
function determineNextVersionByArgument() {
    local updateType="${1:-patch}"
    getNextVersion "${updateType}"
}

# スクリプト実行時に入力された引数またはUPDATE_TYPE環境変数に基づき、
# 次のバージョン番号を決定して表示します。
readonly nextUpdate="${1:-"${UPDATE_TYPE:-""}"}"
if [[ -n "${nextUpdate}" ]]; then
  determineNextVersionByArgument "${nextUpdate}"
else
  selectNextVersion
fi
