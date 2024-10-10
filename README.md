# RemoteX

## TODO
- [x] Basic node connection by TCP
- [x] Get OS info
- [x] List remote node dir 
- [ ] Adapt to UDP, Quic
- [ ] Add file sync(upload, download)
- [ ] Reverse Shell
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
