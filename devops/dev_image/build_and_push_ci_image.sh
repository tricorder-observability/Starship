#!/bin/bash

docker login -u tricorderobservability
make -C devops/dev_image build_and_push_ci_image
