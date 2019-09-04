#! /bin/zsh

# -------------------- docker -------------------- #
yum install -y yum-utils device-mapper-persistent-data lvm2
yum-config-manager --add-repo http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo
yum makecache fast
yum -y install docker-ce
systemctl start docker

rm -rf /etc/docker/daemon.json
echo '
{
	"registry-mirrors": ["http://uwoosppz.mirror.aliyuncs.com"]
}

' > /etc/docker/daemon.json

# -------------------- mysql -------------------- #
docker pull mysql

# -------------------- proxysql -------------------- #
docker pull $(docker search proxysql | awk 'NR==2' | awk '{print $1}')