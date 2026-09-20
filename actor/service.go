package actor

// Service 是一组同 typ Actor 的容器，自身也是一个 Actor（拥有独立的消息循环）。
// Service 负责按需创建/销毁子 Actor、转发 Tick、回收闲置 Actor，
// 并通过 OnInit/OnStart/OnStop/OnCreateActor/OnRemoveActor 等钩子让业务接入生命周期。
type Service interface {
	// Actor 让 Service 复用 Actor 的消息循环和生命周期能力。
	Actor

	// SetChecked 设为 true 时，GetActor 会先调用 Server.IsLocalActor 校验归属，
	// 不属于本机的 Actor 直接返回 false（业务侧需走 RPC）。
	//
	// false 表示该 Service 不参与集群分配，所有 Actor 都在本机创建（如纯本地服务）。
	SetChecked(bCheck bool)
	// IsChecked 返回该 Service 是否启用 Actor 归属校验。
	IsChecked() bool
	// SetTickType 设置该 Service 下所有 Actor 接收的 Tick 类型，取 TickTyp* 位掩码组合。
	SetTickType(tickTyp int)
	// GetTickType 返回该 Service 当前启用的 Tick 类型位掩码。
	GetTickType() int
	// SetHotUpdateEnabled 设置该 Service 是否参与框架热更；默认 true。
	// SetHotUpdateEnabled(enabled bool)
	// IsHotUpdateEnabled 返回该 Service 是否参与框架热更。
	// IsHotUpdateEnabled() bool

	// GetActor 按 ID 获取或按需创建 Actor，是可能阻塞的同步调用。
	//
	// 若调用方本身是 Actor，waiter 应传调用方 Actor，使其等待期间仍可处理只读消息；
	// 非 Actor 调用方传 nil。
	GetActor(id uint64, waiter MsgWaiter) (Actor, bool)
	// RemoveActor 异步停止并移除指定 Actor；若其正在执行消息，会等待当前消息结束。
	RemoveActor(id uint64)
	// IsActorExist 返回指定 Actor 当前是否处于正常运行状态；正在停止或已标记移除时返回 false。
	IsActorExist(id uint64) bool
	// GetActorIds 返回当前所有 Actor ID，包含正在停止中的 Actor 以避免并发创建冲突。
	GetActorIds() []uint64
	// GetActorMsgCnt 返回每个 Actor 的待处理消息数，用于状态上报。
	GetActorMsgCnt() map[uint64]uint32
	// ClearStopActors 扫描并回收闲置 Actor 和已经完全退出的 Actor。
	ClearStopActors()

	// OnInit 在 Service 初始化时调用。
	OnInit() error
	// OnStart 在 Service 启动时调用。
	OnStart()
	// OnStop 在 Service 停止时调用。
	OnStop()
	// OnTick 在 Server 派发 Tick 时调用，tickTyp 取 TickTyp* 常量。
	OnTick(tickTyp int, dt int64)
	// OnCreateActor 在 Actor 创建前调用，业务可用于鉴权或配额检查；返回非 ResCodeSucc 表示拒绝创建。
	OnCreateActor(actorId uint64) int32
	// OnRemoveActor 在 Actor 从 Service 移除时调用。
	OnRemoveActor(actorId uint64)
}
