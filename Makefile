CD = $(shell pwd)

release: build-cli-compiler build-cli-release build-gui-release

fix-permission:
	sudo chown -fR $(shell whoami) dist/* || :
	sudo chown -fR $(shell whoami) release-builds/* || :
	sudo chown -fR $(shell whoami) gui/* || :

build-cli-compiler:
	cd ./docker && docker build -t weatherdump_linux_x64 -f Dockerfile.linux_x64 .
	cd ./docker && docker build -t weatherdump_linux_armv7 -f Dockerfile.linux_armv7 .
	cd ./docker && docker build -t weatherdump_linux_armv6 -f Dockerfile.linux_armv6 .
	cd ./docker && docker build -t weatherdump_win_x64 -f Dockerfile.win_x64 .
	cd ./docker && docker build -t weatherdump_mac_x64 -f Dockerfile.mac_x64 .

build-cli-release:
	mkdir -p release-builds ./dist
	rm -fr ./release-builds/weatherdump-cli-* ./dist/*
	docker run -v $(CD):/home/go/src/weather-dump weatherdump_linux_x64
	docker run -v $(CD):/home/go/src/weather-dump weatherdump_linux_armv7
	docker run -v $(CD):/home/go/src/weather-dump weatherdump_linux_armv6
	docker run -v $(CD):/home/go/src/weather-dump weatherdump_win_x64
	docker run -v $(CD):/home/go/src/weather-dump weatherdump_mac_x64
	make fix-permission
	mv ./dist/export/* ./release-builds
	rm -fr ./dist/export

build-gui-release-compiler:
	cd ./docker && docker build -t weatherdump_gui -f Dockerfile.gui .

build-gui-release:
	docker run -v $(CD):/weather-dump weatherdump_gui

build-gui:
	mkdir -p release-builds
	make build-webpack
	make build-gui-linux
	make build-gui-windows
	make build-gui-mac
	rm -fr ./gui/dist

build-gui-linux:
	electron-builder --project ./gui -l --x64 --armv7l
	mv ./gui/dist/*.AppImage ./release-builds

build-gui-windows:
	electron-builder --project ./gui -w --x64
	mv ./gui/dist/*.exe ./release-builds

build-gui-mac:
	electron-builder --project ./gui -m --x64
	mv ./gui/dist/*.zip ./release-builds

build-webpack:
	cd ./gui && npm i && npm run build && cd -

clean:
	make fix-permission
	rm -rf ./dist ./gui/dist ./gui/node_modules
	rm -rf ./gui/resources/*.css ./gui/resources/*.js
