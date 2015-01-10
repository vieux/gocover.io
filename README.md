[GoCover.io](http://gocover.io)
===============================


What is it ?
------------

[GoCover.io](http://gocover.io) offers the code coverage of any golang package as a service.

The sources are available on github.

How does it works ?
-------------------

Getting code coverage of a sofware requires running it's tests.

As executing all that code (running the tests) could be dangerous, it has to be done in a isolated environment.

That's why, each time you get the coverage of a package, the tests are run inside a [docker](http://docker.io) container.

Development
-----------

### Linux 

```console
docker run -d -p 127.0.0.1:6379:6379 --name redis-master crosbymichael/redis
cd server && go build && ./server
```

### Mac

* Install [Boot2Docker](https://github.com/boot2docker/boot2docker)

```console
boot2docker up
$(boot2docker shellinit)
docker run -d -p 6379:6379 --name redis-master crosbymichael/redis

cd server && go build
./server -r=$(boot2docker ip 2>/dev/null):6379
```

