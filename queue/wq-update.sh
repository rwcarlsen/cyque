#!/bin/bash

tar -xzf $@
rm -f work_queue/*.c work_queue/*.h
cp cctools-*-source/work_queue/src/*.c work_queue/
cp cctools-*-source/work_queue/src/*.h work_queue/
cp cctools-*-source/dttools/src/*.c work_queue/
cp cctools-*-source/dttools/src/*.h work_queue/
cp cctools-*-source/chirp/src/*.c work_queue/
cp cctools-*-source/chirp/src/*.h work_queue/
rm -rf cctools-*-source

