#!/bin/bash -elx

THIS_SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
cd "${THIS_SCRIPT_DIR}/.."

templ repo https://github.com/playtechnique/templ_templates

echo "Catting the fetch charts script..."
templ cat fetch

cat > release-from-tag-config.yaml << EOF
---
Output-Binary: "roflcopter"
EOF

echo "Rendering release-from-tag.yaml"
OUTPUT=$(templ render templ_templates/github_workflows/go/release-from-tag.yaml=release-from-tag-config.yaml)

echo "Validating render..."
echo "${OUTPUT}" | grep "roflcopter"
