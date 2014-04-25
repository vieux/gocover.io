#!/bin/sh
docker run -d -v /var/run/docker.sock:/docker.sock -p 80:8080 --name gocover.io server -r 172.17.0.3:6379 -s /docker.sock
