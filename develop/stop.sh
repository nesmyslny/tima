#!/usr/bin/env bash

DEV="$( cd "$( dirname "$0" )" && pwd )"
DEV_TMP="$DEV/tmp"
GIN_BIN="$DEV_TMP/gnomon-gin"
GIN_LOG="$DEV_TMP/gin.log"
GIN_PID="$DEV_TMP/gin.pid"

cd "$DEV"
vagrant halt

kill -9 $(cat "$GIN_PID")
pkill -9 gnomon-gin
rm "$GIN_PID"
rm "$GIN_LOG"
rm "$GIN_BIN"
