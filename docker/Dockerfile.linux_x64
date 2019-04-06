FROM golang:1.12

RUN apt-get update && apt-get install -y wget build-essential cmake --no-install-recommends

ENV PACKAGE_NAME="weatherdump-cli-linux-x64"
ENV CGO_ENABLED="1"
ENV CGO_CFLAGS="-I/go/libaec/src"
ENV CGO_CXXFLAGS="-I/go/libsathelper/includes -I/go/libcorrect/build/include"
ENV CGO_LDFLAGS="-L/go/libaec/build/src -laec -L/go/libsathelper/build/lib -lSatHelper -L/usr/local/lib -lcorrect"
ENV COMPRESS="tar.gz"
ENV BINARY_NAME="weatherdump"
ENV GOOS="linux"
ENV GOPATH="/home/go"
ENV GO111MODULE=on

RUN git clone https://github.com/erget/libaec.git \
    && cd libaec \
    && mkdir build && cd build \
    && cmake -DBUILD_SHARED_LIBS=OFF .. && make -j$(nproc) && make install \
    && cd ./../..
RUN git clone https://github.com/quiet/libcorrect.git \
    && cd libcorrect \
    && mkdir build && cd build \
    && cmake .. && make -j$(nproc) && make install \
    && cd ./../..
RUN git clone https://github.com/luigifreitas/libsathelper.git \
    && cd libsathelper \
    && mkdir build && cd build \
    && cmake .. && make -j$(nproc) && make install \
    && cd ./../..

WORKDIR /home/go/src/weather-dump

RUN rm -f /usr/local/lib/libcorrect.so

ADD generator.sh /go/generator.sh
RUN chmod +x /go/generator.sh && ls -lh
ENTRYPOINT ["/go/generator.sh"]