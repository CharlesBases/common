#!/bin/bash

set -e

# MySQL
mysql_user=root
mysql_password=123456
mysql_host=127.0.0.1
mysql_port=3306
mysql_charset=utf8mb4

# backup
backup_databases=(database1 database2) 			# 备份的数据库名称，多数据库用空格分开。为空时备份所有数据库
backup_location=/opt/mysql/baks/${mysql_host}	# 备份数据存放位置
backup_expire=1 								# 开启过期数据自动清理
backup_expired_date=2 							# 过期时间天数
backup_day=`date +%F`							# 备份文件夹名称(按天分类)
backup_dir=${backup_location}/${backup_day}		# 备份文件夹全路径
backup_file=`date +%T`							# 备份文件名称
backup_gzip=0									# 是否启用压缩

# 判断 MySQL 是否启动
# mysqld_port=`netstat -tulpn | grep mysqld | wc -l`
# mysqld_process=`ps -ef | grep mysqld | wc -l`
# if [ [$mysqld_port == 0] -o [$mysqld_process == 0] ]; then
#	echo "error: mysql is not running!"
#	echo "error: backup exit!"
#	exit
# fi

# 备份 MySql
mkdir -p ${backup_dir}
if [ ${#backup_databases[*]} -ne 0 ]; then
	if [ ${backup_gzip} -ne 0 ]; then
		mysqldump --opt -u${mysql_user} -p${mysql_password} -h${mysql_host} -P${mysql_port} --quick --extended-insert --single-transaction -B ${backup_databases} | gzip > ${backup_dir}/${backup_file}.sql.gz
	else
		mysqldump --opt -u${mysql_user} -p${mysql_password} -h${mysql_host} -P${mysql_port} --quick --extended-insert --single-transaction -B ${backup_databases} > ${backup_dir}/${backup_file}.sql
	fi
else
	if [ ${backup_gzip} -ne 0 ]; then
		mysqldump --opt -u${mysql_user} -p${mysql_password} -h${mysql_host} -P${mysql_port} --quick --extended-insert --single-transaction -A | gzip > ${backup_dir}/${backup_file}.sql.gz
	else
		mysqldump --opt -u${mysql_user} -p${mysql_password} -h${mysql_host} -P${mysql_port} --quick --extended-insert --single-transaction -A > ${backup_dir}/${backup_file}.sql
	fi
fi

if [ $? -ne 0 ]; then
	printf '┌──────────────────────────────┐\n'
	printf '│    ERROR: mysqldump fail!    │\n'
	printf '└──────────────────────────────┘\n'
	exit
fi

# 删除过期备份
if [ ${backup_expire} -ne 0 ]; then
	find ${backup_location}/ -type d -mtime +$expire_days | xargs rm -rf
fi

printf '┌────────────────────────────────────┐\n'
printf '│    All database backup success!    │\n'
printf '└────────────────────────────────────┘\n'
exit
