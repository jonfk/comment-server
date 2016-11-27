#!/usr/bin/env sh

set -x
set -e

DB_CONTAINER=cs-postgres
DB_VOLUME="/opt/comment-server"

source src/github.com/jonfk/comment-server/.env

docker run --name "$DB_CONTAINER" \
       -e POSTGRES_PASSWORD="$DATABASE_PASSWORD" \
       -e POSTGRES_DB="$DATABASE_NAME" \
       -e POSTGRES_USER="$DATABASE_USER" \
       -p "$DATABASE_PORT":5432 -v `pwd`:"$DB_VOLUME" -d postgres:latest

sleep 5s
docker exec "$DB_CONTAINER" psql -U comment-server -d comment-server -a -f "$DB_VOLUME"/migrations/up.sql
