#!/bin/bash
#cd workerclient
#go build -o workerclient main.go && ./workerclient

#!/bin/bash

export COMPOSE_TLS_VERSION=TLSv1_2

if [ $# -lt 1 ]; then
  echo "Insufficient parameters action up or down required";exit 1
else
if [ $# -gt 1 ]; then
  echo "Warning only first parameter required"
fi
case $1 in
  start) echo "Building and Starting Client"
      docker-compose -f docker-compose-client.yaml up --force-recreate --build;;
  stop) echo "Stopping Client"
      docker-compose -f docker-compose-client.yaml down --remove-orphans;;
  *) echo "Unknown parameter"
     exit 1;;
esac

fi