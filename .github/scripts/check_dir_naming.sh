#!/bin/bash

echo "Checking directory names include only lower case and '-'"
found_breakage=false
for dirname in $(find src helm-charts tools -type d); do
  fname=$(basename ${dirname})
  if ! [[ ${fname} =~ ^[a-z-]+$ ]]; then
    found_breakage=true
    echo "${dirname}"
  fi
done

if [[ "${found_breakage}" == "true" ]]; then
  exit 1
fi
