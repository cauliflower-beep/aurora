## 关于aurora

zinx框架的注释版。

上次更新到：

https://github.com/aceld/zinx/commit/46a2e95a3f3c5998e930792f81207a37ef43cec4

## 目录说明

### aiface

主要提供aurora全部抽象层接口定义。包括：

`IServer`服务module接口；

`IRrouter`路由module接口；

`IConnection`连接module层接口；

`IMessage`消息module层接口；

`IDataPack`消息拆解接口；

`IMsgHandler`消息处理及协程池接口；

`IClient`客户端接口。