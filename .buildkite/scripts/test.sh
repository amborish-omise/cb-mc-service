#!/usr/bin/env bash
. .buildkite/scripts/env.sh

if [ "${BUILDKITE}" = "true" ]; then
  gotestsum \
  --jsonfile ${TEST_FILE} \
  --junitfile ${TEST_FILE_XML} \
  -- -coverprofile=${COVERAGE_FILE} ./...
else
  go test -v ./...
fi
