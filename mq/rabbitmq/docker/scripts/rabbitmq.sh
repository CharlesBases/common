#!/bin/bash

set -e

docker exec -it rmqha_node0 rabbitmqctl set_policy ha-all '^' '{"ha-mode":"all"}'
