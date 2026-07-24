package network

import "errors"

var (
	// ErrConnClose 表示连接关闭后仍调用 Write。
	ErrConnClose = errors.New("conn close")
	// ErrConnWriteFull 表示写队列已满，业务侧应限流或断开连接。
	ErrConnWriteFull = errors.New("write queue full")
	// ErrConnNotImplement 表示当前 ConnMgr 实现不支持该操作。
	ErrConnNotImplement = errors.New("not implement")
	// ErrConnUnderlyConnNil 表示 CreateConn 时传入的底层 net.Conn 为空。
	ErrConnUnderlyConnNil = errors.New("underly conn is nil")
	// ErrConnHandlerNil 表示未设置 ConnHandler 就接入新连接。
	ErrConnHandlerNil = errors.New("handler is nil")
	// ErrPackReaderNil 表示连接创建时未提供拆包器。
	ErrPackReaderNil = errors.New("pack reader is nil")
)
