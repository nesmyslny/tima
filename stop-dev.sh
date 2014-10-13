#!/usr/bin/env bash

vagrant halt

kill -9 $(cat "./tmp/gin.pid")
pkill -9 gnomon-gin
rm "./tmp/gin.pid"
rm "./tmp/gin.log"
rm "./tmp/gnomon-gin"
