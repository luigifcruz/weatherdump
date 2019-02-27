CD = $(shell pwd)

release: build-cli-compiler build-cli-release build-gui-release

fix-permission:
	sudo chown -fR $(shell whoami) dist/* || :
	sudo chown -fR $(shell whoami) release-builds/* || :

build-cli-compiler:
	docker build -t weatherdump_linux_x64 ./xcompilation/linux_x64
	docker build -t weatherdump_linux_armhf ./xcompilation/linux_armhf
	docker build -t weatherdump_win_x64 ./xcompilation/win_x64
	docker build -t weatherdump_mac_x64 ./xcompilation/mac_x64

build-cli-release:
	mkdir -p release-builds ./dist
	rm -fr ./release-builds/weatherdump-cli-* ./dist/*
	docker run -v $(CD):/home/go/src/weather-dump weatherdump_linux_x64
	docker run -v $(CD):/home/go/src/weather-dump weatherdump_linux_armhf
	docker run -v $(CD):/home/go/src/weather-dump weatherdump_win_x64
	docker run -v $(CD):/home/go/src/weather-dump weatherdump_mac_x64
	make fix-permission
	mv ./dist/export/* ./release-builds
	rm -fr ./dist/export

build-gui-release:
	mkdir -p release-builds
	make build-web-resources
	make build-gui-release-linux
	make build-gui-release-windows
	make build-gui-release-mac
	rm -fr ./gui/dist

build-gui-release-linux:
	electron-builder --project ./gui -l --x64
	mv ./gui/dist/*.AppImage ./release-builds

build-gui-release-windows:
	electron-builder --project ./gui -w --x64
	mv ./gui/dist/*.exe ./release-builds

build-gui-release-mac:
	electron-builder --project ./gui -m --x64
	mv ./gui/dist/*.zip ./release-builds

build-web-resources:
	cd ./gui && npm i && npm run build && cd -

clean:
	rm -rf ./dist ./gui/dist ./gui/node_modules
	rm -rf ./gui/resources/*.css ./gui/resources/*.js