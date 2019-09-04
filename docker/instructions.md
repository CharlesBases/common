# images
- 删除所有镜像
	docker rmi -f $(docker images -q)

- 删除指定镜像
	docker rmi -f $(docker images -a | grep 'images' | awk '{print $3}')

- 删除没有使用的镜像
	docker rmi -f $(docker images -a | grep '<none>' | awk '{print $3}')

# container
- 删除所有容器
	docker rm -f $(docker ps -a -q)

- 删除指定容器
	docker rm -f $(docker ps -a -q | grep 'mysql' | awk '{print $1}')

- 杀死所有正在运行的容器
	docker kill $(docker ps -a -q)

- 删除所有已停止的容器
	docker rm -f $(docker ps -a -q)
	docker rm -f $(docker ps --all -q -f status=exited)