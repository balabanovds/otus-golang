#!/bin/sh

PROJ_NAME=tests
DC_DIR=./deployments
DC_FILE_COMMON=$DC_DIR/docker-compose.yaml
DC_FILE_TESTS=$DC_DIR/docker-compose.tests.yml
DOCKER_NETWORK=cal_bridge
DOCKER_VOLUME=${PROJ_NAME}_dbdata

echo Creating docker network $DOCKER_NETWORK
docker network create $DOCKER_NETWORK
docker-compose -f $DC_FILE_COMMON -p $PROJ_NAME up -d --build
sleep 10
docker-compose -f $DC_FILE_TESTS -p $PROJ_NAME up --build
rc=$?
docker-compose -f $DC_FILE_COMMON -p $PROJ_NAME down
echo Removing docker network $DOCKER_NETWORK
docker network rm $DOCKER_NETWORK
echo Removing docker volume $DOCKER_VOLUME
docker volume rm $DOCKER_VOLUME
exit $rc
