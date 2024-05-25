#!/bin/bash

set -euxo pipefail

parent_path=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )
cd "$parent_path"
cd $(git rev-parse --show-toplevel)


git add .
if ! git diff-index --quiet HEAD; then
    git commit -m"quick sync from $1"
fi
git pull -r
git push
