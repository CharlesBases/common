# 加速器
- Linux
```
[root@centos 7]# touch /etc/docker/daemon.json
[root@centos 7]# vi /etc/docker/daemon.json
	{
        "registry-mirrors": ["http://uwoosppz.mirror.aliyuncs.com"]
	}
```

# 交叉编译
- Mac

系统      | 指令
--       | --
Linux    | CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build *.go
Windows  | CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build *.go

- Linux

系统      | 指令
--       | --
Mac      | CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build *.go
Windows  | CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build *.go
