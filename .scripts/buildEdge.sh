#!/bin/bash

edgePrefix="cmd/edge"

cwd=$(pwd)

entries=$(ls $edgePrefix)
for entry in $entries; do
    scripts=(
        "cd $cwd/$edgePrefix/$entry"
        "echo '$(date +%s)' > t.txt"
        "npm run build"
        "mkdir -p $cwd/terraform/assets"
        "mv $entry.zip $cwd/terraform/assets"
        "rm -f $entry.zip"
    )

    for script in "${scripts[@]}"; do
        eval $script
    done
done