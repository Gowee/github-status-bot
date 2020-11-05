#!/bin/sh
gofiles=$(git ls-tree -r HEAD --name-only | grep '\.go$')

togofmt=$(gofmt -l $gofiles)
togolines=$(golines -l $gofiles)
toformat=$( printf "%s\n%s" "$togofmt" "$togolines" | sort | uniq )
if [ -z "$toformat" ]; then
    echo "Nothing to format"
else
    echo "To be formatted:"
    echo "$toformat"
    exit 1
fi

# Ref: https://stackoverflow.com/questions/18535902/concatenating-two-string-variables-in-bash-appending-newline
# Ref: https://stackoverflow.com/questions/22101778/how-to-preserve-line-breaks-when-storing-command-output-to-a-variable
# Ref: https://stackoverflow.com/questions/39792766/checking-to-find-out-if-go-code-has-been-formatted/39796269