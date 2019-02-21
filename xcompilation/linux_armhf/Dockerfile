FROM golang:1.11

RUN apt-get update && \
    apt-get install -y build-essential cmake g++-arm-linux-gnueabihf

RUN git clone https://github.com/ruslo/polly.git

RUN git clone https://github.com/erget/libaec.git \
    && cd libaec \
    && mkdir build && cd build \
    && cmake -DCMAKE_INSTALL_PREFIX=/usr/arm-linux-gnueabihf/ -DBUILD_SHARED_LIBS=OFF -DCMAKE_TOOLCHAIN_FILE=../../polly/linux-gcc-armhf-neon-vfpv4.cmake .. && make -j8 && make install \
    && cd ./../..
RUN git clone https://github.com/quiet/libcorrect.git \
    && cd libcorrect \
    && mkdir build && cd build \
    && cmake -DCMAKE_INSTALL_PREFIX=/usr/arm-linux-gnueabihf/ -DBUILD_SHARED_LIBS=OFF -DCMAKE_TOOLCHAIN_FILE=../../polly/linux-gcc-armhf-neon-vfpv4.cmake .. && make -j8 && make install \
    && cd ./../..
RUN git clone https://github.com/luigifreitas/libsathelper.git \
    && cd libsathelper \
    && mkdir build && cd build \
    && cmake -DCMAKE_INSTALL_PREFIX=/usr/arm-linux-gnueabihf/ -DARCHITECTURE=armv7l -DBUILD_SHARED_LIBS=OFF -DCMAKE_TOOLCHAIN_FILE=../../polly/linux-gcc-armhf-neon-vfpv4.cmake .. && make -j8 && make install && ldconfig \
    && cd ./../..

ENV GOPATH=/home/go
WORKDIR /home/go/src/weather-dump

RUN CXX=arm-linux-gnueabihf-g++ CC=arm-linux-gnueabihf-gcc GOOS=linux GOARCH=arm GOARM=6 CGO_ENABLED=1 go get -v github.com/gorilla/mux github.com/satori/go.uuid github.com/urfave/cli github.com/nfnt/resize github.com/OpenSatelliteProject/libsathelper github.com/gorilla/websocket
ENTRYPOINT CXX=arm-linux-gnueabihf-g++ CC=arm-linux-gnueabihf-gcc GOOS=linux GOARCH=arm GOARM=6 CGO_ENABLED=1 go build -o bin/weatherdump_armhf_linux ./main.go