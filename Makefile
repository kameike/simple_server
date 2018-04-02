hash := $(shell git rev-parse --verify HEAD)
.PHONY: docker_image

run: $(wildcard *.go)
	go run -v main.go
 
docker-image: Dockerfile main.go vendor $(wildcard *.go)
	docker build . -t kameike/simple_server 

