# Mysql

- #### master ( proxysql )

	10.10.10.10

- #### slave

	10.10.10.20

## 1. 创建主库

	- 创建 slave 账号

		mysql[root]> create user 'slave'@'10.10.10.20' identified by 'slave';

		mysql[root]> GRANT REPLICATION SLAVE ON *.* TO 'slave'@'10.10.10.20';

	- 创建 proxysql 账号

		- proxysql [业务账号]

			mysql[root]> create user 'proxysql'@'%' identified by 'proxysql';

			mysql[root]> GRANT ALL ON *.* TO 'proxysql'@'%';

		- monitor [监控账号]

			mysql[root]> create user 'monitor'@'%' identified by 'monitor';

			mysql[root]> GRANT SUPER, REPLICATION CLIENT, SELECT ON *.* TO 'monitor'@'%';

	- 刷新新用户权限

		mysql[root]> flush privileges;

	- 查看 master 状态 (File, Position)

		mysql[root]> show master status;

## 2. 创建从库

	- 验证 master 连接

		mysql -h10.10.10.10 -uslave -pslave

	- 创建 slave 账号

		mysql[root]> create user 'slave'@'%' identified by 'slave';

		mysql[root]> GRANT SELECT ON *.* TO 'slave'@'%';

	- 创建 proxysql 账号

		- monitor [监控账号]

			mysql[root]> create user 'monitor'@'%' identified by 'monitor';

			mysql[root]> GRANT SUPER, REPLICATION CLIENT, SELECT ON *.* TO 'monitor'@'%';

	- 设置复制

		mysql[root]> CHANGE MASTER TO MASTER_HOST='10.10.10.10',MASTER_USER='root',MASTER_PASSWORD='123456',MASTER_PORT=3306,MASTER_LOG_FILE='mysql-master-bin.000003',MASTER_LOG_POS=888,MASTER_CONNECT_RETRY=10;

	- 启动 slave

		mysql[root]> start slave;

	- 查看 slave 状态 (Slave_IO_Running[ YES | Connecting ], Slave_SQL_Running[ YES ])

		mysql[root]> show slave status;

## 3. 测试同步
