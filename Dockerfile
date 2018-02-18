FROM golang:1.10.0-stretch

WORKDIR /go/src/github.com/kameike/simple_server
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 3000 

CMD ["go", "run", "main.go"]

