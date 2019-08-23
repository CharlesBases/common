# [ProxySQL](https://github.com/malongshuai/proxysql/wiki)


##1. 管理端口

	mysql -uadministrator -padministrator -h10.10.10.10 -P6032

	- ### 路由分组

		mysql> select * from main.mysql_replication_hostgroups;

	- ### 读写分离规则

		mysql> select rule_id, active,match_digest, destination_hostgroup, apply from main.mysql_query_rules;

	- ### 读写分离日志

		mysql> select * from stats.stats_mysql_query_digest;

##2. 业务端口

	mysql -uproxysql -pproxysql -h10.10.10.10 -P6033