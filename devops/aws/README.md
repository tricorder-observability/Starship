# AWS

AWS devops operational documentation.

## Search resources on AWS

Use tag edit to find all resources being used on AWS:

https://ap-northeast-1.console.aws.amazon.com/resource-groups/tag-editor
Select all regions, and all resource types:
![image](https://user-images.githubusercontent.com/112656580/217785716-9e41be8a-349e-49f9-8f29-f38afd94cb36.png)

## Install kernel headers on Amazon AMI Linux

Kernel headers need to be present on the node to allow BCC to compile C code.
```
# Make sure EKS nodes has SSH enabled.
sudo yum install kernel-devel-$(uname -r) -y
# Verify that kernel headers are installed
ls /lib/modules/$(uname -r)/build
# The above is a symlink to /usr/src/kernels/$(uname -r)
# Installing kernel-devel basically installed kernel header in this /usr/src/kernels/$(uname -r) dir.
```

TODO: `kernel-devel` package includes additional tools like build tools. Might not be the most compact way of installing
kernel headers.
