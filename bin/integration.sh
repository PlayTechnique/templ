#!/bin/bash -el

THIS_SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
cd "${THIS_SCRIPT_DIR}/.."

go install .

tmpdir="$(mktemp -d)"
export TEMPL_DIR="${tmpdir}/templ"
mkdir -p "${TEMPL_DIR}"
cd "${tmpdir}"

zsh

rm -rf "${tmpdir}"
