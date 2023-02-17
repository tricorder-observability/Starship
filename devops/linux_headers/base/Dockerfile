FROM ubuntu:18.04

ARG KERN_VERSION

RUN apt-get update
RUN apt-get upgrade -y -q
RUN apt-get install -y -q build-essential \
  bc \
  libelf-dev \
  libssl-dev \
  flex \
  bison \
  kmod \
  cpio \
  rsync \
  wget
