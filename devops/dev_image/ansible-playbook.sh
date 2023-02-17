#!/bin/bash -ex

# A wrapper script to run ansible-playbook, which silences warning
ANSIBLE_LOCALHOST_WARNING=False \
ansible-playbook "$@"
