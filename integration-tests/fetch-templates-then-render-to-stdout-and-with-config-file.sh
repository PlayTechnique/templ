#!/bin/bash -elx

#This shell script runs several integration tests in a row.
#1. It verifies `templ fetch` works as expected, cloning a git repository.
#2. It validates that running `templ <an existing template>` renders that template to stdout.
#3. Finally, it verifies that running `templ <an existing template>=<a config file>` will populate the variables in the
#    template from the config file.
cd "${GITHUB_WORKSPACE}"

templ -f https://github.com/PlayTechnique/templ_templates

echo "Output the fetch charts script to stdout."
# This is a partial name match for a file in templ_templates
OUTPUT=$(templ fetch_)
echo "Validating stdout with magic string match..."
echo "${OUTPUT}" | grep "Cowardly refusing"

# Set up a test with a config file
cat > release-from-tag-config.yaml << EOF
---
OutputBinary: "roflcopter"
EOF

echo "Rendering release-from-tag.yaml with a config file"
OUTPUT=$(templ templ_templates/github_workflows/go/release-from-tag.yaml=release-from-tag-config.yaml)

echo "Validating render..."
echo "${OUTPUT}" | grep "roflcopter"
