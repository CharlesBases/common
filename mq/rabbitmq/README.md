# RabbitMQ Server 高可用集群

## 设计集群的目的
- 允许生产者消费者在 RabbitMQ 节点崩溃的情况下继续运行
- 通过增加更多的节点来扩展消息通信的吞吐量

## 集群配置方式
- **cluster**: 不支持跨网段，用于同一个网段内的局域网；可以随意的动态增加或者减少；节点之间需要运行相同版本的 RabbitMQ 和 Erlang。
- **federation**: 应用于广域网。允许单台服务器上的交换机或队列接收发布到另一台服务器上交换机或队列的消息，可以是单独机器或集群。federation 队列类似于单向点对点连接，消息会在联盟队列之间转发任意次，直到被消费者接受。通常使用 federation 来连接 internet 上的中间服务器，用作订阅分发消息或工作队列。
- **shovel**: 连接方式与 federation 的连接方式类似，但它工作在更低层次。可以应用于广域网。

## 节点类型
- **RAM node**: 内存节点将所有的队列、交换机、绑定、用户、权限和 vhost 的元数据定义存储在内存中，好处是可以使得像交换机和队列声明等操作更加的快速。
- **DISK node**: 将元数据存储在磁盘中，单节点系统只允许磁盘类型的节点，防止重启 RabbitMQ 的时候，丢失系统的配置信息。
```
问题说明：RabbitMQ 要求在集群中至少有一个磁盘节点，所有其他节点可以是内存节点，当节点加入或者离开集群时，必须要将该变更通知到至少一个磁盘节点。
	如果集群中唯一的一个磁盘节点崩溃的话，集群仍然可以保持运行，但是无法进行其他操作（增删改查），直到节点恢复。
解决方案：设置两个磁盘节点，至少有一个是可用的，可以保存元数据的更改。
```

## Erlang Cookie
- Erlang Cookie 是保证不同节点可以相互通信的密钥，要保证集群中的不同节点相互通信必须共享相同的 Erlang Cookie。具体的目录存放在 /var/lib/rabbitmq/.erlang.cookie。
```
RabbitMQ 底层是通过 Erlang 架构来实现的，所以 rabbitmqctl 会启动 Erlang 节点，并基于 Erlang 节点来使用 Erlang 系统连接 RabbitMQ 节点，在连接过程中需要正确的 Erlang Cookie 和节点名称，Erlang 节点通过交换 Erlang Cookie 以获得认证。
```

## 镜像队列
- **普通模式**: 默认的集群模式，以两个节点（rabbit01、rabbit02）为例来进行说明。对于 Queue 来说，消息实体只存在于其中一个节点 rabbit01（或者 rabbit02），rabbit01 和 rabbit02 两个节点仅有相同的元数据，即队列的结构。当消息进入 rabbit01 节点的 Queue 后，consumer 从 rabbit02 节点消费时，RabbitMQ 会临时在 rabbit01、rabbit02 间进行消息传输，把 A 中的消息实体取出并经过 B 发送给 consumer。所以 consumer 应尽量连接每一个节点，从中取消息。即对于同一个逻辑队列，要在多个节点建立物理 Queue。否则无论 consumer 连 rabbit01 或 rabbit02，出口总在 rabbit01，会产生瓶颈。当 rabbit01 节点故障后，rabbit02 节点无法取到 rabbit01 节点中还未消费的消息实体。如果做了消息持久化，那么得等 rabbit01 节点恢复，然后才可被消费；如果没有持久化的话，就会产生消息丢失的现象。
- **镜像模式**: 将需要消费的队列变为镜像队列，存在于多个节点，这样就可以实现 RabbitMQ 的 HA 高可用性。作用就是消息实体会主动在镜像节点之间实现同步，而不是像普通模式那样，在 consumer 消费数据时临时读取。缺点就是，集群内部的同步通讯会占用大量的网络带宽。

## 搭建 RabbitMQ Server 高可用集群




scp -rp -C /home/root/rabbitmq/data/.erlang.cookie root@192.168.1.81:/home/root/rabbitmq/data/.erlang.cookie