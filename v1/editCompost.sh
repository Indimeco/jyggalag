#!/bin/bash

parent_path=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )
cd "$parent_path"

date=$(date '+%Y-%m-%d')
path="../compost.md"

date_template='${CURRENT_YEAR}-${CURRENT_MONTH}-${CURRENT_DATE}'
vim $path -c "normal G\$2o" -c "normal a$date: " 