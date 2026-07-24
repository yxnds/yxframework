package actor

import "errors"

var (
	// ErrNoWorker 表示 Actor.Init 时尚未绑定 Worker。
	ErrNoWorker = errors.New("no worker for actor")
	// ErrNotServiceActor 表示将非 Service Actor 当作 Service 使用。
	ErrNotServiceActor = errors.New("actor:not service actor")
	// ErrPckHeadTooSmall 表示 Actor 网络包长度小于 ActorConnHeadLen。
	ErrPckHeadTooSmall = errors.New("actor:pack head is too small")
	// ErrPckLenTooLarge 表示 Actor 网络包声明的包体长度超过 ActorConnMaxPackSize。
	ErrPckLenTooLarge = errors.New("actor:pack head is too large")
	// ErrGetServerAddrFail 表示无法从本地缓存或 Center 查询到 Server 地址。
	ErrGetServerAddrFail = errors.New("actor:get server addr fail")
	// ErrCenterUnavailable 表示 Center 不可用、连接失败或响应异常。
	ErrCenterUnavailable = errors.New("actor:center unavailable")
	// ErrInvalidGroup 表示路由 group 为空或非法。
	ErrInvalidGroup = errors.New("actor:invalid group")
	// ErrNoServer 表示 Center 没有找到可承载目标 Actor 的在线 Server。
	ErrNoServer = errors.New("actor:no server")
	// ErrActorAssignFailed 表示 Center 为 Actor 分配 Server 失败。
	ErrActorAssignFailed = errors.New("actor:assign actor failed")
)
