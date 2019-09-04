# 禁用防火墙
```
[root@centos 7]# systemctl disable firewalld.service
[root@centos 7]# reboot
```

# 修改 IP
```
查看当前网卡名称
[root@centos 7]# ifconfig
	eth0: flags=4163<UP,BROADCAST,RUNNING,MULTICAST>  mtu 1500
            inet 192.168.1.1  netmask 255.255.255.0  broadcast 10.211.55.255
            inet6 fdb2:2c26:f4e4:0:d5a2:332c:f355:9ca0  prefixlen 64  scopeid 0x0<gl
修改网卡配置文件
[root@centos 7]# vi /etc/sysconfig/network-scripts/ifcfg-eth0
    TYPE=Ethernet
	DEFROUTE=yes
	PEERDNS=yes
	PEERROUTES=yes
	BOOTPROTO=static			# 使用静态IP地址, 默认为dhcp
	IPADDR=192.168.1.88			# 设置的静态IP地址
	NETMASK=255.255.255.0		# 子网掩码
	GATEWAY=192.168.1.1			# 网关地址
	DNS1=223.5.5.5				# DNS服务器 1
	DNS2=223.6.6.6				# DNS服务器 2
	NM_CONTROLLED=no 			# 通过配置文件配置
	DEFROUTE=yes
	IPV4_FAILURE_FATAL=no
	IPV6INIT=yes
	IPV6_AUTOCONF=yes
	IPV6_DEFROUTE=yes
	IPV6_PEERROUTES=yes
	IPV6_FAILURE_FATAL=no
	IPV6_ADDR_GEN_MODE=stable-privacy
	NAME=eth0
	UUID=9410076c-bee9-4886-af6c-537c17bcfee0
	DEVICE=eth0
	ONBOOT=yes					# 是否开机启用
重启网络
[root@centos 7]# service network restart
```