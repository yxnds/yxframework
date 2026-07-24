package actor

import "time"

// EventHandleFn 事件处理回调。
// eventId 通过 Worker.SetEventHandler 注册时指定；eventData 由 SendEvent 透传。
type EventHandleFn func(eventId uint32, eventData interface{})

// Actor 是框架的核心并发单元，以 (typ, id) 唯一标识。
// 每个 Actor 独占一个 goroutine 循环消费两条消息队列（读写队列 + 只读队列），
// 通过 Worker 承载业务逻辑。
//
// 发消息接口分三大类：
//   - Go*       异步投递，调用方不等待结果
//   - Call*     同步等待结果，可指定超时
//   - SendEvent 走 MsgIdEvent 通道的特殊广播消息
//
// Call* 的 timeout <= 0 时使用框架默认等待超时。
//
// 后缀含义：
//   - 1/2     1 不带 from 来源信息；2 携带 fromTyp/fromId（用于路由回包、链路追踪）
//   - Data    携带业务数据载荷
//   - Result  返回值带数据（否则只返回错误码）
type Actor interface {
	// MsgWaiter 让 Actor 在同步等待期间继续处理自身的只读消息。
	MsgWaiter

	// SetSelf 设置 Actor 的自身引用，用于基础类型回调到真实子类型实例。
	SetSelf(self Actor)
	// GetSelf 返回 Actor 的自身引用。
	GetSelf() Actor

	// GetService 返回所属 Service，同 typ 的 Actor 由该 Service 管理。
	GetService() Service
	// RemoveFromService 通知所属 Service 将当前 Actor 移除并停止。
	RemoveFromService()

	// SetWorker 绑定承载业务逻辑的 Worker。
	SetWorker(worker Worker)
	// GetWorker 返回当前绑定的 Worker。
	GetWorker() Worker
	// HotUpdateWorker 投递新 Worker 并触发热更流程。
	//
	// Actor 消息循环收到热更请求后调用 newWorker.OnHotUpdate(oldWorker)，
	// 并使用其返回值接管后续消息处理。
	HotUpdateWorker(newWorker Worker)

	// GetId 返回 Actor ID。
	GetId() uint64
	// GetTyp 返回 Actor 类型。
	GetTyp() uint32
	// SetTag 设置 Actor 的业务标签，常用于日志和调试定位。
	SetTag(tag string)
	// GetTag 返回 Actor 的业务标签。
	GetTag() string
	// SetFreeTime 设置闲置回收时间（秒）。
	//
	// 大于 0 时，超过该时长无消息处理的 Actor 会被 Service.ClearStopActors 自动清理。
	SetFreeTime(freeTime int64)
	// GetFreeTime 返回闲置回收时间（秒）。
	GetFreeTime() int64

	// SendEvent 通过 MsgIdEvent 投递事件。需配合 Worker.SetEventHandler 注册的 EventHandleFn。
	SendEvent(eventId uint32, eventData interface{})

	// Go1 异步投递无载荷消息，不携带来源 Actor 信息。
	Go1(msgId uint32, traceId uint32)
	// Go2 异步投递无载荷消息，并携带来源 Actor 类型和 ID。
	Go2(msgId uint32, fromTyp uint32, fromId uint64, traceId uint32)
	// GoData1 异步投递带业务载荷的消息，不携带来源 Actor 信息。
	GoData1(msgId uint32, data interface{}, traceId uint32)
	// GoData2 异步投递带业务载荷的消息，并携带来源 Actor 类型和 ID。
	GoData2(msgId uint32, data interface{}, fromTyp uint32, fromId uint64, traceId uint32)
	// GoBytes1 异步投递已序列化的 RPC 字节，不携带来源 Actor 信息。
	GoBytes1(msgId uint32, data []byte, traceId uint32)
	// GoBytes2 异步投递已序列化的 RPC 字节，并携带来源 Actor 类型和 ID。
	GoBytes2(msgId uint32, data []byte, fromTyp uint32, fromId uint64, traceId uint32)

	// Call1 同步投递无载荷消息，不携带来源 Actor 信息，只返回结果码。
	//
	// 若调用方本身是 Actor，waiter 应传调用方 Actor，使其等待期间仍可处理只读消息；
	// 非 Actor 调用方传 nil。
	Call1(waiter MsgWaiter, msgId uint32, traceId uint32, timeout time.Duration) int32
	// Call2 同步投递无载荷消息，并携带来源 Actor 类型和 ID，只返回结果码。
	Call2(waiter MsgWaiter, msgId uint32, fromTyp uint32, fromId uint64, traceId uint32, timeout time.Duration) int32
	// CallResult1 同步投递无载荷消息，不携带来源 Actor 信息，返回处理结果数据和结果码。
	CallResult1(waiter MsgWaiter, msgId uint32, traceId uint32, timeout time.Duration) (interface{}, int32)
	// CallResult2 同步投递无载荷消息，并携带来源 Actor 类型和 ID，返回处理结果数据和结果码。
	CallResult2(waiter MsgWaiter, msgId uint32, fromTyp uint32, fromId uint64, traceId uint32, timeout time.Duration) (interface{}, int32)

	// CallData1 同步投递带业务载荷的消息，不携带来源 Actor 信息，只返回结果码。
	CallData1(waiter MsgWaiter, msgId uint32, data interface{}, traceId uint32, timeout time.Duration) int32
	// CallData2 同步投递带业务载荷的消息，并携带来源 Actor 类型和 ID，只返回结果码。
	CallData2(waiter MsgWaiter, msgId uint32, data interface{}, fromTyp uint32, fromId uint64, traceId uint32, timeout time.Duration) int32
	// CallDataResult1 同步投递带业务载荷的消息，不携带来源 Actor 信息，返回处理结果数据和结果码。
	CallDataResult1(waiter MsgWaiter, msgId uint32, data interface{}, traceId uint32, timeout time.Duration) (interface{}, int32)
	// CallDataResult2 同步投递带业务载荷的消息，并携带来源 Actor 类型和 ID，返回处理结果数据和结果码。
	CallDataResult2(waiter MsgWaiter, msgId uint32, data interface{}, fromTyp uint32, fromId uint64, traceId uint32, timeout time.Duration) (interface{}, int32)
	// CallBytes1 同步投递已序列化的 RPC 字节，不携带来源 Actor 信息，只返回结果码。
	CallBytes1(waiter MsgWaiter, msgId uint32, data []byte, traceId uint32, timeout time.Duration) int32
	// CallBytes2 同步投递已序列化的 RPC 字节，并携带来源 Actor 类型和 ID，只返回结果码。
	CallBytes2(waiter MsgWaiter, msgId uint32, data []byte, fromTyp uint32, fromId uint64, traceId uint32, timeout time.Duration) int32
	// CallBytesResult1 同步投递已序列化的 RPC 字节，不携带来源 Actor 信息，返回处理结果数据和结果码。
	CallBytesResult1(waiter MsgWaiter, msgId uint32, data []byte, traceId uint32, timeout time.Duration) (interface{}, int32)
	// CallBytesResult2 同步投递已序列化的 RPC 字节，并携带来源 Actor 类型和 ID，返回处理结果数据和结果码。
	CallBytesResult2(waiter MsgWaiter, msgId uint32, data []byte, fromTyp uint32, fromId uint64, traceId uint32, timeout time.Duration) (interface{}, int32)

	// SetService 设置所属 Service，主要供框架内部创建 Actor 时调用。
	SetService(service Service)
	// SetLastActiveTime 设置最后活跃时间戳，主要供框架内部维护闲置回收状态。
	SetLastActiveTime(t int64)
	// GetLastActiveTime 返回最后活跃时间戳。
	GetLastActiveTime() int64
	// GetMsgCnt 返回读写队列待处理消息总数，用于状态上报与拥塞判定。
	GetMsgCnt() uint32
	// Init 初始化 Actor 及其 Worker。
	Init() error
	// Start 启动 Actor 消息循环。
	Start()
	// Stop 请求停止 Actor 消息循环。
	Stop()
	// IsExit 返回 Actor 是否已经退出。
	IsExit() bool
	// Tick 投递定时消息（MsgIdTick），由 Server.OnUpdate 根据 Service.GetTickType 触发。
	//
	// tickTyp 取 TickTyp* 系列常量；dt 为距上一次 Update 的毫秒数（仅 TickTypUpdate 有意义）。
	Tick(tickTyp int, dt int64)
}

