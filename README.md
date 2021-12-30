# serverStatus

> 一个 简易的服务器探针

## PATH

```
|_ web     网页 React 源码
|_ server  服务端
|_ client  服务器客户端
```

## 使用方法

> 需要部署三样, 分别为Client端, Web端, server端. 分别对应了根目录下的三个文件夹

- Client端为客户端, 用于推送服务器当前状态到server端. Client端不需要公网环境. 仅支持Linux服务器
- Server端分别与Client端和Web端构建Websocket通信, 用于转发信息. Server端需要处于公网环境
- Web端为状态监控面板

### Server端

> Server端使用 Workerman构建

将本项目下载到本地, 修改config.php中的token等信息, Server端提供了一个基于Workerman的简易Web服务器, 您可以将Web端编译后, 拷贝至/sever/Web/Applications/Web/ 目录下.

#### 启动Server端

在/Server目录下, 使用`composer install`补全环境依赖

通过`php start.php start`启动服务端. Websocket默认端口为8282, 内建Web端默认端口为8283.

#### 其他Server端常见命令
```shell
      php start.php start -d   # 以守护模式运行
      php start.php stop       # 停止服务端
      php start.php status -d  # 查看服务器状态
```

### Web端

> Web端采用React构建

在/web目录下使用`npm install`补全依赖, 修改/web/src/config.js中的服务器信息, 使用`npm run build`编译打包.

默认会打包至/web/build文件夹内. 您可以选择使用nginx等环境,或直接拷贝至 /server/Applications/Web/ 目录下 使用内建web服务器

### client 端

> 您可以clone本项目后, 自行前往client文件夹编译运行, 也可以直接在Release中下载已经编译过的二进制文件.

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
- 其中host用于指明服务器地址, tag用于表注服务器名称, token为服务端验证token, dutation为服务器信息获取间隔(单位毫秒)
- 命令执行格式: `./serverStatus -host=ws://127.0.0.1:8282 -tag=server1 -token=myToken`

## copyright

- 服务端基于 Gateway-workerman
- 网页使用 React 框架构建
  - 前端样式采用 Ant Design 与 Ant Design Charts
- 客户端采用 Golang 开发