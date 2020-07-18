#!/bin/bash
if [ $# -lt 1 ]; then
  echo "Insufficient parameters action up or down required";exit 1
else
if [ $# -gt 1 ]; then
  echo "warning only first parameter required"
fi
case $1 in
  start) echo "Building and Starting Server"
      docker-compose up --force-recreate;;
  stop) echo "Stopping server"
        docker-compose down;;
  *) echo "Unknown parameter"
     exit 1;;
esac

fi