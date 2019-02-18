CD = $(shell pwd)

docker-build:
	docker build -t weatherdump_linux_amd64 ./xcompilation/linux_amd64
	docker build -t weatherdump_windows_amd64 ./xcompilation/windows_amd64

docker-run:
	docker run -v $(CD):/home/go/src/weather-dump weatherdump_linux_amd64
	docker run -v $(CD):/home/go/src/weather-dump weatherdump_windows_amd64