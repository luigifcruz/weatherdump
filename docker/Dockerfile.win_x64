FROM golang:1.12

RUN apt-get update && apt-get install -y zip build-essential g++-mingw-w64 cmake --no-install-recommends

ENV CC="x86_64-w64-mingw32-gcc"
ENV CXX="x86_64-w64-mingw32-g++"
ENV PACKAGE_NAME="weatherdump-cli-win-x64"
ENV CGO_ENABLED="1"
ENV CGO_CFLAGS="-I/go/libaec/src"
ENV CGO_CXXFLAGS="-I/go/libsathelper/includes -I/go/libcorrect/build/include"
ENV CGO_LDFLAGS="-static -L/go/libaec/build/src -laec -L/go/libsathelper/build/lib -lSatHelper -L/usr/libcorrect/lib -lcorrect"
ENV GOOS="windows"
ENV GOARCH="amd64"
ENV GOPATH="/home/go"
ENV COMPRESS="zip"
ENV BINARY_NAME="weatherdump.exe"
ENV GO111MODULE=on

RUN git clone https://github.com/erget/libaec.git \
    && cd libaec \
    && mkdir build && cd build \
    && cmake -DCMAKE_INSTALL_PREFIX=/usr/x86_64-w64-mingw32 -DBUILD_SHARED_LIBS=OFF -DCMAKE_SYSTEM_NAME=Windows .. && make -j$(nproc) && make install \
    && cd ./../..
RUN git clone https://github.com/racerxdl/libcorrect.git \
    && cd libcorrect \
    && mkdir build && cd build \
    && cmake -DCMAKE_INSTALL_PREFIX=/usr/x86_64-w64-mingw32 -DCMAKE_SYSTEM_NAME=Windows .. && make -j$(nproc) && make install \
    && cd ./../..
RUN git clone https://github.com/luigifreitas/libsathelper.git \
    && cd libsathelper \
    && mkdir build && cd build \
    && cmake -DCMAKE_INSTALL_PREFIX=/usr/x86_64-w64-mingw32 -DCMAKE_SYSTEM_NAME=Windows -DARCHITECTURE=x86_64 .. && make -j$(nproc) && make install \
    && cd ./../..

WORKDIR /home/go/src/weather-dump

ADD generator.sh /go/generator.sh
RUN chmod +x /go/generator.sh && ls -lh
ENTRYPOINT ["/go/generator.sh"]