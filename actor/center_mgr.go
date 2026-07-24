package actor

// ActorInfo Server 上报给 Center 的单个 Actor 运行态信息。
// 用于 Center 计算负载、判断拥塞、参与 Actor 分配决策。
type ActorInfo struct {
	// ActorId 表示 Actor ID。
	ActorId uint64 `json:"actor_id"`
	// MsgCnt 表示当前待处理消息数，包含读写队列中的消息。
	MsgCnt uint32 `json:"msg_cnt"`
	// IsCongestion 表示 Actor 是否拥塞，1 为拥塞，0 为正常。
	IsCongestion uint8 `json:"is_cong"`
}

// ServiceInfo Server 上报的单个 Service（按 typ 聚合）的 Actor 列表。
type ServiceInfo struct {
	// Typ 表示 Service 类型。
	Typ uint32 `json:"typ"`
	// Actors 表示该 Service 当前管理的 Actor 状态列表。
	Actors []*ActorInfo `json:"actors"`
}

// ServerInfo 是 Center 合并 MySQL 静态路由与 Server 运行态后维护的完整快照。
type ServerInfo struct {
	// ServerId 表示 Server 在集群中的唯一 ID。
	ServerId uint32 `json:"svr_id"`
	// SvrClusterId 表示 Server 所属集群 ID，由静态路由配置提供；一个集群对应一个 Center。
	SvrClusterId uint32 `json:"svr_cluster_id,omitempty"`
	// SvrGroupId 表示 Server 所属分组 ID。
	SvrGroupId uint32 `json:"svr_group_id,omitempty"`
	// Group 表示 Server 所属分组路由名；同组内 (typ, actorId) 唯一。
	Group string `json:"group,omitempty"`
	// ServerHost 表示 Server 对外暴露的主机地址，由静态路由配置提供。
	ServerHost string `json:"svr_host"`
	// ServerPort 表示 Server 的 RPC 监听端口，由静态路由配置提供。
	ServerPort uint32 `json:"svr_port"`
	// SupportTypArr 表示该 Server 支持承载的 Service 类型集合，由静态路由配置提供。
	SupportTypArr []uint32 `json:"support_typ"`
	// Typ2ServiceInfoMap 表示按 Service 类型聚合的运行状态。
	Typ2ServiceInfoMap map[uint32]*ServiceInfo `json:"typ_service_dict"`
	// UpdateTime 表示 Center 端记录的最近上报时间戳，超过 ServerActiveTimeoutSec 视为掉线。
	UpdateTime int64 `json:"update_time"`
}

// ActorAssignInfo Center 维护的某个 Actor 被分配到哪台 Server 的记录。
type ActorAssignInfo struct {
	// ServerId 表示 Actor 被分配到的 Server ID。
	ServerId uint32 `json:"svr_id"`
	// AssignTime 表示本次分配发生的时间戳。
	AssignTime int64 `json:"assign_time"`
}

// TypAssignInfo 同一 typ 下所有 ActorId 的分配映射。
type TypAssignInfo struct {
	// ActorId2AssignMap 表示 Actor ID 到分配记录的映射。
	ActorId2AssignMap map[uint64]*ActorAssignInfo `json:"actor_assign_dict"`
}

// AssignInfo Center 端全局视角的「typ -> ActorId -> Server」分配表。
type AssignInfo struct {
	// Typ2AssignMap 是旧版无 group 分配表，仅用于兼容旧快照；新逻辑使用 Group2TypAssignMap。
	Typ2AssignMap map[uint32]*TypAssignInfo `json:"typ_assign_dict"`
	// Group2TypAssignMap 表示 Group 到「typ -> ActorId -> Server」分配表的映射。
	Group2TypAssignMap map[string]map[uint32]*TypAssignInfo `json:"group_assign_dict,omitempty"`
}

// ActorIdSet Server 视角下某个 typ 拥有的 Actor 集合，value 为分配时间戳。
type ActorIdSet = map[uint64]int64

// ServerAssignInfo Center 端反向索引：某台 Server 上分配了哪些 typ 下的哪些 Actor。
// 用于 Server 重启后快速恢复分配状态，以及 Server 离线时清理对应分配。
type ServerAssignInfo struct {
	// ServerId 表示反向索引所属的 Server ID。
	ServerId uint32 `json:"svr_id"`
	// Group 表示反向索引所属分组；为空表示未配置，不会自动映射到默认分组。
	Group string `json:"group,omitempty"`
	// Typ2ActorIdsMap 表示 Service 类型到该 Server 已分配 Actor 集合的映射。
	Typ2ActorIdsMap map[uint32]ActorIdSet `json:"typ_actorids_dict"`
}

// SyncDataReq 表示 Server 向 Center 同步运行状态的请求体。
//
// Server 的 cluster、group、RPC 地址和支持类型属于静态路由配置，由各服务
// 从 MySQL 加载，不通过心跳上报。
type SyncDataReq struct {
	// ServerId 表示 Server 在集群中的唯一 ID。
	ServerId uint32 `json:"svr_id"`
	// Typ2ServiceInfoMap 表示按 Service 类型聚合的运行状态。
	Typ2ServiceInfoMap map[uint32]*ServiceInfo `json:"typ_service_dict"`
}

// SyncDataRsp 表示 Center 响应状态同步请求的空响应体。
type SyncDataRsp struct{}

// ActorServerIdReq 表示查询 Actor 已分配 Server 的请求体。
type ActorServerIdReq struct {
	// Group 表示目标分组路由名；调用方必须显式传入。
	Group string `json:"group,omitempty"`
	// Typ 表示目标 Actor 类型。
	Typ uint32 `json:"typ"`
	// Id 表示目标 Actor ID。
	Id uint64 `json:"id"`
}

// ActorServerIdRsp 表示查询 Actor 已分配 Server 的响应体。
type ActorServerIdRsp struct {
	// ServerId 表示承载目标 Actor 的 Server ID。
	ServerId uint32 `json:"svr_id"`
}

// ServerAssignReq 表示请求 Center 为 Actor 分配 Server 的请求体。
type ServerAssignReq struct {
	// Group 表示目标分组路由名；调用方必须显式传入。
	Group string `json:"group,omitempty"`
	// Typ 表示目标 Actor 类型。
	Typ uint32 `json:"typ"`
	// Id 表示目标 Actor ID。
	Id uint64 `json:"id"`
}

// ServerAssignRsp 表示 Actor 分配结果响应体。
type ServerAssignRsp struct {
	// ServerId 表示 Center 最终分配的 Server ID。
	ServerId uint32 `json:"svr_id"`
}
