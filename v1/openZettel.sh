#!/bin/bash

parent_path=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )
cd "$parent_path"

name="$@"

id=$(ls -l ../zettelkasten | grep ^- | wc -l)
zettel="../zettelkasten/[${id}]${name}.md"
cp -vn '../.foam/templates/default.md' "$zettel"

date=$(date '+%Y-%m-%d')
date_template='${CURRENT_YEAR}-${CURRENT_MONTH}-${CURRENT_DATE}'
cursor_template='$0'
type_template='${1|fleeting,literature,permanent|}'
vim "$zettel" -c "%s/${date_template}/${date}/ge | %s/${type_template}/${type}/ge | %s/${cursor_template}//ge"
