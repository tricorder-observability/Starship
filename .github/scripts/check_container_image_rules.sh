#!/bin/bash -e

errors=()
container_image_bzl="//:bazel/container_image.bzl"

container_image_rules=("go_image" "container_image" "container_push" "pkg_tar")
for rule in "${container_image_rules[@]}"; do
  echo "Checking rule: ${rule} ..."
  pattern="load\\(\"([^,]+)\",.*\"${rule}\".*\\)"
  for build_file in src/**/BUILD.bazel; do
    echo "Check BUILD.bazel file: ${build_file} ..."

    # Find the loaded .bzl files that load container image rules.
    mapfile -t loaded_files < <(sed -nr "s/${pattern}/\\1/p" "${build_file}")
    for loaded_file in "${loaded_files[@]}"; do
      if [[ "${loaded_file}" != "" && "${loaded_file}" != "${container_image_bzl}" ]]; then
        # Record loaded file that is not the desired one
        errors+=("BUILD.bazel:${build_file} Loaded file:${loaded_file} Loaded rule:${rule}")
      fi
    done
  done
done

if [[ ${#errors[@]} != 0 ]]; then
  echo
  echo "Some BUILD.bazel files use wrong container_image rule file ..."
  echo "For all container image related rules, go_image container_image container_push pkg_tar etc.,"
  echo "use ${container_image_bzl}, not the official bazel bzl"
  echo
  echo "Problematic BUILD.bazel:"
  for culprit_build_file in "${errors[@]}"; do
    echo "${culprit_build_file}"
  done

  echo
  echo "Replace the loaded .bzl file with ${container_image_bzl}!"
  exit 1
fi

echo
echo "PASS"
