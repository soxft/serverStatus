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
    -host string
          server ip:port, example:127.0.0.1:8282
    -tag string
          server tag
    -token string
          token of server
```
- 其中host用于指明服务器地址, tag用于表注服务器名称, token为服务端验证token
- 命令执行格式: `./serverStatus -host=127.0.0.1:8282 -tag=server1 -token=myToken`

## copyright

- 服务端基于 Gateway-workerman
- 网页使用 React 框架构建
  - 前端样式采用 Ant Design
- 客户端采用 Golang 开发