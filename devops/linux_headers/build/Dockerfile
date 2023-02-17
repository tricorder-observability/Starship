# This base image was built by devops/linux_headers_build_image/Dockerfile
# Avoids repeatitively install packages
FROM public.ecr.aws/tricorder/linux_headers_build_image

ARG KERN_VERSION

WORKDIR /tmp/src
RUN KERN_VERSION_MAJOR="$(echo "${KERN_VERSION}" | cut -d'.' -f1)"; \
    wget -nv http://mirrors.edge.kernel.org/pub/linux/kernel/v${KERN_VERSION_MAJOR}.x/linux-${KERN_VERSION}.tar.gz
RUN tar zxf linux-${KERN_VERSION}.tar.gz

WORKDIR /tmp/src/linux-${KERN_VERSION}
ADD config .config
RUN make olddefconfig
RUN make -j $(nproc) deb-pkg

WORKDIR /tmp
RUN dpkg -x src/linux-headers-${KERN_VERSION}_${KERN_VERSION}-1_amd64.deb .

RUN find usr/src/linux-headers-${KERN_VERSION} -xtype l -exec rm {} +

# Remove uneeded files to reduce size
# Keep only:
# - usr/src/linux-headers-x.x.x-pl/include
# - usr/src/linux-headers-x.x.x-pl/arch/x86
# This reduces the size by a little over 2x.
RUN rm -rf usr/share
RUN find usr/src/linux-headers-${KERN_VERSION} -maxdepth 1 -mindepth 1 ! \
    -name include ! -name arch -type d -exec rm -rf {} +
RUN find usr/src/linux-headers-${KERN_VERSION}/arch -maxdepth 1 -mindepth 1 ! \
    -name x86 -type d -exec rm -rf {} +
RUN mv usr/src/linux-headers-${KERN_VERSION} \
       usr/src/linux-headers-${KERN_VERSION}-starship
RUN tar zcf linux-headers-${KERN_VERSION}-starship.tar.gz usr

# Runner of this image should mount the host directory to this directory.
VOLUME /output
CMD ["sh", "-c", "cp linux-headers-*.tar.gz /output"]
