FROM golang:1.12

RUN apt-get update && apt-get install -y build-essential cmake g++-arm-linux-gnueabi --no-install-recommends

ENV CC="arm-linux-gnueabi-gcc"
ENV CXX="arm-linux-gnueabi-g++"
ENV PACKAGE_NAME="weatherdump-cli-linux-armv6"
ENV CGO_ENABLED="1"
ENV CGO_CFLAGS="-I/go/libaec/src"
ENV CGO_CXXFLAGS="-I/go/libsathelper/includes -I/go/libcorrect/build/include"
ENV CGO_LDFLAGS="-static -L/go/libaec/build/src -laec -L/go/libsathelper/build/lib -lSatHelper -L/usr/arm-linux-gnueabi/lib -lcorrect"
ENV GOOS="linux"
ENV GOARCH="arm"
ENV GOARM="6"
ENV GOPATH="/home/go"
ENV COMPRESS="tar.gz"
ENV BINARY_NAME="weatherdump"
ENV GO111MODULE=on

RUN git clone https://github.com/erget/libaec.git \
    && cd libaec \
    && mkdir build && cd build \
    && cmake -DCMAKE_INSTALL_PREFIX=/usr/arm-linux-gnueabi -DBUILD_SHARED_LIBS=OFF .. && make -j$(nproc) && make install \
    && cd ./../..
RUN git clone https://github.com/quiet/libcorrect.git \
    && cd libcorrect \
    && mkdir build && cd build \
    && cmake -DCMAKE_INSTALL_PREFIX=/usr/arm-linux-gnueabi .. && make -j$(nproc) && make install \
    && cd ./../..
RUN git clone https://github.com/luigifreitas/libsathelper.git \
    && cd libsathelper \
    && mkdir build && cd build \
    && cmake -DCMAKE_INSTALL_PREFIX=/usr/arm-linux-gnueabi -DARCHITECTURE=armv7l-soft .. && make -j$(nproc) && make install \
    && cd ./../..

WORKDIR /home/go/src/weather-dump

ADD generator.sh /go/generator.sh
RUN chmod +x /go/generator.sh && ls -lh
ENTRYPOINT ["/go/generator.sh"]