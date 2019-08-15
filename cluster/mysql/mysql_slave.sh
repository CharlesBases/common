#!/bin/zsh

set -e

port=3307
name=mysql_slave
password=123456
mysql=/Users/sun/Program/MySql  # 配置文件目录
slave=slave

conf=${mysql}/conf
data=${mysql}/data
logs=${mysql}/logs

mkdir -p ${conf} ${data} ${logs}

cd ${mysql}
rm -rf ${slave}.cnf

# slave
echo '
[mysqld]
server-id=2                         # server ID
log-bin=mysql-slave-bin.log         # open logs
sync_binlog=1                       # refresh binlog frequency
relay-log=mysql-relay-bin           #
read-only=1                         # read only
log-slave-updates=1                 # synchronized update
# lower-case-table-names=1           # case sensitive
log_bin_trust_function_creators=1   # function creation permissions

# binlog-do-db=repl                 # synchronized database
# binlog-ignore-db=mysql            # ignore database

innodb_buffer_pool_size=512M
innodb_flush_log_at_trx_commit=1

character-set-server=utf8mb4
collation-server=utf8mb4_unicode_ci

sql_mode=NO_ENGINE_SUBSTITUTION,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,STRICT_TRANS_TABLES

pid-file    = /var/run/mysqld/mysqld.pid
socket      = /var/run/mysqld/mysqld.sock
datadir     = /var/lib/mysql

secure-file-priv= NULL

symbolic-links=0

# Custom config should go here
!includedir /etc/mysql/conf.d/

' > ${slave}.cnf

# MySql
docker run \
-p ${port}:3306 \
-e MYSQL_ROOT_PASSWORD=${password} \
-v ${mysql}/${slave}.cnf:/etc/mysql/my.cnf \
-d \
--name ${name} \
mysql