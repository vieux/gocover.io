FROM golang:1.8

RUN wget https://get.docker.com/builds/Linux/x86_64/docker-1.13.1.tgz ; tar -xvf docker-1.13.1.tgz ; cp docker/docker /usr/bin/docker ; rm -rf docker


COPY . /go/src/github.com/vieux/gocover.io/server
WORKDIR /go/src/github.com/vieux/gocover.io/server

RUN go get -d -v
RUN go install -v

EXPOSE 8080

ENTRYPOINT ["server"]
