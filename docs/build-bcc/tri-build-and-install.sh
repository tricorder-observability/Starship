#!/bin/bash
#
# Style guide:
# https://google.github.io/styleguide/shellguide.html

# Step 1. Check necessary kernel config options

git co -b v0.25.0 711f0302

mkdir build && cd build

cmake ..

# Using more threads to speedup the compilation
num_cpus=$(grep -c ^processor /proc/cpuinfo)
threads=1
if [ $num_cpus -gt 8 ]; then
    threads=8
elif [ $num_cpus -gt 2 ]; then
    threads=$num_cpus
fi
make -j $threads

sudo make install

cmake -DPYTHON_CMD=python3 .. # build python3 binding
pushd src/python/
make
sudo make install
popd

# Fix "/usr/bin/python: bad interpreter: No such file or directory" error
sudo ln -s /usr/bin/python3 /usr/bin/python
