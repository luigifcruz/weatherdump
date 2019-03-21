#!/bin/bash

# Golang Dependencies
go get -v github.com/fatih/color github.com/schollz/progressbar github.com/luigifreitas/gofast github.com/gorilla/handlers github.com/gorilla/mux github.com/satori/go.uuid gopkg.in/alecthomas/kingpin.v2 github.com/nfnt/resize github.com/luigifreitas/libsathelper github.com/gorilla/websocket
cd /home/go/src/github.com && mkdir ./OpenSatelliteProject ./OpenSatelliteProject/libsathelper && mv ./luigifreitas/libsathelper/* ./OpenSatelliteProject/libsathelper