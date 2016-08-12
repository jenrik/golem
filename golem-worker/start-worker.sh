#!/bin/sh
# For starting golem-worker outside docker while still communicating containerized services
set -e

export PGSSLMODE=disable
export GOLEM_STORAGE=postgresql://golem:golem@127.0.0.1/golem
export GOLEM_QUEUE=amqp://127.0.0.1/
export GOLEM_WORKERS=4

cd "${0%/*}"
go build
echo "> Finished compiling"
./golem-worker
