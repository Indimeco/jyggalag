#!/bin/bash

# addPoemNumber :: (string, number) -> string
function addPoemNumber () {
    path=$1
    iter=$2

    if [[ $iter -eq 0 ]]
    then
        path=$path
    else
        path="$path"
        path="$(echo $path | sed -e s/\.md/\-$iter\.md/)"
    fi

    if test -f "$path"; then
        next_iter=$((1 + $iter))
        addPoemNumber $path $next_iter
    else
        echo $path
    fi
}


parent_path=$( cd "$(dirname "${BASH_SOURCE[0]}")" ; pwd -P )
cd "$parent_path"

year=$(date '+%Y')
date=$(date '+%Y-%m-%d')
path="../composition/${year}/${date}.md"
new_poem="$(addPoemNumber $path 0)"

cp -vn '../.foam/templates/composition.md' $new_poem

date_template='${CURRENT_YEAR}-${CURRENT_MONTH}-${CURRENT_DATE}'
cursor_template='$0'
vim $new_poem -c "%s/${date_template}/${date}/ge | %s/${cursor_template}//ge"
