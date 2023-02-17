FROM ubuntu:22.04

ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update -y
RUN apt-get upgrade -y

# Needed for the next gem install command
RUN apt-get install -y --fix-missing ruby ruby-dev rubygems build-essential

RUN gem install --no-document fpm

# TODO(yzhao): Suspect these installed packages got linked into llvm below.
# Eventually got picked up by BCC. BCC needs zlib1g-dev, libelf-dev.
RUN apt-get install -y --fix-missing bison build-essential cmake flex git libedit-dev \
  clang libclang-dev llvm llvm-dev \
  python3 python3-distutils swig libncurses5-dev zlib1g-dev libelf-dev subversion \
  gcc-multilib

ENV CC=clang
ENV CXX=clang++
ENV BASE_DIR=/opt/tricorder/lib

WORKDIR /llvm_all
RUN git clone --branch llvmorg-14.0.4 --depth 1 https://github.com/llvm/llvm-project.git

# This builds libclang.a, which is needed for BCC static linking.
WORKDIR /llvm_all/build
RUN triple=$(gcc -v 2>&1 | grep "^Target:" | cut -d ' ' -f 2) && \
    cmake -G "Unix Makefiles" -DCMAKE_BUILD_TYPE=Release \
    -DLLVM_BUILD_DOCS=OFF -DCMAKE_INSTALL_PREFIX=${BASE_DIR}/clang-14.0 \
    -DLLVM_INCLUDE_EXAMPLES=OFF \
    -DLLVM_BUILD_32_BITS=OFF \
    -DLLVM_TARGETS_TO_BUILD="BPF;X86;AArch64" \
    -DLLVM_ENABLE_PER_TARGET_RUNTIME_DIR=OFF \
    -DLLVM_ENABLE_SPHINX=OFF \
    -DLLVM_ENABLE_DOXYGEN=OFF \
    -DLLVM_ENABLE_RTTI=ON \
    -DCLANG_INCLUDE_TESTS=OFF \
    -DLIBCLANG_BUILD_STATIC=ON \
    -DLLVM_INCLUDE_TESTS=OFF \
    -DLLVM_DEFAULT_TARGET_TRIPLE=${triple} \
    -DLLVM_ENABLE_PROJECTS="clang;clang-tools-extra;polly;lld;lldb;openmp;libcxx;libcxxabi;compiler-rt" \
    ../llvm-project/llvm

RUN make -j $(nproc)
RUN make -j $(nproc) runtimes
RUN make install

# This is the installation directory
WORKDIR ${BASE_DIR}
VOLUME /image
ENV PKG_NAME clang-14_0-tricorder
ENV DEB_NAME ${PKG_NAME}.deb
CMD ["sh", "-c",  "fpm -p /image/${DEB_NAME} \
        -s dir -t deb -n ${PKG_NAME} -v 14.0 --prefix ${BASE_DIR} clang-14.0"]
