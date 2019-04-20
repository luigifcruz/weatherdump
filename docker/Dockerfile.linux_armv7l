FROM golang:1.12

RUN apt-get update && apt-get install -y build-essential cmake g++-arm-linux-gnueabihf --no-install-recommends

ENV CC="arm-linux-gnueabihf-gcc"
ENV CXX="arm-linux-gnueabihf-g++"
ENV PACKAGE_NAME="weatherdump-cli-linux-armv7l"
ENV CGO_ENABLED="1"
ENV CGO_CFLAGS="-I/go/libaec/src"
ENV CGO_CXXFLAGS="-I/go/libsathelper/includes -I/go/libcorrect/build/include"
ENV CGO_LDFLAGS="-static -L/go/libaec/build/src -laec -L/go/libsathelper/build/lib -lSatHelper -L/usr/arm-linux-gnueabihf/lib -lcorrect"
ENV GOOS="linux"
ENV GOARCH="arm"
ENV GOARM="7"
ENV GOPATH="/home/go"
ENV COMPRESS="tar.gz"
ENV BINARY_NAME="weatherdump"
ENV GO111MODULE=on

RUN git clone https://github.com/erget/libaec.git \
    && cd libaec \
    && mkdir build && cd build \
    && cmake -DCMAKE_INSTALL_PREFIX=/usr/arm-linux-gnueabihf -DBUILD_SHARED_LIBS=OFF .. && make -j$(nproc) && make install \
    && cd ./../..
RUN git clone https://github.com/quiet/libcorrect.git \
    && cd libcorrect \
    && mkdir build && cd build \
    && cmake -DCMAKE_INSTALL_PREFIX=/usr/arm-linux-gnueabihf .. && make -j$(nproc) && make install \
    && cd ./../..
RUN git clone https://github.com/luigifreitas/libsathelper.git \
    && cd libsathelper \
    && mkdir build && cd build \
    && cmake -DCMAKE_INSTALL_PREFIX=/usr/arm-linux-gnueabihf -DARCHITECTURE=armv7l .. && make -j$(nproc) && make install \
    && cd ./../..

WORKDIR /home/go/src/weather-dump

ADD generator.sh /go/generator.sh
RUN chmod +x /go/generator.sh && ls -lh
ENTRYPOINT ["/go/generator.sh"]
