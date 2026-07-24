package actor

// Worker 承载 Actor 的业务逻辑：注册消息处理器、维护业务状态、响应生命周期与时间回调。
// 一个 Actor 持有一个 Worker；通过 HotUpdateWorker 可以在不停 Actor 的情况下替换 Worker。
//
// 消息处理分两类：
//   - 写消息（RegisterMsgHandler）：可改变 Actor 状态，串行执行
//   - 只读消息（RegisterReadMsgHandler）：不改状态，写消息处理完后批量消费，
//     且在 Call 等待期间允许执行——这是 waiter==self 不死锁的关键
//
// handler 支持以下函数签名（按是否带 data / 是否返回 (data, code) / 仅 code 分组）：
//
//	func(msgId, fromTyp, fromId, traceId)
//	func(msgId, data, fromTyp, fromId, traceId)
//	func(msgId, fromTyp, fromId, traceId) int32
//	func(msgId, fromTyp, fromId, traceId) (interface{}, int32)
//	func(msgId, data, fromTyp, fromId, traceId) int32
//	func(msgId, data, fromTyp, fromId, traceId) (interface{}, int32)
//	func(msg)
//	func(msg) int32
//	func(msg) (interface{}, int32)
type Worker interface {
	// SetSelf 设置 Worker 的真实实例引用，用于基础类型回调到子类型重写的钩子。
	SetSelf(self Worker)
	// GetSelf 返回 Worker 的真实实例引用。
	GetSelf() Worker

	// Wrap 设置热更 Worker 内嵌的静态基类对象。
	//
	// 创建热更 Worker 时，框架会把原静态 Worker 传入，让热更对象可以内嵌并复用
	// 静态对象的状态与能力；Actor.HotUpdateWorker 不会再次调用 Wrap。
	Wrap(base Worker)
	// Unwrap 返回热更 Worker 当前内嵌的静态基类对象。
	Unwrap() (Worker, bool)

	// SetVer 设置 Worker 热更版本号。
	SetVer(ver uint32)
	// GetVer 返回 Worker 热更版本号。
	GetVer() uint32

	// SetActor 绑定当前 Worker 所属的 Actor。
	SetActor(actor Actor)
	// GetActor 返回当前 Worker 所属的 Actor。
	GetActor() Actor

	// RegisterMsgHandler 注册写消息处理器。msgId 为 0 或重复注册都会被框架忽略并打 Warn。
	RegisterMsgHandler(msgId uint32, handler interface{})
	// UnregisterMsgHandler 取消注册指定写消息处理器。
	UnregisterMsgHandler(msgId uint32)
	// GetMsgHandler 返回指定写消息处理器。
	GetMsgHandler(msgId uint32) (interface{}, bool)
	// SetDefaultMsgHandler 设置默认写消息处理器，未注册具体 msgId 时使用；传 nil 清除默认处理器。
	SetDefaultMsgHandler(handler interface{})
	// RegisterReadMsgHandler 注册只读消息处理器；处理器不应修改 Actor 状态。
	RegisterReadMsgHandler(msgId uint32, handler interface{})
	// UnregisterReadMsgHandler 取消注册指定只读消息处理器。
	UnregisterReadMsgHandler(msgId uint32)
	// GetReadOnlyMsgHandler 返回指定只读消息处理器。
	GetReadOnlyMsgHandler(msgId uint32) (interface{}, bool)
	// SetEventHandler 注册事件处理函数，用于响应 SendEvent 投递的 MsgIdEvent。
	SetEventHandler(eventId uint32, handler EventHandleFn)
	// RemoveEventHandler 移除指定事件处理函数。
	RemoveEventHandler(eventId uint32)
	// HandleMsg 由 Actor 消息循环调用，按 msgId 路由到处理器并在结束后调用 msg.Done。
	HandleMsg(msg ActorMsg)

	// OnInit 在 Worker 初始化时调用。
	OnInit() error
	// OnStart 在 Worker 启动时调用。
	OnStart()
	// OnStop 在 Actor 请求停止时调用。
	OnStop()
	// OnDestroy 在 Worker 销毁前调用。
	OnDestroy()
	// OnHandleMsgBegin 在每条消息处理前调用，返回 false 表示拦截并跳过后续 handler。
	OnHandleMsgBegin(msg ActorMsg) bool
	// OnRpcMsgKeepBytes 在 Bytes 投递消息进入 HandleMsg 前调用。
	//
	// 返回 true 表示当前消息跳过 RpcBytesToData，业务通过 ActorMsg.GetBytesData 读取原始载荷；
	// 返回 false 表示框架先调用 RpcBytesToData，并通过 ActorMsg.GetData 传给 handler。
	OnRpcMsgKeepBytes(msg ActorMsg) bool
	// OnHandleMsgEnd 在每条消息处理后调用。
	OnHandleMsgEnd(msg ActorMsg)
	// OnHotUpdate 在热更新时调用。
	//
	// old 是当前 Actor 正在使用的旧 Worker；返回值会成为新的活动 Worker。
	// 业务需要在这里完成状态迁移和必要的包装逻辑。
	OnHotUpdate(old Worker) Worker

	// DoTick 处理 MsgIdTick，并按 tickTyp 分派到对应时间变化钩子。
	DoTick(tickTyp int, dt int64)
	// OnUpdate 在 TickTypUpdate 启用时按定时器周期调用，dt 为距离上次调用的毫秒数。
	OnUpdate(dt int64)
	// OnSecondChanged 在 TickTypeSecond 启用且自然秒变化时调用。
	OnSecondChanged()
	// OnMinuteChanged 在 TickTypMinute 启用且自然分钟变化时调用。
	OnMinuteChanged()
	// OnHourChanged 在 TickTypHour 启用且自然小时变化时调用。
	OnHourChanged()
	// OnDayChanged 在 TickTypDay 启用且自然日期变化时调用。
	OnDayChanged()
	// OnWeekChanged 在 TickTypWeek 启用且自然周变化时调用。
	OnWeekChanged()
	// OnMonthChanged 在 TickTypMonth 启用且自然月份变化时调用。
	OnMonthChanged()

	// OnInnerInit 框架内部初始化点，WorkerBase 在此注册 MsgIdTick / MsgIdEvent 处理器；
	//
	// 子类如需扩展应调用父类实现。
	OnInnerInit() error
	// OnInnerDestroy 框架内部销毁点，WorkerBase 在此释放内部处理器。
	OnInnerDestroy()
	// OnInnerMsgBegin 框架内部消息前置钩子，在 OnHandleMsgBegin 之前执行，返回 false 表示拦截并跳过后续 handler。
	OnInnerMsgBegin(msg ActorMsg) bool
	// OnInnerMsgEnd 框架内部消息后置钩子，在 OnHandleMsgEnd 之后执行。
	OnInnerMsgEnd(msg ActorMsg)

	// 热更状态迁移由 OnHotUpdate 负责；框架不会自动调用序列化/反序列化钩子。
}
