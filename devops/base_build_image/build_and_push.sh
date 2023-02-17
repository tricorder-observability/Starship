#!/bin/bash -ex

aws ecr-public get-login-password --region us-east-1 |\
    docker login --username AWS --password-stdin public.ecr.aws/tricorder
image="public.ecr.aws/tricorder/base_build_image:v0.1"
docker build . -t ${image}
docker push ${image}