// RpcActor 用于跨 Server 的 Actor 调用。
// 当目标 Actor 不在本地（IsLocalActor 返回 false）时，框架通过 RpcActor 走 TCP 网络转发消息。
// 业务侧通常通过 framework_api.Go*/Call* 间接使用，不直接持有 RpcActor。
type RpcActor interface {
	// Actor 让 RpcActor 复用普通 Actor 的消息循环和生命周期能力。
	Actor

	// Connect 建立到目标 server 的 TCP 长连接；连接断开时框架会自动重连，超过 RpcRetryCnt 后放弃。
	Connect() error
	// RpcGo 跨服异步投递；data 会经 ServerConf.RpcDataToBytes 序列化后发送。
	RpcGo(msgId uint32, data interface{}, fromTyp uint32, fromId uint64, toTyp uint32, toId uint64, traceId uint32)
	// RpcGoBytes 跨服异步投递已序列化的 RPC 字节。
	RpcGoBytes(msgId uint32, data []byte, fromTyp uint32, fromId uint64, toTyp uint32, toId uint64, traceId uint32)
	// RpcCall 跨服同步等待返回；waiter 语义与 Actor.Call* 一致。
	RpcCall(waiter MsgWaiter, msgId uint32, data interface{}, fromTyp uint32, fromId uint64, toTyp uint32, toId uint64, traceId uint32, timeout time.Duration) (interface{}, int32)
	// RpcCallBytes 跨服同步投递已序列化的 RPC 字节，并等待返回。
	RpcCallBytes(waiter MsgWaiter, msgId uint32, data []byte, fromTyp uint32, fromId uint64, toTyp uint32, toId uint64, traceId uint32, timeout time.Duration) (interface{}, int32)

	// OnStop 在 RpcActor 停止时释放网络连接和等待中的调用。
	OnStop()
}
