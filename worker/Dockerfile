FROM crosbymichael/golang

RUN echo "deb http://security.debian.org testing/updates main" >> /etc/sources.list

RUN apt-get update
RUN apt-get install subversion bzr build-essential --no-install-recommends -y

# Custom libs for some projects
RUN apt-get install libreadline-dev -y

RUN go get code.google.com/p/go.tools/cmd/cover

ADD gocover.sh /
RUN chmod +x /gocover.sh

WORKDIR /go/src

ENTRYPOINT ["/gocover.sh"]
