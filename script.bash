#!/bin/sh
git filter-branch --env-filter '
CORRECT_NAME="badc0re"
CORRECT_EMAIL="dame.jovanoski@gmail.com"
export GIT_AUTHOR_NAME="$CORRECT_NAME"
export GIT_AUTHOR_EMAIL="$CORRECT_EMAIL"
' --tag-name-filter cat -- --branches --tags
