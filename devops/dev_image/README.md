# DevImage

This docker image is same as the development environment. Used in
```
.github/workflows/build-and-test.yml
```
Every time the dev image is updated, fill in the version tag to the above file.

## Ansible

Ansible machine provisioning scripts, and dev image builder (used on github
actions).

**Ansible** cannot be executed without sudo.
It seems the original installation is done with `root` user in chef.
And that setup a Python environment that requires root permission.

* `make build_and_push_ci_image` builds the image used in GitHub CI, which uses
  `packer_ci_image.json`, which in turn uses ansible playbook `apt.yaml` and `install.yaml`
  included in `ci.yaml`.
* `make build_and_push_dev _image` builds the dev image on top of the base image,
  which uses `packer_dev_image.json`, which in turn uses the above playbook plus
  `apt_extra.yaml` and `install_extra.yaml` included in `dev.yaml`.
* `sudo ansible-playbook devops/dev_image/dev.yaml` to provision a dev workstation.
  Everything will be installed under `/opt/tricorder`.
