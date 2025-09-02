#!/bin/sh
curl -H "Authorization: Bearer $TOKEN" -X POST https://api.buildkite.com/v2/organizations/omise/pipelines \
  -d '{
  "name": "$1",
  "repository": "git@git.omise.co:omise/$1",
  "steps": [
    {
      "type": "script",
      "name": "buildkite",
      "command": "buildkite-agent pipeline upload",
      "branch_configuration": "master",
      "agent_query_rules": ["queue=omise", "docker-builder=true"]
    }
  ]
}'

