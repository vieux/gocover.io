#!/bin/sh
go build -a -tags "netgo static_build"
docker build -t server .
