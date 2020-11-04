#!/bin/sh
gofiles=$(git ls-tree -r HEAD --name-only | grep '\.go$')
[ -z "$gofiles" ] && exit 0

# echo $gofiles

toformat=$(gofmt -l $gofiles)
if [ -z "$toformat" ]; then
    echo "Nothing to format"
else
    echo "To be formatted:"
    echo $toformat
    exit 1
fi
