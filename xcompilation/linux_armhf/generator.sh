#!/bin/bash

PACKAGE_NAME=weatherdump_linux_armhf

mkdir ./bin ./bin/$PACKAGE_NAME ./bin/export
CXX=arm-linux-gnueabihf-g++ CC=arm-linux-gnueabihf-gcc GOOS=linux GOARCH=arm GOARM=6 CGO_ENABLED=1 go build -o bin/$PACKAGE_NAME/weatherdump ./main.go
cd ./bin && tar --xform s:'./':: -czvf $PACKAGE_NAME.tar.gz ./$PACKAGE_NAME
cd - && mv ./bin/$PACKAGE_NAME.tar.gz ./bin/export