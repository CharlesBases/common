# Mysql

- master
	192.168.1.1
- slave
	192.168.1.2

1.  创建主库
	- 添加 slave 账号
		mysql> GRANT REPLICATION SLAVE ON *.* TO 'root'@'192.168.1.81' IDENTIFIED BY '123456';
	- 查看 master 状态 (File, Position)
		mysql> show master status;
	- 刷新新用户权限
		mysql> flush privileges;

2.  创建从库
	- 验证连接
		mysql -h192.168.1.1 -uroot -p123456
	- 用户权限
		mysql> show grants for root@192.168.1.2;
	- 设置复制
		mysql> CHANGE MASTER TO MASTER_HOST='192.168.1.1',MASTER_USER='root',MASTER_PASSWORD='123456',MASTER_PORT=3306,MASTER_LOG_FILE='File',MASTER_LOG_POS=Position,MASTER_CONNECT_RETRY=10;
	- 启动 slave
		mysql> start slave
	- 查看 slave 状态 (Slave_IO_Running[YES], Slave_SQL_Running[YES])
		mysql> show slave status;

3. 测试同步