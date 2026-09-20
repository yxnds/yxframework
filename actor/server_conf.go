package actor

// ServerConf Server 启动配置。由业务侧填好后通过 Server.SetConfig 注入。
type ServerConf struct {
	// ServerId 表示本机在集群中的唯一 ID。
	ServerId uint32
	// SvrClusterId 表示本机所属集群 ID，一个集群对应一个 Center。
	SvrClusterId uint32
	// SvrGroupId 表示本机所属分组 ID。
	SvrGroupId uint32
	// Group 表示本机所属分组路由名；为空表示未配置，不会自动映射到默认分组。
	Group string
	// ServerHost 表示上报给 Center 的对外地址；填 UnknownServerAddr 时由 Center 使用 RemoteAddr 替代。
	ServerHost string
	// RpcPort 表示跨服 RPC 监听端口。
	RpcPort uint32
	// StatPort 表示 /status HTTP 健康检查端口。
	StatPort uint32

	// TimerIntvMs 表示 GlobalTimer 心跳间隔，单位毫秒，决定 TickTypUpdate 频率。
	TimerIntvMs int64

	// MsgQueueSize 表示单个 Actor 读写队列容量，写满时返回 ResCodeMsgQueueFull。
	MsgQueueSize int32
	// MsgTimeoutSec 表示 Service 创建 Actor 和本地 Actor Call 默认等待超时时间，单位秒。
	MsgTimeoutSec int64
	// ShutdownTimeoutMSec 表示服务停服时等待 Actor 退出的最长时间，单位毫秒。
	ShutdownTimeoutMSec uint32
	// ReadMsgCntOnceHandle 表示每处理一条写消息后最多批量消费的只读消息数量。
	ReadMsgCntOnceHandle int32
	// CongestionTimeoutCnt 表示拥塞判定阈值，时间窗内累计超时达到该值即视为拥塞。
	CongestionTimeoutCnt uint32
	// CongestionMsgCostMSec 表示单条消息从入队到处理完成的耗时阈值，单位毫秒。
	CongestionMsgCostMSec int64
	// CongestionClearSec 表示拥塞计数衰减时间窗，单位秒，超过则计数清零。
	CongestionClearSec int64

	// NilServiceTyp 表示内置 NilService 的类型，用于无归属消息的兜底处理。
	NilServiceTyp uint32
	// CongestionServiceTyp 表示内置拥塞兜底 Service 的类型，用于目标 Actor 拥塞时的兜底处理。
	CongestionServiceTyp uint32

	// RpcServiceTyp 表示框架内置 RPC 路由 Service 的类型，仅用于内部注册和路由。
	RpcServiceTyp uint32
	// RpcActorFreeTime 表示 RpcActor 闲置回收时间，单位秒。
	RpcActorFreeTime int64
	// RpcWriteQueueSize 表示 RPC 连接写队列容量。
	RpcWriteQueueSize int32
	// RpcWriteTimeoutSec 表示 RPC 连接写超时时间，单位秒。
	RpcWriteTimeoutSec int64
	// RpcDialTimeoutSec 表示 RPC 主动拨号超时时间，单位秒。
	RpcDialTimeoutSec uint32
	// RpcRetryCnt 表示 RPC 连接失败重试次数，超过后放弃并移除 RpcActor。
	RpcRetryCnt int
	// RpcWaitTimeoutMSec 表示等待对端 RPC 响应的超时时间，单位毫秒。
	RpcWaitTimeoutMSec uint32
	// RpcInboundMaxConcurrent 表示本机同时处理的入站 RPC 请求上限。
	RpcInboundMaxConcurrent int

	// CenterAddr 表示 Center 的网络访问地址。
	CenterAddr string
	// CenterDialTimeoutSec 表示连接 Center 的拨号超时时间，单位秒。
	CenterDialTimeoutSec uint32
	// CenterWriteQueueSize 表示 Center 连接写队列容量。
	CenterWriteQueueSize int32
	// CenterWriteTimeoutSec 表示 Center 连接写超时时间，单位秒。
	CenterWriteTimeoutSec int64
	// CenterRetryCnt 表示连接 Center 失败后的最大重试次数。
	CenterRetryCnt int
	// CenterWaitTimeoutMSec 表示等待 Center 响应的兼容默认超时时间，单位毫秒。
	CenterWaitTimeoutMSec uint32
	// CenterRequestWaitTimeoutMSec 表示等待 Center 查询/分配响应的超时时间，单位毫秒。
	CenterRequestWaitTimeoutMSec uint32
	// CenterSyncWaitTimeoutMSec 表示等待 Center 心跳同步响应的超时时间，单位毫秒。
	CenterSyncWaitTimeoutMSec uint32
	// CenterRequestPoolSize 表示 Center 查询/分配请求连接池容量。
	CenterRequestPoolSize int
	// CenterFailFastMSec 表示 Center 查询/分配失败后的短熔断时间，单位毫秒。
	CenterFailFastMSec uint32
	// CenterSyncIntvMSec 表示向 Center 上报本机状态的周期，单位毫秒。
	CenterSyncIntvMSec uint32

	// GetRpcRspMsgId 根据 RPC 请求消息 ID 推导响应消息 ID。
	GetRpcRspMsgId func(msgId uint32) uint32
	// RpcBytesToData 将 RPC 网络包字节反序列化为业务数据。
	RpcBytesToData func(msgId uint32, data []byte) (interface{}, error)
	// RpcDataToBytes 将 RPC 业务数据序列化为网络包字节。
	RpcDataToBytes func(msgId uint32, data interface{}) ([]byte, error)

	// GetMainVer 返回主程序版本号，用于状态上报。
	GetMainVer func() string
	// GetHotVer 返回当前热更版本号，用于状态上报。
	GetHotVer func() string
	// RegisterServices 在 Server.Init 注册完内置 Nil/RPC 路由 Service 后调用，用于注册业务 Service。
	RegisterServices func(Server) error

	// Typ2WorkerCreateFn 保存普通 Actor Worker 工厂，键为 Service 类型。
	Typ2WorkerCreateFn map[uint32]WorkerCreateFn
	// Typ2WorkshopCreateFn 保存 Service 自身 Worker 工厂，键为 Service 类型。
	Typ2WorkshopCreateFn map[uint32]WorkerCreateFn

	Typ2Hotfixer map[uint32]ActorHotfixer
}
