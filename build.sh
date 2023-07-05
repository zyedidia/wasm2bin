#!/bin/sh

cd w2c2/w2c2
cmake -B build
cmake --build build
cd -

cd w2c2/wasi
cmake -B build
cmake --build build
cd -
