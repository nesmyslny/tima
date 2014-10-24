#!/usr/bin/env bash

DEV="$( cd "$( dirname "$0" )" && pwd )"
DEV_TMP="$DEV/tmp"
BASE="$( cd "$DEV" && cd .. && pwd )"
GIN_BIN=".${DEV_TMP#$BASE}/tima-gin" # must be relative
GIN_LOG="$DEV_TMP/gin.log"
GIN_PID="$DEV_TMP/gin.pid"

mkdir -p "$DEV_TMP"

cd "$DEV"
vagrant up

(
cd "$BASE"
gin --appPort '8080' --path "$BASE" --bin "$GIN_BIN" > "$GIN_LOG" &
echo $! > "$GIN_PID"
)
