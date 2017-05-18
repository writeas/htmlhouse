#!/bin/bash

#
# Compile the application and put it all in one place.
# NOTE: also requires generated keys to run.
#

L="\033[0;35m==\033[00m"
S="\033[0;32m==\033[00m"

echo -e "$L Cleaning up previous build..."
rm -rf build

echo -e "$L Creating directories..."
mkdir -p build/static/css/
mkdir -p build/static/js/

echo -e "$L Generating CSS..."
make

echo -e "$L Copying files..."
cp static/*.html build/static/
cp static/css/*.css build/static/css/
cp -r static/js/* build/static/js/
cp -r templates/ build/

echo -e "$L Building application..."
cd main/
go build -o htmlhouse
mv htmlhouse ../build/

echo -e "$S Done."
