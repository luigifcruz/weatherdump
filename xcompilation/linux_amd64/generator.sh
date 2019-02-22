#!/bin/bash

PACKAGE_NAME=weatherdump_linux_amd64

mkdir ./bin ./bin/$PACKAGE_NAME ./bin/export
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o bin/$PACKAGE_NAME/weatherdump ./main.go
cd ./bin && tar --xform s:'./':: -czvf $PACKAGE_NAME.tar.gz ./$PACKAGE_NAME
cd - && mv ./bin/$PACKAGE_NAME.tar.gz ./bin/export