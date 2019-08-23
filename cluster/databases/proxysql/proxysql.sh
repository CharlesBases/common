#!/bin/bash

admin_port=6032
mysql_port=6033

name=proxysql                                       # container name
proxysql=/home/root/MySql/proxysql                  # proxysql dir

data=${proxysql}/data
conf=${proxysql}/proxysql.cnf

rm -rf ${proxysql}
mkdir -p ${data}

echo '
datadir="/var/lib/proxysql"                         # 数据目录

admin_variables =
{
	admin_credentials = "administrator:administrator" # admin 凭证
	mysql_ifaces = "0.0.0.0:6032"                     # admin 管理端口
	refresh_interval = 2000
	# debug = true
}

mysql_variables =
{
	threads = 4                                       # 转发端口的线程数量，建议与CPU核心数相等
	max_connections = 2048
	default_query_delay = 0
	default_query_timeout = 36000000
	have_compress = true
	poll_timeout = 2000
	interfaces = "0.0.0.0:6033;/tmp/proxysql.sock"    # mysql 代理端口
	default_schema = "information_schema"
	stacksize = 1048576
	server_version = "8.0.17"                         # mysql 版本
	connect_timeout_server = 10000
	ping_interval_server_msec = 10000
	ping_timeout_server = 200
	commands_stats = true
	sessions_sort = true
    monitor_username = "monitor"
    monitor_password = "monitor"
    monitor_history = 60000
    monitor_connect_interval = 60000
    monitor_connect_timeout = 3000
	monitor_ping_max_failures = 3
    monitor_ping_interval = 5000
    monitor_ping_timeout = 3000
    monitor_read_only_interval = 1000
    monitor_read_only_timeout = 500
}

mysql_servers =
(
    {
		hostgroup = 10              # master group
	    address = "10.10.10.10"
	    port = 3306
		weight = 1
	    status = "ONLINE"
	    max_connections = 1000      # 最大连接
    },
    {
		hostgroup = 20              # slave group
		address = "10.10.10.20"
		port = 3306
		weight = 1
		status = "ONLINE"
	    max_connections = 1000      # 最大连接
	    max_replication_lag = 10    # 最大延迟 (只适用于从节点)
    },
    {
		hostgroup = 20
	    address = "10.10.10.10"
	    port = 3306
		weight = 10
	    status = "ONLINE"
	    max_connections = 1000
    }
)

mysql_users =
(
    {
        username = "proxysql"
        password = "proxysql"
        default_hostgroup = 10
        max_connections = 1000
        active = 1
    },
    {
        username = "monitor"
        password = "monitor"
        default_hostgroup = 20
        max_connections = 1000
        active = 1
    }
)

mysql_query_rules =
(
    {
        rule_id = 10
        active = 1
        match_pattern = "^SELECT .* FOR UPDATE"
        destination_hostgroup = 10
        apply = 1
    },
    {
        rule_id = 20
        active = 1
        match_pattern = "^SELECT .*"
        destination_hostgroup = 20
        apply = 1
        # cache_ttl = 10        #
    }
)

mysql_replication_hostgroups =
(
    {
        writer_hostgroup = 10
        reader_hostgroup = 20
        comment = "MySql Ver 8.0.17"
   }
)

#	scheduler =
#	(
#		{
#			id = 1
#			active = 0
#			interval_ms = 3600
#			filename = "/home/root/MySql/proxysql/proxysql_timer.sh"
#			arg1 = "0"
#			arg2 = "0"
#			arg3 = "0"
#			arg4 = "1"
#			comment = "this is a timer task."
#		}
#	)

' > ${proxysql}/proxysql.cnf

container_id=$(docker ps -a | grep ${name} | awk '{print $1}')
if [ ${#container_id[@]} -gt 0 ]
then
	docker rm -f ${container_id}
fi

# proxysql
docker run \
	-p ${admin_port}:6032 -p ${mysql_port}:6033 \
	-v ${data}:/var/lib/proxysql/ \
	-v ${conf}:/etc/proxysql.cnf \
	-d \
	--name=${name} \
	proxysql/proxysql

# docker logs -f $(docker ps -l -q)
