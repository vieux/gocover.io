FROM golang:1.8

RUN wget https://download.docker.com/linux/static/stable/x86_64/docker-19.03.8.tgz ; tar -xvf docker-19.03.8.tgz ; cp docker/docker /usr/bin/docker ; rm -rf docker


COPY . /go/src/github.com/vieux/gocover.io/server
WORKDIR /go/src/github.com/vieux/gocover.io/server

RUN go get -d -v
RUN go install -v

EXPOSE 8080

ENTRYPOINT ["server"]
