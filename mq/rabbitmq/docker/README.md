# 目录结构

docker                                          // 主目录
  ├── scripts                                   // 宿主机使用的一些脚本
  │       └── rabbitmq.sh                       // 配置 rabbitmq 账户和开启主从复制
  ├── volumes                                   // 各个容器的挂载数据卷
  │       ├── rabbitmq_proxy
  │       │       └── haproxy.cfg               // haproxy 配置
  │       └── rabbitmq_slave
  │               └── cluster_entrypoint.sh     // 集群入口文件
  ├── parameters.env                            // 账号密码等环境参数
  └── compose.yml                               // 编排配置
