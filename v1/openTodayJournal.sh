#!/bin/bash

parent_path=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )
cd "$parent_path"

date=$(date '+%Y-%m-%d')
new_journal="../journal/${date}.md"
cp -vn '../.foam/templates/journal.md' $new_journal

date_template='${CURRENT_YEAR}-${CURRENT_MONTH}-${CURRENT_DATE}'
cursor_template='$0'
vim $new_journal -c "%s/${date_template}/${date}/ge | %s/${cursor_template}//ge"
