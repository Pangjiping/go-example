# 发布订阅模式

使用发布订阅模式实现数据库迁移和数据同步

## 发布者

- 发布者设置为从旧的数据库中获取数据并且format
- 两级发布者，先从本地clusterInfo表中获取数据并且format，再从tag service获取，
- 两级发布者使用channel通信，限定channel大小

## 订阅者

- 订阅者只负责向新的tag table中插入数据?