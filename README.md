# edda #

----

## what this? ##
- edda 用于给 [odin](https://github.com/offer365/odin) 生成授权码。 
- 使用mongodb存储数据。

## 安装运行 ##
#### 安装edda

```
unzip edda-v1.0-linux-amd64.zip
cd edda
sh install.sh
/home/admin/mongodb/bin/mongo
# 非auth 模式下创建用户
use admin
db.createUser({user:"admin",pwd:"eddaedda",roles:["root"]})
use edda
db.createUser({user:"edda",pwd:"edda",roles:[{role:"dbOwner",db:"edda"}]})
exit
# 配置文件添加 auth=true 重启mongodb
echo "auth=true" >> /home/admin/mongodb/conf/mongodb.conf
systemctl restart mongodb
# use edda
# db.auth("edda","edda") # 返回1

systemctl start edda
```

> 访问 http://127.0.0.1:1999


#### 相关说明
> 配置文件是 edda.json 
>
> 修改 edda.service 可以指定程序与配置文件的位置
>

## 使用介绍 ##
1. 先安装 edda 并运行。访问web端口，默认账号密码：admin:123 可在配置文件 edda.json 中修改。
2. 使用 web 或访问 api 接口生成 序列号。
3. 在 [edda](https://github.com/offer365/edda) 里根据约定新建应用，并配置该应用的属性。
4. 在 [odin](https://github.com/offer365/odin) 中生产序列号，在 edda 中解析，并生成license。
5. 激活 [odin](https://github.com/offer365/odin)。

## License
[Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0.html).