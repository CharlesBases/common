#!/bin/bash

# Docker
name=proxysql
admin_port=6032     # ProxySQL 管理端口
mysql_port=6033     # ProxySQL 业务端口( 用于连接 MySQL )

# ProxySQL
proxysql_admin=administrator                    # ProxySQL admin 凭证
proxysql_admin_password=administrator
proxysql_writer_server=192.168.1.80              # MySQL master 服务器
proxysql_writer_server_port=3306
proxysql_reader_server=192.168.1.81              # MySQL slave 服务器
proxysql_reader_server_port=3306
proxysql_master_writer=proxysql                 # MySQL master 业务账户(用于写数据，以及 slave 崩溃时读数据)
proxysql_master_writer_password=proxysql
proxysql_master_monitor=monitor                 # MySQL master 监控账户
proxysql_master_monitor_password=monitor
proxysql_slave_reader=slave                     # MySQL slave 业务账户
proxysql_slave_reader_password=slave

# MySQL
mysql_version=5.7.27

# config
proxysql_dir=/home/root/mysql/proxysql

data=${proxysql_dir}/data
conf=${proxysql_dir}/proxysql.cnf

rm -rf ${proxysql_dir}
mkdir -p ${data}

echo '
datadir="/var/lib/proxysql"

admin_variables =
{
	admin_credentials = "'${proxysql_admin}':'${proxysql_admin_password}'"
	mysql_ifaces = "0.0.0.0:6032"
	refresh_interval = 2000
	# debug = true
}

mysql_variables =
{
	threads = 4
	max_connections = 2048
	default_query_delay = 0
	default_query_timeout = 36000000
	have_compress = true
	poll_timeout = 2000
	interfaces = "0.0.0.0:6033;/tmp/proxysql.sock"
	default_schema = "information_schema"
	stacksize = 1048576
	server_version = "'${mysql_version}'"
	connect_timeout_server = 10000
	ping_interval_server_msec = 10000
	ping_timeout_server = 200
	commands_stats = true
	sessions_sort = true
    monitor_username = "'${proxysql_master_monitor}'"
    monitor_password = "'${proxysql_master_monitor_password}'"
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
		hostgroup = 10
	    address = "'${proxysql_writer_server}'"
	    port = '${proxysql_writer_server_port}'
		weight = 1
	    status = "ONLINE"
	    max_connections = 1000
    },
    {
		hostgroup = 20
		address = "'${proxysql_reader_server}'"
		port = '${proxysql_reader_server_port}'
		weight = 1
		status = "ONLINE"
	    max_connections = 1000
	    max_replication_lag = 10
    },
    {
		hostgroup = 20
	    address = "'${proxysql_writer_server}'"
	    port = '${proxysql_writer_server_port}'
		weight = 10
	    status = "ONLINE"
	    max_connections = 1000
    }
)

mysql_users =
(
    {
        username = "'${proxysql_master_writer}'"
        password = "'${proxysql_master_writer_password}'"
        default_hostgroup = 10
        max_connections = 1000
        active = 1
    },
    {
        username = "'${proxysql_slave_reader}'"
        password = "'${proxysql_slave_reader_password}'"
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
        comment = "MySQL Ver 8.0.17"
   }
)

#	scheduler =
#	(
#		{
#			id = 1
#			active = 0
#			interval_ms = 3600
#			filename = "/home/root/mysql/proxysql/proxysql_timer.sh"
#			arg1 = "0"
#			arg2 = "0"
#			arg3 = "0"
#			arg4 = "1"
#			comment = "this is a timer task."
#		}
#	)

' > ${proxysql_dir}/proxysql.cnf

container_id=$(docker ps -a | grep ${name} | awk '{print $1}')
if [ ${#container_id} -ne 0 ]
then
	docker rm -f ${container_id}
fi

# ProxySQL
docker run \
	-p ${admin_port}:6032 -p ${mysql_port}:6033 \
	-e TZ="Asia/Shanghai" \
	-v ${data}:/var/lib/proxysql/ \
	-v ${conf}:/etc/proxysql.cnf \
	-d \
	--name=${name} \
	proxysql/proxysql

# docker logs -f $(docker ps -l -q)