#!/bin/sh

DSN="postgres://${CAL_STORAGE_USER}:${CAL_STORAGE_PASSWORD}@${CAL_STORAGE_HOST}:${CAL_STORAGE_PORT}/${CAL_STORAGE_DBNAME}?sslmode=disable"

for i in $(seq 1 5); do
  echo "attempt $i: $DSN"
  ./goose postgres "$DSN" up && break
  sleep 1
done
