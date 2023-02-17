#!/bin/bash -ex

absolute_path="/etc/nginx/conf.d/default.conf"

# Backup the original config
mv ${absolute_path} ${absolute_path}.orig
cat ${absolute_path}.orig |
    sed s#HELM_RELEASE_NAME#${HELM_RELEASE_NAME}#g  >  ${absolute_path}

