FROM golang:1.5.3


COPY . /go/src/github.com/vieux/gocover.io/server
WORKDIR /go/src/github.com/vieux/gocover.io/server

RUN go get -d -v
RUN go install -v

EXPOSE 8080

ENTRYPOINT ["server"]
