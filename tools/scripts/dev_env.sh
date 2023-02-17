#!/bin/bash

if [[ $(basename "$(pwd)") != "scripts" ]]; then
  echo "Must be run inside scripts directory, currently in $(pwd), exiting ..."
  exit 1
fi

echo "Installing ansible ..."
sudo apt-add-repository -y ppa:ansible/ansible
sudo apt-get update -y
sudo apt-get install -y ansible

echo "Ansible updating local dev environment ..."
sudo ansible-playbook ansible/dev.yaml

# TODO(yzhao): This is not used, as we setup dev environment with Pixie's chef.
# Replicate Pixie's chef in ansible, and remove this or adjust accordingly.
echo "Install pre-built llvm for BCC ..."
deb_name="clang-14_0-tricorder-0.deb"
s3_path="s3://tricorder-dev/${deb_name}"
aws s3 cp ${s3_path} ${deb_name}
sudo dpkg --install ${deb_name}

echo "Put the text below to your rc file to include the new directories ..."
echo "  export PATH=\"/opt/tricorder/bin:\$PATH\""
