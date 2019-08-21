# Mysql

- master ( proxysql )
	192.168.1.1
- slave
	192.168.1.2

1.  创建主库
	- 添加 slave 账号
		mysql> create user 'user'@'192.168.1.2' identified by '123456';
		mysql> GRANT REPLICATION SLAVE ON *.* TO 'user'@'192.168.1.2';
	- 添加 proxysql 账号
		mysql> create user 'proxysql'@'192.168.1.1' identified by 'proxysql';
		mysql> GRANT ALL ON *.* TO 'proxysql'@'192.168.1.1';
	- 添加 monitor 账号
		mysql> create user 'monitor'@'192.168.1.1' identified by 'monitor';
		mysql> GRANT SELECT ON *.* TO 'monitor'@'192.168.1.1';
	- 刷新新用户权限
		mysql> flush privileges;
	- 查看 master 状态 (File, Position)
		mysql> show master status;

2.  创建从库
	- 验证连接
		mysql -h192.168.1.1 -uuser -p123456
	- 创建用户
		- 创建用户
			mysql> create user 'user'@'localhost' identified by '123456';
            mysql> create user 'user'@'%' identified by '123456';
		- 授予权限
			mysql> grant select ON *.* TO 'user'@'localhost';
			mysql> grant select ON *.* TO 'user'@'%';
	- 设置复制
		mysql> CHANGE MASTER TO MASTER_HOST='192.168.1.1',MASTER_USER='root',MASTER_PASSWORD='123456',MASTER_PORT=3306,MASTER_LOG_FILE='mysql-master-bin.000003',MASTER_LOG_POS=710,MASTER_CONNECT_RETRY=10;
	- 启动 slave
		mysql> start slave;
	- 查看 slave 状态 (Slave_IO_Running[YES], Slave_SQL_Running[YES])
		mysql> show slave status;

3. 测试同步
