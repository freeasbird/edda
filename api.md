# Odin Api 文档


# *授权码*
### 获取授权码
`curl -X POST -d {"""""}  http://127.0.0.1:9999/asgard/api/server/license`
###### 返回:
`200 {"key": "xxxxxxxxx","date": "2019-07-23 13:36:47","msg":["xxxxx","xxxxxxx"]}`
`401 Unauthorized`

### 查找序列号
`curl -X POST http://127.0.0.1:9999/odin/api/server/code`
###### 返回:
`200 {"key": "xxxxxxxxx","date": "2019-07-23 13:36:47","msg":"xxxxx"}`
`401 Unauthorized`

### 获取序列号二维码
`curl -X GET http://127.0.0.1:9999/odin/api/server/qr-code`
###### 返回:
`二维码图片`

---

# *授权码*
### 导入授权码
`curl -X PUT http://127.0.0.1:9999/odin/api/server/license`
###### 返回:
`{"code":"xxxxxxxxx"}`

### 查看授权信息
`curl -X GET http://127.0.0.1:9999/odin/api/server/license`
###### 返回:
`{"code":"xxxxxxxxx"}`

### 删除授权
`curl -X DELETE http://127.0.0.1:9999/odin/api/server/license`
###### 返回:
`{"code":"xxxxxxxxx"}`

---

# *节点状态*
### 查看节点状态
`curl -X GET http://127.0.0.1:9999/odin/api/server/nodes`
###### 返回:
`[{"id":"a","online":"xxxxxxxxxx"},{"id":"b","online":"xxxxxxxxxx"},{"id":"c","online":"xxxxxxxxxx"}]`

---

# *配置接口*
### 新增配置
`curl -X POST http://127.0.0.1:9999/odin/api/client/conf/{ID}`
###### 返回:
`{"code":"xxxxxxxxx"}`
### 删除配置
`curl -X DELETE http://127.0.0.1:9999/odin/api/client/conf/{ID}`
###### 返回:
`{"code":"xxxxxxxxx"}`
### 修改配置
`curl -X PUT http://127.0.0.1:9999/odin/api/client/conf/{ID}`
###### 返回:
`{"code":"xxxxxxxxx"}`
### 获取配置
`curl -X GET http://127.0.0.1:9999/odin/api/client/conf/{ID}`
###### 返回:
`{"code":"xxxxxxxxx"}`

---

# *Client接口*
### 获取认证
`curl -X POST http://127.0.0.1:9999/odin/api/client/{Product}/{ID}`
###### 返回:
`{"code":200,"lease":0,"msg":"xxxxxx}`
### 心跳
`curl -X PUT http://127.0.0.1:9999/odin/api/client/{Product}/{ID}`
###### 返回:
`{"code":200,"lease":0,"msg":"xxxxx"}`
### 关闭
`curl -X DELETE http://127.0.0.1:9999/odin/api/client/{Product}/{ID}`
###### 返回:
`{"code":200,"lease":0,"msg":"xxxxx"}`

---

# *Client在线信息接口*
### 在线信息
`curl -X GET http://127.0.0.1:9999/odin/api/client/online/{Product}`
`curl -X GET http://127.0.0.1:9999/odin/api/client/online`
###### 返回:
`{"code":"xxxxxxxxx"}`






