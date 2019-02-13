CD = $(shell pwd)

docker-build:
	docker build -t compiler ./build

docker-run:
	docker run -v $(CD):/home/go/src/weather-dump compiler