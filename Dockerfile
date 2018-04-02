FROM golang:1.10.0
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh && \
      mkdir /var/log/app

  

WORKDIR /go/src/github.com/kameike/simple_server

COPY Gopkg.lock Gopkg.toml ./
RUN dep ensure -v -vendor-only

COPY . .

RUN go build -v

EXPOSE 3000 
CMD ["go", "run", "main.go"]

