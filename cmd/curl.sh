#!/bin/bash

# Infinite loop to continuously curl the URL and print status and time taken
while true; do
  curl -o /dev/null -s -w "Status: %{http_code}, Time: %{time_total}s\n" https://dockman.htmgo.dev
  sleep 0.5 # Optional: Pause for 1 second between requests
done
