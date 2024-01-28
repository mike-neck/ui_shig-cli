#!/usr/bin/env bash

readonly binaryName="$(basename "${PWD%-*}")"

readonly self="${0}"
readonly scriptDir="$(cd "${self%/*}" && pwd)"
readonly rootDir="${scriptDir%/*}"

readonly osName="${1:-"$(go env GOOS)"}"
readonly archName="${2:-"$(go env GOARCH)"}"
readonly versionName="$(git describe --tags --abbrev 2>/dev/null || echo "v0.0.0")"
if [[ -z "${versionName}" ]]; then
  echo "version is not set" >> /dev/stderr
  exit 1
fi

readonly binaryDir="${rootDir}/bin/${osName}/${archName}"
if [[ ! -d "${binaryDir}" ]]; then
  echo "binary directory does not exist: ${binaryDir}" >> /dev/stderr
  exit 0
fi

declare -a zipParams=()

readonly archiveName="${binaryName}-${osName}-${archName}-${versionName}.zip"
if [[ -f "${binaryDir}/${archiveName}" ]]; then
  rm "${binaryDir}/${archiveName}"
  echo "removed old archive file: ${binaryDir}/${archiveName}"
fi
zipParams+=("${archiveName}")

if [[ ! -f "${binaryDir}/${binaryName}" || ! -x "${binaryDir}/${binaryName}" ]]; then
  echo "binary file does not exist or is not executable: ${binaryDir}/${binaryName}" >> /dev/stderr
  exit 2
fi

if [[ "${osName}" == "windows" ]]; then
  binaryName="${binaryName}.exe"
fi

zipParams+=("${binaryName}")

zipParams+=("README.md")

cd "${binaryDir}" && zip "${zipParams[@]}"

if [[ -f "${binaryDir}/${archiveName}" ]]; then
  echo "archive file created: ${binaryDir}/${archiveName}"
  exit 0
else
  echo "failed to create archive file: ${binaryDir}/${archiveName}" >> /dev/stderr
  exit 4
fi
