#!/bin/bash

set -e

# Docker
name=rabbitmq
admin_port=15672                # 控制台端口
visit_port=5672                 # 业务端口
RABBITMQ_DEFAULT_USER=admin
RABBITMQ_DEFAULT_PASS=admin

# config
rabbitmq_dir=/Users/sun/Program/rabbitmq

conf=${rabbitmq_dir}/conf
data=${rabbitmq_dir}/data
logs=${rabbitmq_dir}/logs

rm -rf ${rabbitmq_dir}
mkdir -p ${conf} ${data} ${logs}

echo '
[rabbitmq_management].
' > ${conf}/enabled_plugins

container_id=$(docker ps -a | grep ${name} | awk '{print $1}')
if [ ${#container_id} -ne 0 ]
then
	docker rm -f ${container_id}
fi

# RabbitMQ
docker run \
	-p ${admin_port}:15672 -p ${visit_port}:5672 \
	-e RABBITMQ_DEFAULT_USER=${RABBITMQ_DEFAULT_USER} \
	-e RABBITMQ_DEFAULT_PASS=${RABBITMQ_DEFAULT_PASS} \
	-v ${conf}:/etc/rabbitmq  \
    -v ${logs}:/var/log/rabbitmq \
    -v ${data}:/var/lib/rabbitmq \
	-d \
	--name=${name} \
	rabbitmq:3-management