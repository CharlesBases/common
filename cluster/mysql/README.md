# Mysql

- master
	192.168.1.1
- slave
	192.168.1.2

1.  创建主库
	- 添加 slave 账号
		mysql> CREATE USER 'user'@'192.168.1.2' IDENTIFIED BY '123456';
		mysql> GRANT REPLICATION SLAVE ON *.* TO 'user'@'192.168.1.2';
	- 查看 master 状态 (File, Position)
		mysql> show master status;
	- 刷新新用户权限
		mysql> flush privileges;

2.  创建从库
	- 验证连接
		mysql -h192.168.1.80 -uuser -p123456
	- 创建用户
		- 创建用户
			mysql> CREATE USER 'user'@'localhost' IDENTIFIED BY '123456';
            mysql> CREATE USER 'user'@'%' IDENTIFIED BY '123456';
		- 授予权限
			mysql> GRANT SELECT ON *.* TO 'user'@'localhost';
			mysql> GRANT SELECT ON *.* TO 'user'@'%';
	- 用户权限
		mysql> SHOW GRANTS for user@192.168.1.2;
	- 设置复制
		mysql> CHANGE MASTER TO MASTER_HOST='192.168.1.80',MASTER_USER='root',MASTER_PASSWORD='123456',MASTER_PORT=3306,MASTER_LOG_FILE='mysql-master-bin.000003',MASTER_LOG_POS=710,MASTER_CONNECT_RETRY=10;
	- 启动 slave
		mysql> start slave
	- 查看 slave 状态 (Slave_IO_Running[YES], Slave_SQL_Running[YES])
		mysql> show slave status;

3. 测试同步