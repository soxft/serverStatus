# serverStatus
a simple serverStatus monitor

> stil under developing

## PATH

```
|_ web     网页 React 源码
|_ server  服务端
|_ client  服务器客户端
```



## client 端使用方法
```shell
  Usage of serverStatus:
      -duration int
            Data send Duration, Unit:ms (default 5000)
      -host string
            Server address (including protocol,ip and port), ex: ws://127.0.0.1:8282
      -tag string
            Server tag (default "xcsoftsMBP")
      -token string
            Server token
```
- 其中host用于指明服务器地址, tag用于表注服务器名称, token为服务端验证token
- 命令执行格式: `./serverStatus -host=ws://127.0.0.1:8282 -tag=server1 -token=myToken`

## copyright

- 服务端基于 Gateway-workerman
- 网页使用 React 框架构建
  - 前端样式采用 Ant Design 与 Ant Design Charts
- 客户端采用 Golang 开发