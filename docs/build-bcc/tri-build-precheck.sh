#!/bin/bash
#
# Style guide:
# https://google.github.io/styleguide/shellguide.html

# Utils
err() {
    echo "[$(date +'%Y-%m-%dT%H:%M:%S%z')]: $*" >&2
}

fatal() {
    echo "[$(date +'%Y-%m-%dT%H:%M:%S%z')]: $*" >&2
    exit 1
}

# Script starts
echo "Pre-checking system information of current node ..."
echo

# Step 1. Check necessary kernel configurations
kernel_version=$(uname -r)
kernel_config_file=/boot/config-$kernel_version
echo "Kernel version: $kernel_version"
echo

# https://github.com/iovisor/bcc/blob/master/INSTALL.md#kernel-configuration
# CONFIG_BPF=y
# CONFIG_BPF_SYSCALL=y
# CONFIG_NET_CLS_BPF=m   # [optional, for tc filters]
# CONFIG_NET_ACT_BPF=m   # [optional, for tc actions]
# CONFIG_BPF_JIT=y
# CONFIG_HAVE_EBPF_JIT=y # [for Linux kernel versions 4.7 and later]
# CONFIG_BPF_EVENTS=y    # [optional, for kprobes]
# CONFIG_IKHEADERS=y     # Need kernel headers through /sys/kernel/kheaders.tar.xz
#
# CONFIG_NET_SCH_SFQ=m
# CONFIG_NET_ACT_POLICE=m
# CONFIG_NET_ACT_GACT=m
# CONFIG_DUMMY=m
# CONFIG_VXLAN=m
expected_config="CONFIG_BPF=y CONFIG_BPF_SYSCALL=y CONFIG_NET_CLS_BPF=m CONFIG_NET_ACT_BPF=m CONFIG_BPF_JIT=y CONFIG_HAVE_EBPF_JIT=y CONFIG_BPF_EVENTS=y"
optional_config="CONFIG_IKHEADERS=y CONFIG_NET_SCH_SFQ=m CONFIG_NET_ACT_POLICE=m CONFIG_NET_ACT_GACT=m CONFIG_DUMMY=m CONFIG_VXLAN=m"

echo "Checking kernel configurations in $kernel_config_file"
for c in $expected_config; do
    grep $c $kernel_config_file || fatal "$c not meet, please check it by hand"
done

echo "BPF related kernel configurations OK"
echo

# Step 2. Check or install build toolchain
# https://github.com/iovisor/bcc/blob/master/INSTALL.md#ubuntu---source

echo "Installing packages ..."
# For Ubuntu 22.04
sudo apt update -y # a necessay step otherwise you may get errors like "libllvm14-dev not found"
sudo apt install -y bison build-essential cmake flex git libedit-dev \
    libllvm14 llvm-14-dev libclang-14-dev python3 zlib1g-dev libelf-dev libfl-dev python3-distutils
