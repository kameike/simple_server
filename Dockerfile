FROM golang:1.10.0 AS build-env

# RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN go get -u github.com/golang/dep/cmd/dep

WORKDIR /go/src/github.com/kameike/simple_server
COPY Gopkg.lock Gopkg.toml ./
RUN dep ensure -v -vendor-only
COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -v -o main

FROM alpine
COPY --from=build-env /go/src/github.com/kameike/simple_server/main main
EXPOSE 3000 
CMD ["./main"]

