# Go-SecKill
基于Go语言的秒杀商品系列

# 1. 单机模式

# 2. 分布式模式

## 2.1 环境搭建

1. 搭建三主三从Cluster模式的Redis集群,配置Redisson
2. 搭建zookeeper集群,导入Curator依赖


### case1:基于Redisson的Redis分布式锁，正常

接口：/seckillDistributed/handleWithRedisson?gid=1197

注意要用Redis Lock把整个事务提交都包住。这里仅仅使用了Redis分布式提供的锁功能，秒杀数据处理还是直接访问数据库来完成

### case2:基于缓存的ETCD分布式锁，正常

接口：/seckillDistributed/handleWithZk?gid=1197

类似于之前使用BlockingQueue时编写了一个单例模式的工具类来全局使用的形式相同，这里也采用静态内部类形式的单例模式编写一个Curator框架的分布式锁功能工具类ZkLockUtil来实现全局调用

注意这里也要用Zookeeper分布式锁把整个事务提交都包住。这里只用了zookeeper的分布式锁功能，秒杀数据处理也是直接访问数据库来完成

### case3:Redis的List队列，正常

接口：/seckillDistributed/handleWithRedisList?gid=1197

这里利用Redis分布式队列的方式是，在秒杀活动初始化阶段时有多少库存就在Redis的List中初始化多少个商品元素。然后每有一个用户进行秒杀，就从List队列中取出一个商品元素分配给该用户。同时将该用户信息存入到Redis的Set类型中，防止用户多次秒杀的情况。在秒杀结束之后，在Redis中数据写入到数据库中进行保存。可参考下图：


### case4:Redis原子递减,正常

接口：/seckillDistributed/handleWithRedisIncr?gid=1197

这里先将秒杀商品的库存数量，写入到redis中，利用redis的incr来实现原子递减。假如有100件商品，这里相当于准备好了100个钥匙，有人没有抢到钥匙，就返回库存不够，有人抢到了钥匙，就进行下一步处理，先将秒杀订单的信息写入到redis中，等空闲下来后在写入到数据库中。这里其实与case3差不多

### 其他
