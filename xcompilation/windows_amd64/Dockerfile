FROM golang:1.11

ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update && \
    apt-get install -y wget build-essential mingw-w64 g++-mingw-w64 wget cmake --no-install-recommends

RUN wget https://raw.githubusercontent.com/zyga/cmake-toolchains/master/Toolchain-Ubuntu-mingw32.cmake

RUN git clone https://github.com/erget/libaec.git \
    && cd libaec \
    && mkdir build && cd build \
    && cmake -DCMAKE_TOOLCHAIN_FILE=../../Toolchain-Ubuntu-mingw32.cmake -DBUILD_SHARED_LIBS=OFF .. && make -j8 && make install \
    && cd ./../..
RUN git clone https://github.com/quiet/libcorrect.git \
    && cd libcorrect \
    && mkdir build && cd build \
    && cmake -DCMAKE_TOOLCHAIN_FILE=../../Toolchain-Ubuntu-mingw32.cmake .. && make -j8 && make install \
    && cd ./../..
RUN git clone https://github.com/opensatelliteproject/libsathelper.git \
    && cd libsathelper \
    && mkdir build && cd build \
    && cmake -DCMAKE_TOOLCHAIN_FILE=../../Toolchain-Ubuntu-mingw32.cmake .. && make -j8 && make install \
    && cd ./../..

ENV GOPATH=/home/go
WORKDIR /home/go/src/weather-dump

RUN go get -v github.com/urfave/cli github.com/nfnt/resize github.com/OpenSatelliteProject/libsathelper github.com/gorilla/websocket

ENTRYPOINT CGO_ENABLED=1 CXX=i686-w64-mingw32-g++ CC=i686-w64-mingw32-gcc GOOS=windows GOARCH=386 go build -o bin/weatherdump_amd64_linux ./main.go