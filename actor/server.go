package actor

import (
	"net/http"

	"github.com/wampee-sweet/yxframework/network"
)

// WorkerCreateFn Worker 工厂函数。
// 由 Server.AddService 注册到对应 typ；每次创建 Actor 时通过 CreateWorker 调用得到新实例。
type WorkerCreateFn func() Worker

// Server 是进程级单例，统一管理本机所有 Service、与 Center 通信、
// 维护集群 Actor 分配视图、提供本地/远程路由判定与热更新入口。
//
// 启动流程：SetConfig -> 配置 ServerConf.RegisterServices/WorkerCreateFn -> Init -> Start。
type Server interface {
	// SetServerId 设置本进程 Server 在集群中的唯一 ID。
	SetServerId(id uint32)
	// GetServerId 返回本进程 Server 在集群中的唯一 ID。
	GetServerId() uint32
	// SetGroup 设置本进程所在 group 路由名；空字符串表示未配置，不会自动映射默认分组。
	SetGroup(group string)
	// GetGroup 返回本进程所在 group 路由名；空字符串表示未配置。
	GetGroup() string

	// AddService 注册业务 Service。
	//
	// 框架内置 NilService 和 RPC 路由 Service 在 Init 时自动注册，不需要业务介入。
	// RPC 路由 Service 没有对外接口，业务只应通过 Go/Call 系列 API 触发跨服调用。
	// 业务 Service 推荐在 ServerConf.RegisterServices 钩子中注册，确保内置服务已存在。
	AddService(typ uint32, service Service)
	// GetService 根据 Service 类型返回已注册的 Service。
	GetService(typ uint32) (Service, bool)

	// SetConfig 设置 Server 启动配置。
	SetConfig(conf *ServerConf)
	// Init 初始化 Server、内置 Service、网络监听和 Center 连接。
	Init() error
	// Start 启动 Server 的定时器、网络监听和所有已注册 Service。
	Start()
	// Stop 停止 Server 并关闭本进程管理的全部 Service 和连接。
	Stop()
	// IsStop 返回 Server 是否已进入停止状态。
	IsStop() bool
	// HotUpdate 用一组新的 Worker 工厂替换全部业务 Service 的 Worker。
	//
	// 仅 typSet 中的业务 Service 参与热更；每个目标 Actor 会创建新的 Worker，
	// 并通过 Actor.HotUpdateWorker 触发 newWorker.OnHotUpdate(oldWorker)。
	HotUpdate(cfg *ServerConf, typSet map[uint32]bool)
	// SyncCenterData 向 Center 同步本 Server 当前 Actor 运行态信息。
	SyncCenterData() error

	// CreateWorkshop 创建 Service 自身的 Worker，Service 本身也是 Actor。
	CreateWorkshop(typ uint32) (Worker, bool)
	// CreateWorker 为指定 Actor 创建普通业务 Worker。
	CreateWorker(typ uint32, id uint64) (Worker, bool)

	// GetSupportTypArr 返回本 Server 支持的 Service 类型列表，可通过 filterTypSet 过滤内部类型。
	GetSupportTypArr(filterTypSet map[uint32]bool) []uint32
	// GetLocalActorIdsByTyp 返回本 Server 指定 Service 下的本地 Actor ID 列表。
	GetLocalActorIdsByTyp(typ uint32) []uint64
	// GetActorMsgCntByTyp 返回本 Server 指定 Service 下每个 Actor 的待处理消息数。
	GetActorMsgCntByTyp(typ uint32) map[uint64]uint32
	// BroadcastLocalEvent 向本 Server 上全部业务 Actor 广播事件，返回投递数量。
	BroadcastLocalEvent(eventId uint32, eventData interface{}) int

	// SetClusterInfo 用从 Center 拿到的全量数据重置本地缓存的集群视图。
	//
	// 通常在 Server 重启恢复或 Center 主动推送时调用。
	SetClusterInfo(svrInfoArr []*ServerInfo, svrAssignInfoArr []*ServerAssignInfo)
	// GetActorAssignServerId 使用本机已配置 group 查询目标 Actor 应路由到哪台 Server。
	//
	// 本机 group 为空时返回错误；先查本地缓存，若未命中再请求 Center；bCheckOnly=true 时仅查询不分配。
	GetActorAssignServerId(typ uint32, id uint64, bCheckOnly bool) (uint32, error)
	// GetActorAssignServerIdInGroup 查询指定 group 下目标 Actor 应路由到哪台 Server；group 必须非空。
	GetActorAssignServerIdInGroup(group string, typ uint32, id uint64, bCheckOnly bool) (uint32, error)
	// IsLocalActor 使用本机已配置 group 判定目标 Actor 是否在本 Server。
	//
	// 本机 group 为空时返回 false；Service.IsChecked 为 false 时仅在 group 等于本机 group 时返回 true。
	IsLocalActor(typ uint32, id uint64) bool
	// IsLocalActorInGroup 判定指定 group 下的目标 Actor 是否在本 Server；group 必须非空。
	IsLocalActorInGroup(group string, typ uint32, id uint64) bool
	// SetServerAddr 设置本地缓存的 Server RPC 地址。
	SetServerAddr(serverId uint32, addr string)
	// GetServerAddr 查询指定 Server 的 RPC 地址，缓存未命中时可向 Center 请求。
	GetServerAddr(serverId uint32) (string, error)
	// RecordTimeout 记录一次 Actor 慢消息事件，用于 IsCongestion 拥塞判定。
	RecordTimeout(typ uint32, id uint64)
	// IsCongestion 在 CongestionClearSec 时间窗内累计慢消息达 CongestionTimeoutCnt 时返回 true。
	IsCongestion(typ uint32, id uint64) bool

	// OnUpdate 由 GlobalTimer 周期回调，根据时间跨度向 Service 派发各级 Tick。
	OnUpdate(dt int64)

	// OnCreatePackReader 为新的 RPC 连接创建拆包器；由框架网络层内部调用。
	OnCreatePackReader(mgr network.ConnMgr) network.PackReader
	// OnConnOpen 在 RPC 连接建立后触发；由框架网络层内部调用。
	OnConnOpen(c network.Conn) error
	// OnConnClose 在 RPC 连接关闭后触发；由框架网络层内部调用。
	OnConnClose(c network.Conn)
	// OnConnError 在 RPC 连接读写发生底层 IO 错误后触发；由框架网络层内部调用。
	OnConnError(c network.Conn)
	// OnConnRecv 在 RPC 连接收到完整网络包后触发；由框架网络层内部调用。
	OnConnRecv(c network.Conn, pack []byte)
	// ServeHTTP 处理 Server 暴露的 HTTP 状态请求。
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}
