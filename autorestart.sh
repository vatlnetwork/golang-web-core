#!/bin/bash

# this requires inotify-tools and a linux system in order to run
while true; do
  go build
  $@ &
  PID=$!
  echo "PID: $PID"
  inotifywait -r -e modify .
  kill $PID
  rm golang-web-core
done