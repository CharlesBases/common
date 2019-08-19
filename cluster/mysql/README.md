# Mysql

- master
	192.168.1.1
- slave
	192.168.1.2

1.  创建主库
	- 添加 slave 账号
		mysql[root]> create user 'user'@'192.168.1.2' identified by '123456';
		mysql[root]> GRANT REPLICATION SLAVE ON *.* TO 'user'@'192.168.1.2';
	- 刷新新用户权限
		mysql[root]> flush privileges;
	- 查看 master 状态 (File, Position)
		mysql[root]> show master status;

2.  创建从库
	- 验证连接
		mysql -h192.168.1.1 -uuser -p123456
	- 创建用户
		- 创建用户
			mysql[root]> create user 'user'@'localhost' identified by '123456';
            mysql[root]> create user 'user'@'%' identified by '123456';
		- 授予权限
			mysql[root]> grant select ON *.* TO 'user'@'localhost';
			mysql[root]> grant select ON *.* TO 'user'@'%';
	- 设置复制
		mysql[root]> CHANGE MASTER TO MASTER_HOST='192.168.1.1',MASTER_USER='root',MASTER_PASSWORD='123456',MASTER_PORT=3306,MASTER_LOG_FILE='mysql-master-bin.000003',MASTER_LOG_POS=710,MASTER_CONNECT_RETRY=10;
	- 启动 slave
		mysql[root]> start slave;
	- 查看 slave 状态 (Slave_IO_Running[Connecting], Slave_SQL_Running[YES])
		mysql[root]> show slave status;

3. 测试同步