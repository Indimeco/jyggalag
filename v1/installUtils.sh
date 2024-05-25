#!/bin/bash

DIR=$(realpath -s "${BASH_SOURCE%/*}")
if [[ ! -d "$DIR" ]]; then DIR="$PWD"; fi
alias_file="$HOME/.notes_aliases.indimeco"

cat <<- EOF > $alias_file
alias n='cd $DIR && cd .. && if command -v code &> /dev/null; then code .; fi'
alias ns='$DIR/quickNotesSync.sh $(hostname)'
alias np='$DIR/openNewPoem.sh'
alias nj='$DIR/openTodayJournal.sh'
alias nz='$DIR/openZettel.sh'
alias nc='$DIR/editCompost.sh'
EOF

shell_profile="$HOME/.bashrc"
source_aliases=". $alias_file"
if [[ -f ~/.zshrc ]]; then     
    shell_profile="$HOME/.zshrc"
    source_aliases="source $alias_file";
fi

if [[ $(grep -c $alias_file $shell_profile) -eq 0 ]]; then 
    echo $source_aliases >> $shell_profile; 
fi 