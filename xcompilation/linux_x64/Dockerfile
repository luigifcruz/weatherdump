FROM golang:1.11

RUN apt-get update && apt-get install -y wget build-essential cmake --no-install-recommends

RUN git clone https://github.com/erget/libaec.git \
    && cd libaec \
    && mkdir build && cd build \
    && cmake -DBUILD_SHARED_LIBS=OFF .. && make -j8 && make install \
    && cd ./../..
RUN git clone https://github.com/quiet/libcorrect.git \
    && cd libcorrect \
    && mkdir build && cd build \
    && cmake .. && make -j8 && make install \
    && cd ./../..
RUN git clone https://github.com/opensatelliteproject/libsathelper.git \
    && cd libsathelper \
    && mkdir build && cd build \
    && cmake .. && make -j8 && make install \
    && cd ./../..

ENV GOPATH=/home/go
WORKDIR /home/go/src/weather-dump

RUN go get -v github.com/gorilla/handlers github.com/gorilla/mux github.com/satori/go.uuid gopkg.in/alecthomas/kingpin.v2 github.com/nfnt/resize github.com/OpenSatelliteProject/libsathelper github.com/gorilla/websocket

ADD generator.sh /go/generator.sh
RUN chmod +x /go/generator.sh && ls -lh
ENTRYPOINT ["/go/generator.sh"]