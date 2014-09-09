#!/bin/bash

if [ $# -eq 0 ]; then
    echo "Usage: wq-update.sh <cctools-source-tarball>"
    exit 1
fi

cctools="$(pwd)/$(echo "$@" | sed -r 's/\.[[:alnum:]]+\.[[:alnum:]]+$//')"

wqdir=$(pwd)

# untar and remove possibly prev version of work queue files
tar -xzf $@

# copy over needed files

cd $cctools
sed -i '/HAS_SYSTEMD_JOURNAL/d' configure
sed -i '/-lsystemd/{N;d;}' configure
./configure
make

