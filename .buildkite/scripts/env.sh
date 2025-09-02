#!/usr/bin/env bash
set -euo pipefail

BUILDKITE=${BUILDKITE:-false}
ARTIFACT_DIR=artifacts

if [ "$BUILDKITE" != "true" ]; then
  ARTIFACT_DIR=/tmp
fi

# Ensure directory exists
[ ! -d $ARTIFACT_DIR ] && mkdir -p $ARTIFACT_DIR || true

# Report files
LINT_FILE=$ARTIFACT_DIR/lint.txt
TEST_FILE=$ARTIFACT_DIR/test.txt
TEST_FILE_XML=$ARTIFACT_DIR/test.xml
COVERAGE_FILE=$ARTIFACT_DIR/cover.txt
COVERAGE_FILE_XML=$ARTIFACT_DIR/cover.xml
