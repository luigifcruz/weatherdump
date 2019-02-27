#!/bin/bash

PACKAGE_NAME=weatherdump-cli-mac-x64

mkdir -p ./dist ./dist/$PACKAGE_NAME ./dist/export

CGO_ENABLED=1 CGO_CFLAGS="-I/go/libaec/src" CGO_CXXFLAGS="-I/go/libsathelper/includes -I/go/libcorrect/build/include" \
MACOSX_DEPLOYMENT_TARGET=10.9 CGO_LDFLAGS="-L/go/libaec/build/src -laec -L/go/libsathelper/build/lib -lsathelper -L/go/libcorrect/build/lib -lcorrect" \
CC=o64-clang CXX=o64-clang++ GOOS=darwin GOARCH=amd64 \
go build -o dist/$PACKAGE_NAME/weatherdump ./main.go

cd ./dist && tar --xform s:'./':: -czvf $PACKAGE_NAME.tar.gz ./$PACKAGE_NAME
cd - && mv ./dist/$PACKAGE_NAME.tar.gz ./dist/export