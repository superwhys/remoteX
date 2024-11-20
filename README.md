# RemoteX

## TODO
- [x] Basic node connection by TCP
- [x] Get OS info
- [x] List remote node dir 
- [x] Adapt to UDP, Quic
- [x] Add file sync(upload, download)
- [x] Tunnel Forward
- [x] Tunnel Reverse
- [ ] Screenshot
- [ ] Run Command

## Message Body
```
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|         Message Length        |
|           (32 bits)           |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
/                               /
\             Message           \
/                               /
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
```

message read-write utils: `pkg/protoutils/reader.go`, `pkg/protoutils/writer.go`

## Process
1. Listen local port
2. Dial remote node or Accept remote connection
3. Exchange node info
4. Start heartbeat send
5. Wait for command

## Usage


## Develop
该项目许多结构都使用的protobuf，你可以自定义修改它们(`./internal/proto`)

修改了之后需要运行下面的指令生成对应的文件:

`go generate internal/proto/generate.go`

项目使用了 `DDD` 架构进行开发, 分为了四个领域: `Auth`, `Command`, `Connection`, `Node`
具体各目录如下：
```bash
.
├── api                   // 接口层
├── config                // 配置文件
├── domain                // 领域层
│   ├── auth
│   ├── command
│   ├── connection
│   └── node
├── internal
│   ├── connection
│   ├── filesync
│   ├── filesystem
│   └── proto
├── main.go
├── pkg                   // 基础设施层
├── server                // 应用层
│   ├── auth
│   ├── command
│   ├── command.go
│   ├── connection
│   ├── connection.go
│   ├── handshake.go
│   ├── heartbeat.go
│   ├── node
│   ├── option.go
│   └── server.go
└── tools.go 
```
