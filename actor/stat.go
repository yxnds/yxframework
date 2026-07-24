package actor

// SvrStatInfo 表示 Server 对外暴露的状态摘要。
type SvrStatInfo struct {
	// SvrId 表示 Server 在集群中的唯一 ID。
	SvrId uint32 `json:"svr_id"`
	// MainVer 表示主程序版本号。
	MainVer string `json:"main_ver"`
	// HotVer 表示当前热更版本号。
	HotVer string `json:"hot_ver"`
	// ActorCnt 表示本 Server 当前承载的 Actor 总数。
	ActorCnt uint32 `json:"actor_cnt"`
	// MsgCnt 表示本 Server 当前待处理的消息总数。
	MsgCnt uint32 `json:"msg_cnt"`
}
