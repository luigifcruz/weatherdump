FROM golang:1.12

RUN apt-get update && apt-get install -y clang libxml2-dev patch build-essential cmake --no-install-recommends

RUN git clone https://github.com/tpoechtrager/osxcross.git 
COPY tarball/* ./osxcross/tarballs/
RUN UNATTENDED=1 ./osxcross/build.sh 

ENV PATH="/go/osxcross/target/bin:${PATH}"
ENV CC="x86_64-apple-darwin15-clang"
ENV CXX="x86_64-apple-darwin15-clang++"
ENV MACOSX_DEPLOYMENT_TARGET="10.9"
ENV PACKAGE_NAME="weatherdump-cli-mac-x64"
ENV CGO_ENABLED="1"
ENV CGO_CFLAGS="-I/go/libaec/src"
ENV CGO_CXXFLAGS="-I/go/libsathelper/includes -I/go/libcorrect/build/include"
ENV CGO_LDFLAGS="-L/go/libaec/build/src -laec -L/go/libsathelper/build/lib -lSatHelper -L/go/libcorrect/build/lib -lcorrect"
ENV GOOS="darwin"
ENV GOARCH="amd64"
ENV GOPATH="/home/go"
ENV COMPRESS="tar.gz"
ENV BINARY_NAME="weatherdump"
ENV GO111MODULE=on

RUN git clone https://github.com/erget/libaec.git \
    && cd libaec \
    && mkdir build && cd build \
    && cmake -DBUILD_SHARED_LIBS=OFF -DCMAKE_SYSTEM_NAME=Darwin .. && make -j$(nproc) && make install \
    && cd ./../..
RUN git clone https://github.com/racerxdl/libcorrect.git \
    && cd libcorrect \
    && mkdir build && cd build \
    && cmake -DCMAKE_SYSTEM_NAME=Darwin .. && make -j$(nproc) && make install \
    && cd ./../..
RUN git clone https://github.com/luigifreitas/libsathelper.git \
    && cd libsathelper \
    && mkdir build && cd build \
    && cmake -DCMAKE_SYSTEM_NAME=Darwin -DARCHITECTURE=x86_64 .. && make -j$(nproc) && make install \
    && cd ./../..

RUN ls -l /go/libsathelper/build/lib/

WORKDIR /home/go/src/weather-dump

RUN rm -f /go/libaec/build/src/libaec.dylib
RUN rm -f /go/libcorrect/build/lib/libcorrect.dylib
RUN rm -f /go/libsathelper/build/lib/libsathelper.dylib

ADD generator.sh /go/generator.sh
RUN chmod +x /go/generator.sh && ls -lh
ENTRYPOINT ["/go/generator.sh"]