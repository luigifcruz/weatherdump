CD = $(shell pwd)
VER = $(shell git describe --tags)

all: compiler build

fix-permission:
	sudo chown -fR $(shell whoami) dist/* || :
	sudo chown -fR $(shell whoami) release-builds/* || :
	sudo chown -fR $(shell whoami) gui/* || :

compiler:
	cd ./docker && docker build -t weatherdump_linux_x64 -f Dockerfile.linux_x64 .
	cd ./docker && docker build -t weatherdump_linux_armv7l -f Dockerfile.linux_armv7l .
	cd ./docker && docker build -t weatherdump_linux_armv6 -f Dockerfile.linux_armv6 .
	cd ./docker && docker build -t weatherdump_win_x64 -f Dockerfile.win_x64 .
	cd ./docker && docker build -t weatherdump_mac_x64 -f Dockerfile.mac_x64 .

build:
	mkdir -p release-builds ./dist
	rm -fr ./release-builds/weatherdump-cli-* ./dist/*
	docker run -v $(CD):/home/go/src/weather-dump weatherdump_linux_x64
	docker run -v $(CD):/home/go/src/weather-dump weatherdump_linux_armv7l
	docker run -v $(CD):/home/go/src/weather-dump weatherdump_linux_armv6
	docker run -v $(CD):/home/go/src/weather-dump weatherdump_win_x64
	docker run -v $(CD):/home/go/src/weather-dump weatherdump_mac_x64
	make fix-permission
	mv ./dist/export/* ./release-builds
	rm -fr ./dist/export

release:
	ghr -n $(VER) $(VER) ./release-builds

prerelease:
	ghr -prerelease -n $(VER) $(VER) ./release-builds 

draft:
	ghr -draft -n $(VER) $(VER) ./release-builds 

clean:
	make fix-permission
	rm -rf ./dist ./release-builds
