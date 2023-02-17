#!/bin/bash -e

# Build kernel headers for various LTS versions.
# A list of LTS linux kernel versions can be found at
# https://en.wikipedia.org/wiki/Linux_kernel_version_history

function build_linux_header() {
  build_dir="$1"
  output_dir="$2"
  ker_ver="$3"

  output_file_path="${output_dir}/linux-headers-${ker_ver}-starship.tar.gz"
  if [[ -f ${output_file_path} ]]; then
    echo "${output_file_path} already exists, skipping version ${ker_ver} ..."
    return
  fi

  echo "Building and outputting headers to $ToT/devops/linux_headers/output/"
  echo "for Kernel Version ${ker_ver} ..."

  builder_image="${ker_ver}-builder"
  docker build --build-arg KERN_VERSION="${ker_ver}" "${build_dir}" \
    -t "${builder_image}"
  docker run -it --rm -v ${root}/devops/linux_headers/output:/output \
    "${builder_image}"
}

root=$(git rev-parse --show-toplevel)
output_dir=${root}/devops/linux_headers/output
build_dir=${root}/devops/linux_headers/build

# Reads kernel versions file into array
readarray -t ker_vers < ${build_dir}/KERNEL_VERSIONS

for ker_ver in "${ker_vers[@]}"; do
  build_linux_header "${build_dir}" "${output_dir}" "${ker_ver}"
done

echo "Archiving kernel headers for versions: ${ker_vers[@]} ..."
all_headers_filename=linux-headers.tar.gz
(cd ${output_dir} && tar --transform "s/^/starship\/linux_headers\//" \
  -czf ${all_headers_filename} *.tar.gz timeconst_*.h)

all_headers_s3path="s3://tricorder-dev/${all_headers_filename}"
echo "Uploading all kernel headers archive to ${all_headers_s3path} ..."
aws s3 cp "${output_dir}/${all_headers_filename}" "${all_headers_s3path}"
