#!/bin/bash

set -e

# Docker
port=3306
name=mysql_master

# MySQL
mysql_root_password=123456

# config
mysql_dir=/home/root/mysql
master_tag=master

baks=${mysql_dir}/baks
conf=${mysql_dir}/conf
data=${mysql_dir}/data
logs=${mysql_dir}/logs

rm -rf ${mysql_dir}
mkdir -p ${baks} ${conf} ${data} ${logs}

# master
echo '
[mysqld]
server-id                       = 1
log-bin                         = mysql-master-bin.log
sync_binlog                     = 1
log_bin_trust_function_creators = 1
key_buffer_size                 = 16M
max_allowed_packet              = 16M
thread_stack                    = 256K
thread_cache_size               = 8
symbolic-links                  = 0
skip_name_resolve               = ON

character-set-server            = utf8mb4
collation-server                = utf8mb4_unicode_ci
default_authentication_plugin   = mysql_native_password

innodb_file_per_table           = ON
innodb_buffer_pool_size         = 512M
innodb_flush_log_at_trx_commit  = 1

sql_mode                        = NO_ENGINE_SUBSTITUTION,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,STRICT_TRANS_TABLES

pid-file    = /var/run/mysqld/mysqld.pid
socket      = /var/run/mysqld/mysqld.sock
datadir     = /var/lib/mysql

secure-file-priv = NULL

[mysqld_safe]
!includedir /etc/mysql/conf.d/
log-error = /logs/mysql/server.log

' > ${mysql_dir}/${master_tag}.cnf

container_id=$(docker ps -a | grep ${name} | awk '{print $1}')
if [ ${#container_id} -ne 0 ]
then
	docker rm -f ${container_id}
fi

# MySQL
docker run \
	-p ${port}:3306 \
	-e TZ="Asia/Shanghai" \
	-e MYSQL_ROOT_PASSWORD=${mysql_root_password} \
	-v ${baks}:/opt/mysql/baks  \
	-v ${conf}:/etc/mysql/conf.d  \
	-v ${logs}:/logs/mysql \
	-v ${data}:/var/lib/mysql \
	-v ${mysql_dir}/${master_tag}.cnf:/etc/mysql/my.cnf \
	-d \
	--name ${name} \
	mysql
