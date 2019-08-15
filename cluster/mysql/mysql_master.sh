#!/bin/zsh

set -e

port=3306
name=mysql_master
password=123456
mysql=/Users/sun/Program/MySql  # 配置文件目录
master=master

conf=${mysql}/conf
data=${mysql}/data
logs=${mysql}/logs

mkdir -p ${conf} ${data} ${logs}

cd ${mysql}
rm -rf ${master}.cnf

# master
echo '
[mysqld]
server-id=1                         # server ID
log-bin=mysql-master-bin.log        # open logs
sync_binlog=1                       # refresh binlog frequency
log-slave-updates=1                 # synchronized update
# lower-case-table-names=1            # case sensitive
log_bin_trust_function_creators=1   # function creation permissions

# binlog-do-db=mysql                # synchronized database
# binlog-ignore-db=mysql            # ignore database

innodb_buffer_pool_size=512M
innodb_flush_log_at_trx_commit=1

character-set-server=utf8mb4
collation-server=utf8mb4_unicode_ci

sql_mode=NO_ENGINE_SUBSTITUTION,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,STRICT_TRANS_TABLES

pid-file    = /var/run/mysqld/mysqld.pid
socket      = /var/run/mysqld/mysqld.sock
datadir     = /var/lib/mysql

secure-file-priv=NULL

# Custom config should go here
!includedir /etc/mysql/conf.d/

' > ${master}.cnf

# docker rm -f $(docker ps -a | grep ${name} | awk '{print $1}')

# MySql
docker run \
-p ${port}:3306 \
-e MYSQL_ROOT_PASSWORD=${password} \
-v ${conf}:/etc/mysql/conf.d  \
-v ${logs}:/logs \
-v ${data}:/var/lib/mysql \
-v ${mysql}/${master}.cnf:/etc/mysql/my.cnf \
-d \
--name ${name} \
mysql
