#!/bin/zsh

set -e

port=3306
name=mysql_slave
password=123456
mysql=/Users/sun/Program/MySql  # 配置文件目录
slave=slave

conf=${mysql}/conf
data=${mysql}/data
logs=${mysql}/logs

rm -rf ${mysql}
mkdir -p ${conf} ${data} ${logs}

# slave
echo '
[mysqld]
server-id                       = 2
read_only                       = 1
relay-log                       = mysql-relay-bin
log-slave-updates               = 1
key_buffer_size                 = 16M
max_allowed_packet              = 16M
thread_stack                    = 256K
thread_cache_size               = 8
symbolic-links                  = 0
skip_name_resolve               = ON

character-set-server            = utf8mb4
collation-server                = utf8mb4_unicode_ci

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

' > ${mysql}/${slave}.cnf

# docker rm -f $(docker ps -a | grep ${name} | awk '{print $1}')

# MySql
docker run \
-p ${port}:3306 \
-e MYSQL_ROOT_PASSWORD=${password} \
-v ${conf}:/etc/mysql/conf.d  \
-v ${logs}:/logs/mysql/ \
-v ${data}:/var/lib/mysql \
-v ${mysql}/${slave}.cnf:/etc/mysql/my.cnf \
-d \
--name ${name} \
mysql
