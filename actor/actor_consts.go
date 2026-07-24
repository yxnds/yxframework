package actor

const (
	// ActorConnHeadLen 表示 Actor RPC 网络包长度头占用的字节数。
	ActorConnHeadLen uint32 = 4
	// ActorConnMaxPackSize 表示 Actor RPC 网络包允许的最大包体字节数。
	ActorConnMaxPackSize uint32 = 16 * 1024 * 1024
)

const (
	// CenterMsgIdSyncDataReq 表示 Server 向 Center 同步运行状态的请求消息 ID。
	CenterMsgIdSyncDataReq uint32 = 101
	// CenterMsgIdSyncDataRsp 表示 Center 响应 Server 状态同步的消息 ID。
	CenterMsgIdSyncDataRsp uint32 = 102

	// CenterMsgIdActorServerIdReq 表示查询 Actor 已分配 Server 的请求消息 ID。
	CenterMsgIdActorServerIdReq uint32 = 103
	// CenterMsgIdActorServerIdRsp 表示查询 Actor 已分配 Server 的响应消息 ID。
	CenterMsgIdActorServerIdRsp uint32 = 104

	// CenterMsgIdServerAssignReq 表示请求 Center 为 Actor 分配 Server 的消息 ID。
	CenterMsgIdServerAssignReq uint32 = 105
	// CenterMsgIdServerAssignRsp 表示 Center 返回 Actor 分配结果的消息 ID。
	CenterMsgIdServerAssignRsp uint32 = 106
)

const (
	// MsgIdCreateActor 表示 Service 创建子 Actor 的内部消息 ID。
	MsgIdCreateActor uint32 = 1
	// MsgIdRemoveActor 表示 Service 移除子 Actor 的内部消息 ID。
	MsgIdRemoveActor uint32 = 2
	// MsgIdTick 表示 Server.OnUpdate 派发的内部定时消息 ID。
	MsgIdTick uint32 = 3
	// MsgIdEvent 表示 Actor.SendEvent 投递事件时使用的内部消息 ID。
	MsgIdEvent uint32 = 4
	// MsgIdNoActorGo 表示异步调用缺失 Actor 时交给 NilService 兜底处理的消息 ID。
	MsgIdNoActorGo uint32 = 8
	// MsgIdNoActorCall 表示同步调用缺失 Actor 时交给 NilService 兜底处理的消息 ID。
	MsgIdNoActorCall uint32 = 9
)

const (
	// TickTypNone 表示 Service 不接收任何定时回调。
	TickTypNone int = 0
	// TickTypUpdate 表示按 ServerConf.TimerIntvMs 周期触发 Worker.OnUpdate。
	TickTypUpdate int = 1
	// TickTypeSecond 表示自然秒变化时触发 Worker.OnSecondChanged。
	TickTypeSecond int = 2
	// TickTypMinute 表示自然分钟变化时触发 Worker.OnMinuteChanged。
	TickTypMinute int = 4
	// TickTypHour 表示自然小时变化时触发 Worker.OnHourChanged。
	TickTypHour int = 8
	// TickTypDay 表示自然日期变化时触发 Worker.OnDayChanged。
	TickTypDay int = 16
	// TickTypWeek 表示自然周变化时触发 Worker.OnWeekChanged。
	TickTypWeek int = 32
	// TickTypMonth 表示自然月份变化时触发 Worker.OnMonthChanged。
	TickTypMonth int = 64
)

const (
	// ResCodeSucc 表示 Actor 消息处理成功。
	ResCodeSucc int32 = 0
	// ResCodeInvalidMsgId 表示消息 ID 无效、为零或没有注册处理器。
	ResCodeInvalidMsgId int32 = -1
	// ResCodeMsgQueueFull 表示目标 Actor 消息队列已达到 ServerConf.MsgQueueSize。
	ResCodeMsgQueueFull int32 = -2
	// ResCodeMsgTimeout 表示同步 Actor 调用等待超时。
	ResCodeMsgTimeout int32 = -3
	// ResCodeInvalidActorId 表示 Actor ID 为零或不满足当前操作要求。
	ResCodeInvalidActorId int32 = -4
	// ResCodeActorNotExist 表示目标 Actor 不存在或创建失败。
	ResCodeActorNotExist int32 = -5
	// ResCodeWorkerCreateFailed 表示对应 WorkerCreateFn 未配置或创建出的 Worker 为空。
	ResCodeWorkerCreateFailed int32 = -6
	// ResCodeActorStopping 表示目标 Actor 正在停止，已拒绝新消息。
	ResCodeActorStopping int32 = -7
	// ResCodeRpcWriteDataErr 表示 RPC 数据序列化或网络写出失败。
	ResCodeRpcWriteDataErr int32 = -8
	// ResCodeDataFormatErr 表示消息序列化、反序列化或处理器数据格式错误。
	ResCodeDataFormatErr int32 = -9
	// ResCodeSvrNotExist 表示 Center 没有找到可承载目标 Actor 的在线 Server。
	ResCodeSvrNotExist int32 = -10
	// ResCodeActorNotService 表示需要 Service 的场景传入了非 Service Actor。
	ResCodeActorNotService int32 = -11
	// ResCodeActorAssignFailed 表示 Center 为 Actor 分配 Server 失败。
	ResCodeActorAssignFailed int32 = -12
	// ResCodeCenterUnavailable 表示 Center 不可用、连接失败或响应异常。
	ResCodeCenterUnavailable int32 = -13
	// ResCodeInvalidGroup 表示路由 group 为空或非法。
	ResCodeInvalidGroup int32 = -14
	// ResCodeRpcInboundCongestion 表示目标 Server 入站 RPC 并发已满。
	ResCodeRpcInboundCongestion int32 = -15
	// ResCodeActorSelfCall 表示 Actor 同步调用自身写消息，框架已拒绝以避免自等待。
	ResCodeActorSelfCall int32 = -16
)

const (
	// UnknownServerAddr 表示 ServerHost 未知，中心服会使用连接的远端地址替换。
	UnknownServerAddr string = "0.0.0.0"
)
