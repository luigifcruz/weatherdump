FROM golang:1.11

RUN apt-get update && \
    apt-get install -y build-essential g++-mingw-w64 cmake 

RUN git clone https://github.com/ruslo/polly.git

RUN git clone https://github.com/erget/libaec.git \
    && cd libaec \
    && mkdir build && cd build \
    && cmake -DCMAKE_INSTALL_PREFIX=/usr/x86_64-w64-mingw32/  -DCMAKE_TOOLCHAIN_FILE=../../polly/linux-mingw-w64-gnuxx11.cmake -DBUILD_SHARED_LIBS=OFF .. && make -j8 && make install \
    && cd ./../..
RUN git clone https://github.com/racerxdl/libcorrect.git \
    && cd libcorrect \
    && mkdir build && cd build \
    && cmake -DCMAKE_INSTALL_PREFIX=/usr/x86_64-w64-mingw32/ -DCMAKE_TOOLCHAIN_FILE=../../polly/linux-mingw-w64-gnuxx11.cmake -DBUILD_SHARED_LIBS=OFF .. && make -j8 && make install \
    && cd ./../..
RUN git clone https://github.com/luigifreitas/libsathelper.git \
    && cd libsathelper \
    && mkdir build && cd build \
    && cmake -DCMAKE_INSTALL_PREFIX=/usr/x86_64-w64-mingw32/ -DARCHITECTURE=x86_64 -DCMAKE_TOOLCHAIN_FILE=../../polly/linux-mingw-w64-gnuxx11.cmake -DBUILD_SHARED_LIBS=OFF .. && make -j8 && make install \
    && cd ./../..
RUN git clone https://github.com/luigifreitas/mingw-std-threads.git \
    && cd mingw-std-threads \
    && mkdir /usr/x86_64-w64-mingw32/include/mingw32 \
    && mv ./*.h /usr/x86_64-w64-mingw32/include/mingw32

ENV GOPATH=/home/go
WORKDIR /home/go/src/weather-dump

RUN CGO_ENABLED=1 CXX=x86_64-w64-mingw32-g++ CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go get -v github.com/gorilla/mux github.com/satori/go.uuid github.com/urfave/cli github.com/nfnt/resize github.com/luigifreitas/libsathelper github.com/gorilla/websocket
ENTRYPOINT CGO_ENABLED=1 CXX=x86_64-w64-mingw32-g++ CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -o bin/weatherdump_amd64_windows ./main.go