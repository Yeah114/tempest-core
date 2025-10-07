# tempest-core

`tempest-core` 提供一个基于 [Fatalder](https://github.com/Yeah114/Fatalder) 控制器的 gRPC 接入点，实现与 FateArk 相同的 API 定义，便于在网易租赁服环境下复用既有自动化工具。

## 功能特性

- ✅ gRPC 协议完全兼容 FateArk：沿用相同的 proto、消息结构与服务名称。
- ✅ 使用 Fatalder `Control` 进行机器人登录与生命周期管理。
- ✅ 打通指令、监听、玩家管理、工具等核心能力。
- ✅ 支持多路监听（数据包、聊天、命令方块等）并通过本地消息总线分发。
- ✅ 提供统一的 `tempestd` 守护进程，默认监听 `0.0.0.0:20919`。

## 目录结构

```
cmd/tempestd          # 服务器入口
internal/app          # 共享状态与消息总线
internal/server       # 各 gRPC Service 实现
network_api           # protoc 生成的 Go 代码
proto                 # gRPC proto 定义
modules/Fatalder      # Fatalder 内联依赖
```

## 构建与运行

环境需求：

- Go 1.20+（项目 go.mod 指定为 1.25.1）

编译：

```bash
go build ./cmd/tempestd
```

运行：

```bash
./tempestd -a 0.0.0.0 -p 20919
```

启动后即可通过 FateArk 兼容的 gRPC 客户端调用。连接租赁服需在 `FateReversalerService.NewFateReversaler` 请求中提供认证信息；连接断开时可通过 `WaitDead` 订阅退出原因。

## 注意事项

- 项目内联 Fatalder 及其依赖（FunShuttler、WaterStructure、blocks 等），无需额外拉取仓库。
- 当前 `InterceptPlayerJustNextInput` 为占位实现，只返回成功状态。
- 监听类接口内部采用非阻塞队列，若消费速度不足可能丢弃事件，请按需在客户端侧处理。

欢迎在现有基础上继续扩展命令封装、事件缓存或认证机制。
