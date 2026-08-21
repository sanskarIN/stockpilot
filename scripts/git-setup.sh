#!/usr/bin/env sh
set -eu
git config user.email "sanskarin@outlook.in"
printf 'Configured repository Git email: %s\n' "$(git config user.email)"
