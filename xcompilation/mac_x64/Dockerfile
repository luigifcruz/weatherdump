FROM golang:1.11

RUN apt-get update && apt-get install -y libxml2-dev patch build-essential cmake --no-install-recommends

RUN git clone https://github.com/tpoechtrager/osxcross.git 

COPY tarball/* ./osxcross/tarballs/

RUN ./osxcross/tools/get_dependencies.sh
RUN UNATTENDED=1 ./osxcross/build_clang.sh
RUN UNATTENDED=1 ./osxcross/build.sh 
ENV PATH="/go/osxcross/target/bin:${PATH}"

COPY toolchain.cmake .

RUN git clone https://github.com/erget/libaec.git \
    && cd libaec \
    && mkdir build && cd build \
    && cmake -DBUILD_SHARED_LIBS=OFF -DCMAKE_TOOLCHAIN_FILE=../../toolchain.cmake .. && make -j8 && make install \
    && cd ./../..
RUN git clone https://github.com/racerxdl/libcorrect.git \
    && cd libcorrect \
    && mkdir build && cd build \
    && cmake -DBUILD_SHARED_LIBS=OFF -DCMAKE_TOOLCHAIN_FILE=../../toolchain.cmake .. && make -j8 && make install \
    && cd ./../..
RUN git clone https://github.com/luigifreitas/libsathelper.git \
    && cd libsathelper \
    && mkdir build && cd build \
    && cmake -DBUILD_SHARED_LIBS=OFF -DARCHITECTURE=x86_64 -DCMAKE_TOOLCHAIN_FILE=../../toolchain.cmake .. && make -j8 && make install \
    && cd ./../..

ENV GOPATH=/home/go
WORKDIR /home/go/src/weather-dump

RUN rm -f /go/libcorrect/build/lib/libcorrect.dylib
RUN rm -f /go/libsathelper/build/lib/libsathelper.dylib

RUN CGO_ENABLED=1 CGO_CXXFLAGS="-I/go/libsathelper/includes -I/go/libcorrect/build/include" \
    MACOSX_DEPLOYMENT_TARGET=10.9 CGO_LDFLAGS="-L/go/libsathelper/build/lib -lsathelper -L/go/libcorrect/build/lib -lcorrect" \
    CC=o64-clang CXX=o64-clang++ GOOS=darwin GOARCH=amd64 \
    go get -v github.com/gorilla/handlers github.com/gorilla/mux github.com/satori/go.uuid gopkg.in/alecthomas/kingpin.v2 github.com/nfnt/resize github.com/luigifreitas/libsathelper github.com/gorilla/websocket

RUN cd /home/go/src/github.com && mkdir ./OpenSatelliteProject ./OpenSatelliteProject/libsathelper && mv ./luigifreitas/libsathelper/* ./OpenSatelliteProject/libsathelper

ADD generator.sh /go/generator.sh
RUN chmod +x /go/generator.sh && ls -lh
ENTRYPOINT ["/go/generator.sh"]