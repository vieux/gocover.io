#!/bin/sh

docker run -d -p 127.0.0.1:6379:6379 --name redis-master crosbymichael/redis
cd worker && docker build -t worker .
cd server && go build && ./server
