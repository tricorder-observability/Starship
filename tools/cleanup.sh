#!/bin/bash

tools_root_dir=$(dirname "$(pwd)/$0")

function print_divider() {
  echo "============================"
}

print_divider
echo "Running gazelle ..."
print_divider

${tools_root_dir}/scripts/gazelle.sh

echo
print_divider
echo "Running buildifier ..."
print_divider

${tools_root_dir}/scripts/buildifier.sh

echo
print_divider
echo "Running golangci-lint ..."
print_divider
# Run golangci-lint and fix issues. Config file is .golangci.yml
golangci-lint run --fix

echo
print_divider
echo "Running golines ..."
print_divider
golines --max-len=120 --write-output src/**/*.go

echo
print_divider
echo "Running clang-format ..."
print_divider
find ./ -name '*.h' -o -name '*.c' | xargs clang-format -i

echo
print_divider
echo "Running check_readme ..."
print_divider
.github/scripts/check_readme.sh

echo
print_divider
echo "Running check_markdown_naming ..."
print_divider
.github/scripts/check_markdown_filename.sh

echo
print_divider
echo "Running check_dir_naming ..."
print_divider
.github/scripts/check_dir_naming.sh

echo
print_divider
echo "Running check_todo ..."
print_divider
.github/scripts/check_todo.sh

echo
print_divider
echo "Running check_license ..."
print_divider
make -C devops/license/ addlicense
