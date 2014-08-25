#!/bin/bash

if [ $# -eq 0 ]; then
    echo "Usage: wq-update.sh <cctools-source-tarball>"
    exit 1
fi

# untar and remove possibly prev version of work queue files
tar -xzf $@
rm -f work_queue
mkdir work_queue

# copy over needed files
cp cctools-*-source/work_queue/src/*.c work_queue/
cp cctools-*-source/work_queue/src/*.h work_queue/
cp cctools-*-source/dttools/src/*.c work_queue/
cp cctools-*-source/dttools/src/*.h work_queue/
cp cctools-*-source/chirp/src/*.c work_queue/
cp cctools-*-source/chirp/src/*.h work_queue/

# comment out internal includes
list=$(ls work_queue/*.h)
while read -r header; do
    header=$(echo $header | sed 's/\//\\\//g')
    sed -i "s/#include *.\+$header/\/\/\0/" $(ls work_queue/*)
done <<< "$list"

# amalgamate and cleanup
cat work_queue/*.c > wq_all.c
cat work_queue/*.h > wq_all.h

rm -rf cctools-*-source
rm -rf work_queue
