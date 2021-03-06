#!/bin/bash
export COMPOSE_TLS_VERSION=TLSv1_2
if [ $# -lt 1 ]; then
  echo "Insufficient parameters action up or down required";exit 1
else
if [ $# -gt 1 ]; then
  echo "Warning only first parameter required"
fi
case $1 in
  start) echo "Building and Starting Server"
      docker-compose -f docker-compose-server.yaml up --force-recreate --build;;
  stop) echo "Stopping server"
      docker-compose -f docker-compose-server.yaml down --remove-orphans;;
  *) echo "Unknown parameter"
     exit 1;;
esac

fi