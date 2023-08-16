#!/bin/bash -elx

cd "${GITHUB_WORKSPACE}"

templ repo https://github.com/playtechnique/templ_templates

echo "Catting the fetch charts script..."
templ cat fetch

cat > release-from-tag-config.yaml << EOF
---
OutputBinary: "roflcopter"
EOF

echo "Rendering release-from-tag.yaml"
OUTPUT=$(templ render templ_templates/github_workflows/go/release-from-tag.yaml=release-from-tag-config.yaml)

echo "Validating render..."
echo "${OUTPUT}" | grep "roflcopter"
