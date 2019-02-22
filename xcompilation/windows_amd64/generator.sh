#!/bin/bash

PACKAGE_NAME=weatherdump_windows_amd64

mkdir ./bin ./bin/$PACKAGE_NAME ./bin/export
cp /usr/lib/gcc/x86_64-w64-mingw32/6.3-win32/libstdc++-6.dll ./bin/$PACKAGE_NAME
cp /usr/lib/gcc/x86_64-w64-mingw32/6.3-win32/libgcc_s_seh-1.dll ./bin/$PACKAGE_NAME
CGO_ENABLED=1 CXX=x86_64-w64-mingw32-g++ CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -o bin/$PACKAGE_NAME/weatherdump.exe ./main.go
cd ./bin && zip $PACKAGE_NAME.zip ./$PACKAGE_NAME/*
cd - && mv ./bin/$PACKAGE_NAME.zip ./bin/export