package network

// Conn 抽象一条双向连接（TCP、WebSocket、或其他实现）。
// 框架保证：
//   - 写入通过内部 writePackChan 异步发送，调用 Write 不阻塞底层 IO，队列满返回 ErrConnWriteFull
//   - Write 默认复制传入数据；buffRef 非空时会在写完后尝试归还到缓冲池
//   - 主动 Close 后会尽量将队列剩余数据写完再断开
//
// 业务层只需要面向 Conn 编程，无需关心底层是 net.Conn 还是 websocket.Conn。
type Conn interface {
	// GetConnId 返回框架内自增生成的连接 ID，用于日志和统计。
	GetConnId() uint32
	// SetAgentId 设置业务层绑定的代理 ID，例如玩家 ID。
	SetAgentId(id uint64)
	// GetAgentId 返回业务层绑定的代理 ID。
	GetAgentId() uint64

	// GetRemoteAddr 返回连接对端地址。
	GetRemoteAddr() string
	// Read 阻塞读取，反复调用底层读取直到填满 b 或出错。
	Read(b []byte) error
	// Write 将数据投递到发送队列，队列满时返回 ErrConnWriteFull。
	//
	// copyBuffer 为 true 时会先复制一份再入队，调用方可立刻复用 data；
	// 为 false 时直接使用传入 buffer。
	// buffRef 非空时，连接写完后会尝试把对应 buffer 归还到池中。
	Write(data []byte, buffRef *[]byte, copyBuffer bool) error
	// Close 主动关闭连接，并尽量写完发送队列中的剩余数据。
	Close()

	// SetManager 设置管理当前连接的 ConnMgr。
	SetManager(mgr ConnMgr)
	// Start 在独立 goroutine 启动读写循环，reader 决定如何拆包。
	Start(reader PackReader)
	// Stop 停止连接读写循环。
	Stop()
}

// PackReader 业务自定义的拆包器。
// 框架的 Conn.readLoop 调用顺序为 ReadHead -> ReadData -> GetPack -> OnConnRecv -> Reset，循环往复。
// 这样既可以处理「定长包头 + 变长包体」也可以处理「分隔符」等格式，由实现决定。
type PackReader interface {
	// ReadHead 读取包头，并解析包体长度等元信息。
	ReadHead(conn Conn) error
	// ReadData 按包头信息读取剩余包体数据。
	ReadData(conn Conn) error
	// GetPack 返回当前完整包的字节数据，通常不包含包头。
	GetPack() []byte
	// Reset 清空内部状态，准备读取下一个包。
	Reset()
}
