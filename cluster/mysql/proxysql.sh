#!/bin/bash

admin_port=6032
mysql_port=6033

container_name=proxysql

mysql_root_password=123456
mysql_proxy_user=proxysql
mysql_proxysql_password=123456

cluster_name=proxysql_cluster
etcd_host=10.20.2.4:2379
network=proxysql_net

proxysql=/Users/sun/Program/MySql/proxysql          # proxysql dir
data=${proxysql}/data
conf=${proxysql}/proxysql.cnf

rm -rf ${proxysql}
mkdir -p ${data}

echo '
datadir="/var/lib/proxysql"                         # 数据目录

admin_variables =
{
	admin_credentials="proxysql:proxysql"           # admin 凭证
	mysql_ifaces="0.0.0.0:6032"                     # admin 管理端口
	refresh_interval=2000
	# debug=true
}

mysql_variables =
{
	threads=4                                       # 转发端口的线程数量，建议与CPU核心数相等
	max_connections=2048
	default_query_delay=0
	default_query_timeout=36000000
	have_compress=true
	poll_timeout=2000
	interfaces="0.0.0.0:6033;/tmp/proxysql.sock"    # mysql 代理端口
	default_schema="information_schema"
	stacksize=1048576
	server_version="8.0.17"                         # mysql 版本
	connect_timeout_server=10000
	ping_interval_server_msec=10000
	ping_timeout_server=200
	commands_stats=true
	sessions_sort=true
    monitor_username="monitor"
    monitor_password="monitor"
    monitor_history=600000
    monitor_connect_interval=60000
    monitor_ping_interval=10000
    monitor_read_only_interval=1500
    monitor_read_only_timeout=500
}

mysql_servers =
(
    {
		hostgroup = 1               # master group
	    address = "192.168.1.99"
	    port = 3306
	    status = "ONLINE"
	    weight = 1
	    compression = 0
	    max_connections = 200
    },
    {
		hostgroup = 2               # slave group
		address = "192.168.1.80"
		port = 3306
		status = "ONLINE"
		weight = 1
		compression = 0
		max_connections=1000
    }
)

mysql_users =
(
    {
        username = "proxysql"
        password = "proxysql"
        default_hostgroup = 1
        max_connections=1000
        active = 1
    },
    {
        username = "user"
        password = "123456"
        default_hostgroup = 2
        max_connections=1000
        active = 1
    }
)

mysql_query_rules =
(
#	{
#		rule_id=1
#		active=1
#		match_pattern="^SELECT .* FOR UPDATE"
#		destination_hostgroup=1
#		apply=1
#	},
#	{
#		rule_id=2
#		active=1
#		match_pattern="^SELECT .*"
#		destination_hostgroup=2
#		apply=1
#	},
#	{
#		rule_id=3
#		active=1
#		match_pattern=".*"
#		destination_hostgroup=1
#		apply=1
#	}
)

mysql_replication_hostgroups=
(
    {
        writer_hostgroup = 1
        reader_hostgroup = 2
        comment = "MySql Ver 8.0.17"
   }
)

scheduler =
(
#	{
#		id=1
#		active=0
#		interval_ms=10000
#		filename="/var/lib/proxysql/proxysql_galera_checker.sh"
#		arg1="0"
#		arg2="0"
#		arg3="0"
#		arg4="1"
#		arg5="/var/lib/proxysql/proxysql_galera_checker.log"
#	}
)


' > ${proxysql}/proxysql.cnf

docker rm -f $(docker ps -a | grep ${container_name} | awk '{print $1}')

# proxysql
docker run \
	-p ${admin_port}:6032 -p ${mysql_port}:6033 \
	-v ${data}:/var/lib/proxysql/ \
	-v ${conf}:/etc/proxysql.cnf \
	-d \
	--name=${container_name} \
	proxysql/proxysql

# docker logs -f $(docker ps -l -q)