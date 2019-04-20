#!/bin/bash

mkdir -p ./dist ./dist/$PACKAGE_NAME ./dist/export
go build -ldflags "-X main.version=$(git describe --tags)" -o dist/$PACKAGE_NAME/$BINARY_NAME

if [ $COMPRESS = "tar.gz" ]; then
    cd ./dist && tar --xform s:'./':: -czvf $PACKAGE_NAME.tar.gz ./$PACKAGE_NAME
    cd - && mv ./dist/$PACKAGE_NAME.tar.gz ./dist/export
fi

if [ $COMPRESS = "zip" ]; then
    cd ./dist && zip $PACKAGE_NAME.zip ./$PACKAGE_NAME/*
    cd - && mv ./dist/$PACKAGE_NAME.zip ./dist/export
fi
