#!/bin/sh

docker run -d -p 127.0.0.1:6379:6379 --name redis-master crosbymichael/redis
cd workers
docker build -t worker:1.2.2 1.2.2/
docker build -t worker:1.3rc1 1.3rc1/
docker tag worker:1.2.2 worker:latest
cd ../server
go build && ./server
