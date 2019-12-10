# edda #

----

## what this? ##
- 这个项目仅仅是给 [odin](https://github.com/offer365/odin) 生成授权码。 

## 只有4g gRpc 接口
```
service Authorization {
    rpc Resolved (Cipher) returns (SerialNum); // 解析序列号
    rpc Authorized (AuthReq) returns (AuthResp); // 授权
    rpc Untied (UntiedReq) returns (Cipher); // 解绑
    rpc Cleared (Cipher) returns (Clear);  // 清除
}

```


## 安装运行 ##
```
go get github.com/offer365/edda
cd $GOPATH/src/github.com/offer365/edda
go build
./edda 
# 或 -l 指定监听端口
./edda -l 4567
```

> 访问 http://127.0.0.1:19527


## 使用介绍 ##
1. 先安装 edda 并运行。访问web端口，默认账号密码：admin:123 可在配置文件 edda.json 中修改。
2. 使用 gRpc 或 Restful api 接口生成 序列号。
3. 在 [edda](https://github.com/offer365/edda) 里根据约定新建应用，并配置该应用的属性。
5. 激活 [odin](https://github.com/offer365/odin)。

## License
[Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0.html).