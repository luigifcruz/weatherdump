#!/bin/bash

PACKAGE_NAME=weatherdump-cli-linux-x64

mkdir -p ./dist ./dist/$PACKAGE_NAME ./dist/export
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o dist/$PACKAGE_NAME/weatherdump ./main.go
cd ./dist && tar --xform s:'./':: -czvf $PACKAGE_NAME.tar.gz ./$PACKAGE_NAME
cd - && mv ./dist/$PACKAGE_NAME.tar.gz ./dist/export