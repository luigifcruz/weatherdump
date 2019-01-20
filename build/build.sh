#!/usr/bin/env bash

export CGO_CFLAGS_ALLOW='-fopenmp'
export CGO_CFLAGS="`pkg-config --cflags MagickWand MagickCore` -I/usr/include/ImageMagick-6/magick"
export CGO_LDFLAGS="\
-lstdc++ -lXext -lX11 \
-Wl,-Bstatic \
    `pkg-config --libs MagickWand MagickCore` \
    -lgomp -laec -lIlmImf -lImath -lHalf -llcms -lIlmThread -lIex -lfreetype -lxml2 -lexpat -lfftw3 -lbz2 -lm -lz \
    -lwebp -ljpeg -lpng16 -ltiff -lgif \
-Wl,-Bdynamic"

go get gopkg.in/gographics/imagick.v2/imagick
go build -tags no_pkgconfig gopkg.in/gographics/imagick.v2/imagick
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o bin/weatherapp ./main.go