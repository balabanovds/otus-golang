#!/bin/sh

PROJ_NAME=calendar
DC_DIR=./deployments
DC_FILE_COMMON=$DC_DIR/docker-compose.yaml
DC_FILE_PROD=$DC_DIR/docker-compose.prod.yml
DOCKER_NETWORK=cal_bridge
DOCKER_VOLUME=${PROJ_NAME}_dbdata

case "$1" in
up)
  echo Creating external network $DOCKER_NETWORK
  docker network create $DOCKER_NETWORK
  docker-compose -f $DC_FILE_COMMON -f $DC_FILE_PROD -p $PROJ_NAME up -d --build
  ;;
down)
  docker-compose -f $DC_FILE_COMMON -f $DC_FILE_PROD -p $PROJ_NAME down
  echo Removing docker network $DOCKER_NETWORK
  docker network rm $DOCKER_NETWORK
  echo Removing docker volume $DOCKER_VOLUME
  docker volume rm $DOCKER_VOLUME
  ;;
*)
  echo Usage ./deploy.sh up|down
  ;;
esac
