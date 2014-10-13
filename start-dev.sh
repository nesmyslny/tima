#!/usr/bin/env bash

vagrant up

mkdir -p "./tmp"

(
gin -a '8080' -b "./tmp/gnomon-gin" > "./tmp/gin.log" &
echo $! > "./tmp/gin.pid"
)
