hash := $(shell git rev-parse --verify HEAD)
.PHONY: docker_image


run: $(wildcard *.go)
	go run main.go
 
docker_image: Dockerfile main.go vendor $(wildcard *.go)
	docker build . -t kameike/simple_server && docker push kameike/simple_server

