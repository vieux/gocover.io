FROM golang:1.6

RUN wget https://test.docker.com/builds/Linux/x86_64/docker-1.10.3 -O /usr/bin/docker
RUN chmod +x /usr/bin/docker


COPY . /go/src/github.com/vieux/gocover.io/server
WORKDIR /go/src/github.com/vieux/gocover.io/server

RUN go get -d -v
RUN go install -v

EXPOSE 8080

ENTRYPOINT ["server"]
