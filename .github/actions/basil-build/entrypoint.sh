#!/usr/bin/env bash

set -euo pipefail

# FIXME: https://github.com/actions/checkout/issues/1169
git config --system --add safe.directory "$(pwd)"

# Changing working directory to project path
cd "$INPUT_PATH"

# Build arguments for basil command
args=""
[ "$INPUT_CROSS_COMPILE" = "true" ] && args+=" -cross-compile"
[ -n "$INPUT_PLATFORMS" ] && args+=" -platforms=${INPUT_PLATFORMS}"

# Run basil build command
output=$(eval basil project build $args)
output=$(echo $output | sed 's/ $//g')

# Set action output parameters
echo "artifacts=$output" >> $GITHUB_OUTPUT
