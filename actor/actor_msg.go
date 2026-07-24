package actor

import "time"

// ActorMsg 表示 Actor 之间传递并交给 Worker 处理的消息信封。
//
// ActorMsg 保存消息 ID、业务载荷、来源 Actor、目标 Actor、链路 ID、同步调用结果，
// 以及 Actor.Call 类接口等待完成时使用的通知通道。具体实现由 yxcore 在内部创建。
type ActorMsg interface {
	// GetMsgId 返回业务或框架内部消息 ID。
	GetMsgId() uint32
	// GetData 返回消息携带的业务载荷。
	GetData() interface{}
	// GetBytesData 返回 Bytes 投递接口携带的原始 RPC 载荷。
	GetBytesData() []byte
	// GetFromTyp 返回来源 Actor 类型。
	GetFromTyp() uint32
	// GetFromId 返回来源 Actor ID。
	GetFromId() uint64
	// SetToTyp 设置目标 Actor 类型，主要用于 RPC 路由。
	SetToTyp(toTyp uint32)
	// GetToTyp 返回目标 Actor 类型。
	GetToTyp() uint32
	// SetToId 设置目标 Actor ID，主要用于 RPC 路由。
	SetToId(toId uint64)
	// GetToId 返回目标 Actor ID。
	GetToId() uint64
	// SetTraceId 设置随 Actor 调用链透传的链路 ID。
	SetTraceId(traceId uint32)
	// GetTraceId 返回随 Actor 调用链透传的链路 ID。
	GetTraceId() uint32
	// GetResultCode 返回 Done 写入的处理结果码。
	GetResultCode() int32
	// GetResultData 返回 Done 写入的处理结果数据。
	GetResultData() interface{}
	// SetAsyncMode 标记该消息是否来自 Go 类异步投递接口。
	SetAsyncMode(bAsyncMode bool)
	// IsAsyncMode 返回该消息是否为不需要唤醒等待方的异步消息。
	IsAsyncMode() bool
	// WaitTimeout 等待 Done 通知完成，超时返回 ResCodeMsgTimeout。
	WaitTimeout(timeout time.Duration) int32
	// GetDoneChan 返回同步消息完成时由 Done 通知的通道。
	GetDoneChan() chan bool
	// Done 写入处理结果，并唤醒正在等待该同步消息的调用方。
	Done(code int32, data interface{})
}

// MsgWaiter 表示同步等待 ActorMsg 完成的等待器。
//
// Actor 实现 MsgWaiter 后，可以在同步调用其他 Actor 时继续处理自己的只读队列，
// 避免只读消息被长时间阻塞。非 Actor 调用方通常传 nil，由 yxcore 直接等待
// ActorMsg.GetDoneChan。
type MsgWaiter interface {
	// WaitMsgDone 等待 doneChan 收到完成通知或等待超时。
	WaitMsgDone(doneChan chan bool, timeout time.Duration) int32
}
