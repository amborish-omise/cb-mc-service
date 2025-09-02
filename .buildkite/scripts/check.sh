#!/usr/bin/env bash
. .buildkite/scripts/env.sh

args=()

if [ "$BUILDKITE" = "true" ]; then
  args+=(--out-format=checkstyle)
else
  LINT_FILE=/dev/stdout
fi

golangci-lint run -c .buildkite/.golangci.yml --exclude-use-default=false "${args[@]}" > $LINT_FILE
