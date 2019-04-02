FROM quay.io/bitriseio/android-ndk

RUN apt-get update && apt-get install -y build-essential cmake --no-install-recommends

ENV PATH="/opt/android-ndk/toolchains/llvm/prebuilt/linux-x86_64/bin:${PATH}"
ENV CC="armv7a-linux-androideabi22-clang"
ENV CXX="armv7a-linux-androideabi22-clang++"
ENV LD="/opt/android-ndk/toolchains/llvm/prebuilt/linux-x86_64/bin/arm-linux-androideabi-ld"
ENV AR="/opt/android-ndk/toolchains/llvm/prebuilt/linux-x86_64/bin/arm-linux-androideabi-ar"
ENV CGO_LDFLAGS="-Wl,-Bstatic -L/usr/local/lib -laec -L/usr/local/lib -lSatHelper -L/usr/local/lib -lcorrect -Wl,-Bdynamic -L/opt/android-ndk/platforms/android-22/arch-arm/usr/lib -llog"
ENV PACKAGE_NAME="weatherdump-cli-android-armv7a"
ENV CXXFLAGS="-I/usr/local/include"
ENV CGO_ENABLED="1"
ENV CGO_CFLAGS="-I/usr/local/include"
ENV CGO_CXXFLAGS="-I/usr/local/include"
ENV GOOS="android"
ENV GOARCH="arm"
ENV GOPATH="/home/go"
ENV COMPRESS="tar.gz"
ENV BINARY_NAME="weatherdump"

RUN git clone https://github.com/erget/libaec.git \
    && cd libaec \
    && mkdir build && cd build \
    && cmake -DCMAKE_AR=/opt/android-ndk/toolchains/llvm/prebuilt/linux-x86_64/bin/arm-linux-androideabi-ar -DCMAKE_SYSTEM_NAME=Android -DBUILD_SHARED_LIBS=OFF .. && make -j$(nproc) && make install \
    && cd ./../..
RUN git clone https://github.com/quiet/libcorrect.git \
    && cd libcorrect \
    && mkdir build && cd build \
    && cmake -DCMAKE_AR=/opt/android-ndk/toolchains/llvm/prebuilt/linux-x86_64/bin/arm-linux-androideabi-ar -DCMAKE_SYSTEM_NAME=Android -DBUILD_SHARED_LIBS=OFF .. && make -j$(nproc) && make install \
    && cd ./../..
RUN git clone https://github.com/luigifreitas/libsathelper.git \
    && cd libsathelper \
    && mkdir build && cd build \
    && cmake -DARCHITECTURE=armv7l-soft .. && make -j$(nproc) && make install \
    && cd ./../.. \
    && mkdir /usr/local/include/SatHelper && cp -R /usr/local/include/sathelper/* /usr/local/include/SatHelper

WORKDIR /home/go/src/weather-dump

ADD generator.sh /go/generator.sh
RUN chmod +x /go/generator.sh && ls -lh
ENTRYPOINT ["/go/generator.sh"]