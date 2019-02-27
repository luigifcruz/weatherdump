#!/bin/bash

PACKAGE_NAME=weatherdump-cli-win-x64

mkdir -p ./dist ./dist/$PACKAGE_NAME ./dist/export
cp /usr/lib/gcc/x86_64-w64-mingw32/6.3-win32/libstdc++-6.dll ./dist/$PACKAGE_NAME
cp /usr/lib/gcc/x86_64-w64-mingw32/6.3-win32/libgcc_s_seh-1.dll ./dist/$PACKAGE_NAME
CGO_ENABLED=1 CXX=x86_64-w64-mingw32-g++ CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -o dist/$PACKAGE_NAME/weatherdump.exe ./main.go
cd ./dist && zip $PACKAGE_NAME.zip ./$PACKAGE_NAME/*
cd - && mv ./dist/$PACKAGE_NAME.zip ./dist/export