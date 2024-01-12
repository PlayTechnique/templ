#!/bin/bash -elx

cd "${GITHUB_WORKSPACE}"

templ -fetch https://github.com/PlayTechnique/templ_templates

echo "Output the fetch charts script to stdout."
# This is a partial name match for a file in templ_templates
templ fetch_

cat > release-from-tag-config.yaml << EOF
---
OutputBinary: "roflcopter"
EOF

echo "Rendering release-from-tag.yaml"
OUTPUT=$(templ templ_templates/github_workflows/go/release-from-tag.yaml=release-from-tag-config.yaml)

echo "Validating render..."
echo "${OUTPUT}" | grep "roflcopter"
