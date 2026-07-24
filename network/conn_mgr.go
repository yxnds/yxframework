package network

import (
	"net"
	"time"
)

// ConnHandler 业务侧的连接生命周期回调。
// 注意：OnConnRecv 在连接的读 goroutine 中被调用，业务回调内部不应做重 IO 阻塞操作，
// 通常的做法是把 pack 投递到 Actor 消息队列后立即返回。
type ConnHandler interface {
	// OnCreatePackReader 在每条新连接建立时调用，返回该连接专用的拆包器实例。
	OnCreatePackReader(mgr ConnMgr) PackReader

	// OnConnOpen 在连接建立并加入管理器后调用。
	OnConnOpen(c Conn) error
	// OnConnClose 在连接关闭并从管理器移除后调用。
	OnConnClose(c Conn)
	// OnConnError 在连接读写发生底层 IO 错误时调用。
	OnConnError(c Conn)
	// OnConnRecv 在连接收到完整网络包后调用。
	OnConnRecv(c Conn, pack []byte)
}

// ConnMgr 长连接管理器（TCP / WebSocket 等）。
// 一个 ConnMgr 可同时承担 Listen（被动接入）和 Dial（主动外联）两种角色，
// 内部维护连接列表并统一关闭。
type ConnMgr interface {
	// SetHandler 设置连接生命周期回调处理器。
	SetHandler(handler ConnHandler)
	// GetHandler 返回连接生命周期回调处理器。
	GetHandler() ConnHandler
	// SetGate 设置关联的 ConnGate。
	SetGate(gate ConnGate)
	// GetGate 返回当前关联的 ConnGate。
	GetGate() ConnGate
	// GetWriteQueueSize 返回每条连接的写队列容量。
	GetWriteQueueSize() int32
	// GetWriteTimeoutSec 返回每条连接的写超时时间，单位秒。
	GetWriteTimeoutSec() int64

	// Listen 启动长连接监听，protoTye 通常为 tcp 或 ws。
	Listen(protoTye string, address string) error
	// CloseListen 关闭监听端口，但不直接关闭已经建立的连接。
	CloseListen() error

	// DialTimeout 主动拨号建立外联连接；timeout 为连接超时。
	DialTimeout(protoTye string, address string, timeout time.Duration) (Conn, error)

	// CreateConn 将底层 net.Conn 包装为框架 Conn，一般由 Listen 内部使用。
	CreateConn(conn net.Conn) (Conn, error)
	// CloseConn 关闭并移除指定连接。
	CloseConn(c Conn)
	// CloseAllConn 关闭并移除当前管理的全部连接。
	CloseAllConn()
}
