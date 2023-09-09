#!/bin/bash -el

PROMPT=
RPROMPT=
THIS_SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
cd "${THIS_SCRIPT_DIR}/.."

go install .

tmpdir="$(mktemp -d)"
export TEMPL_DIR="${tmpdir}/templ"
mkdir -p "${TEMPL_DIR}"

# Not needed, but by cd-ing the prompt changes so it's clear we're in a subshell in the script.
cd ${TEMPL_DIR}

cat > roflcopter << EOF
---
Hello {{ .name }}
EOF

 cat > roflcopter.yaml << EOF
---
name: "World"
EOF


#echo "Setting logs to debug..."
#export TEMPL_LOG_LEVEL="debug"
export TEMPL_DEBUG_BREAK="true"

zsh

cd "${THIS_SCRIPT_DIR}"
rm -rf "${tmpdir}"
