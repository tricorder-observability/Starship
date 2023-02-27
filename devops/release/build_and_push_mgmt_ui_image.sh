#!/bin/bash -e

if [[ $# -lt 1 ]]; then
  echo "Need 1 argument as the tag, should be '$0 <tag>', exit ..."
  exit 1
fi

tag="$1"
REGISTRY=docker.io/tricorderobservability

# Execute in the sub-shell, not affecting the pwd
(cd ui && yarn install && yarn run build)

mgmt_ui_version_tag="${REGISTRY}/tricorder/mgmt-ui:${tag}"
docker build ui/ -f ui/docker/Dockerfile -t ${mgmt_ui_version_tag}
docker push ${mgmt_ui_version_tag}
