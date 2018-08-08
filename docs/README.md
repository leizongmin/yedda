# 设计文档

## 通讯协议

### 数据包结构

包长度 = 10 + data.length

| 字段名 | 类型 | 长度 | 说明 |
| --- | --- | --- | --- |
| version | uint16 | 2 | 协议版本，当前为 `1` |
| id | uint32 | 4 | 数据包序号 |
| op | uint16 | 2 | 当前操作类型 |
| length | uint16 | 2 | 操作数据长度 |
| data | buffer | 可变 | 操作数据内容 |


| 操作类型 | 值 | 说明 |
| --- | --- | --- |
| PING | 0x1 | 发送PING |
| PONG | 0x2 | 回应PING |
| GET | 0x3 | 查询指定Key的当前计数 |
| GET_RESULT | 0x4 | GET指令的结果 |
| INCR | 0x5 | 增加指定Key的技术 |
| INCR_RESULT | 0x6 | INCR指令的结果 |


### PING / PONG 操作数据结构

包长度 = 64

| 字段名 | 类型 | 长度 | 说明 |
| --- | --- | --- | --- |
| time | uint32 | 8 | 当前毫秒时间戳 |

### GET / INCR 操作数据结构

包长度 = 14 + ns.length + key.length

| 字段名 | 类型 | 长度 | 说明 |
| --- | --- | --- | --- |
| db | uint32 | 4 | 数据库号 |
| nsLength | uint8 | 1 | 命名空间名称长度，最长255 |
| ns | buffer | 可变 | 命名空间名称 |
| milliseconds | uint32 | 4 | 数据有效期毫秒数 |
| keyLength | uint8 | 1 | 键长度 |
| key | buffer | 可变 | 键内容，最长255 |
| count | uint32 | 4 | 增加的数量 |

