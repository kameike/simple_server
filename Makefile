hash := $(shell git rev-parse --verify HEAD)
.PHONY: docker_image

docker_image: Dockerfile main.go vendor
	docker build . -t simple_server:${hash}



