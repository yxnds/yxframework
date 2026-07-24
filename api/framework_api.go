package api

import (
	"net/http"
	"time"

	"github.com/yxnds/yxframework/actor"
	"github.com/yxnds/yxframework/db"
	"github.com/yxnds/yxframework/network"
	"github.com/yxnds/yxframework/object"
	"github.com/yxnds/yxframework/yx"

	"github.com/gorilla/websocket"
)

// FrameworkApi 是框架的统一聚合入口，把 actor / object / network / db / yx 五个子包的
// 工厂方法和管理器收拢在一个接口里，业务侧只需要持有它就能拿到框架的全部能力。
//
// 实现在独立仓库（如 yxcore）中，业务通过 GetFrameworkApi() 注入。
type FrameworkApi interface {
	// GetMysqlMgr 返回 MySQL 连接管理器。
	GetMysqlMgr() db.MysqlMgr
	// GetMysqlHandle 返回无状态 MySQL 操作集合。
	GetMysqlHandle() db.MysqlHandle
	// GetMongoMgr 返回 MongoDB 连接管理器。
	GetMongoMgr() db.MongoMgr
	// GetMongoHandle 返回无状态 MongoDB 操作集合。
	GetMongoHandle() db.MongoHandle
	// GetRedisMgr 返回 Redis 连接管理器。
	GetRedisMgr() db.RedisMgr
	// GetRedisHandle 返回无状态 Redis 操作集合。
	GetRedisHandle() db.RedisHandle

	// NewSockMgr 创建 TCP 长连接管理器。
	NewSockMgr(writeQueueSize int32, writeTimeoutSec int64) network.ConnMgr
	// NewWebsockMgr 创建 WebSocket 长连接管理器。
	NewWebsockMgr(writeQueueSize int32, writeTimeoutSec int64, readLimit int64, upgrader *websocket.Upgrader) network.ConnMgr
	// GetDefaultWebsockUpgrader 返回默认 WebSocket 升级器。
	GetDefaultWebsockUpgrader() *websocket.Upgrader
	// NewHttpGate 创建 HTTP 网关。
	NewHttpGate(bAllowOrigin bool, handler http.Handler) network.ConnGate

	// NewObjectBase 创建基础 Object 实例。
	NewObjectBase() object.Object
	// GetObjectBuilder 返回全局对象构建器。
	GetObjectBuilder() object.ObjectBuilder
	// NewComponentBase 创建基础 Component 实例。
	NewComponentBase() object.Component
	// NewNode 创建指定名称的 Node。
	NewNode(name string) object.Node
	// GetCmdRouter 返回全局命令路由器。
	GetCmdRouter() object.CmdRouter
	// CreateObject 通过已注册的 ObjectCreateFn 创建对象实例。
	CreateObject(cls string) (object.Object, bool)
	// HotUpdate 对单个 Object 应用最新热更逻辑，返回热更后的实例。
	HotUpdate(obj object.Object) object.Object

	// GetGlobalTimer 返回进程级全局时钟。
	GetGlobalTimer() actor.GlobalTimer
	// GetServer 返回进程级 Actor Server 单例。
	GetServer() actor.Server
	// NewActorBase 创建基础 Actor 实例。
	NewActorBase(typ uint32, id uint64) actor.Actor
	// NewServiceBase 创建基础 Service 实例。
	NewServiceBase(typ uint32) actor.Service
	// NewWorkerBase 创建基础 Worker 实例。
	NewWorkerBase() actor.Worker
	// NewWorkshopBase 创建基础 Service Worker 实例。
	NewWorkshopBase() actor.Worker
	// NewRpcWorker 创建框架内置 RPC 路由 Actor 的 Worker 实例，业务通常不需要直接调用。
	NewRpcWorker() actor.Worker
	// NewRpcWorkshop 创建框架内置 RPC 路由 Service 的 Worker 实例，业务通常不需要直接调用。
	NewRpcWorkshop() actor.Worker
	// NewActorConnReader 创建 Actor RPC 网络包拆包器。
	NewActorConnReader() network.PackReader
	// GetServerConf 返回当前 Server 配置副本；修改后需通过 Server.SetConfig 显式回写。
	GetServerConf() *actor.ServerConf
	// CloneServerConf 返回当前 Server 配置副本。
	CloneServerConf() *actor.ServerConf
	// Go1 以无来源 Actor 信息异步投递消息，内部会自动判定本地调用或 RPC。
	Go1(msgId uint32, data interface{}, toTyp uint32, toId uint64, traceId uint32)
	// GoBytes1 以无来源 Actor 信息异步投递已序列化的 RPC 字节；本地目标会先反序列化。
	GoBytes1(msgId uint32, data []byte, toTyp uint32, toId uint64, traceId uint32)
	// Go2 以指定来源 Actor 信息异步投递消息。
	Go2(msgId uint32, data interface{}, fromTyp uint32, fromId uint64, toTyp uint32, toId uint64, traceId uint32)
	// GoBytes2 以指定来源 Actor 信息异步投递已序列化的 RPC 字节。
	GoBytes2(msgId uint32, data []byte, fromTyp uint32, fromId uint64, toTyp uint32, toId uint64, traceId uint32)
	// GoGroup1 向指定 group 内 Actor 异步投递消息。
	GoGroup1(group string, msgId uint32, data interface{}, toTyp uint32, toId uint64, traceId uint32)
	// GoGroupBytes1 向指定 group 内 Actor 异步投递已序列化的 RPC 字节。
	GoGroupBytes1(group string, msgId uint32, data []byte, toTyp uint32, toId uint64, traceId uint32)
	// GoGroup2 向指定 group 内 Actor 异步投递消息，并携带来源 Actor 信息。
	GoGroup2(group string, msgId uint32, data interface{}, fromTyp uint32, fromId uint64, toTyp uint32, toId uint64, traceId uint32)
	// GoGroupBytes2 向指定 group 内 Actor 异步投递已序列化的 RPC 字节，并携带来源 Actor 信息。
	GoGroupBytes2(group string, msgId uint32, data []byte, fromTyp uint32, fromId uint64, toTyp uint32, toId uint64, traceId uint32)
	// Call1 以无来源 Actor 信息同步投递消息，并返回结果数据和结果码。
	Call1(waiter actor.MsgWaiter, msgId uint32, data interface{}, toTyp uint32, toId uint64, traceId uint32, timeout time.Duration) (interface{}, int32)
	// CallBytes1 以无来源 Actor 信息同步投递已序列化的 RPC 字节。
	CallBytes1(waiter actor.MsgWaiter, msgId uint32, data []byte, toTyp uint32, toId uint64, traceId uint32, timeout time.Duration) (interface{}, int32)
	// Call2 以指定来源 Actor 信息同步投递消息，并返回结果数据和结果码。
	Call2(waiter actor.MsgWaiter, msgId uint32, data interface{}, fromTyp uint32, fromId uint64, toTyp uint32, toId uint64, traceId uint32, timeout time.Duration) (interface{}, int32)
	// CallBytes2 以指定来源 Actor 信息同步投递已序列化的 RPC 字节。
	CallBytes2(waiter actor.MsgWaiter, msgId uint32, data []byte, fromTyp uint32, fromId uint64, toTyp uint32, toId uint64, traceId uint32, timeout time.Duration) (interface{}, int32)
	// CallGroup1 向指定 group 内 Actor 同步投递消息，并返回结果数据和结果码。
	CallGroup1(group string, waiter actor.MsgWaiter, msgId uint32, data interface{}, toTyp uint32, toId uint64, traceId uint32, timeout time.Duration) (interface{}, int32)
	// CallGroupBytes1 向指定 group 内 Actor 同步投递已序列化的 RPC 字节。
	CallGroupBytes1(group string, waiter actor.MsgWaiter, msgId uint32, data []byte, toTyp uint32, toId uint64, traceId uint32, timeout time.Duration) (interface{}, int32)
	// CallGroup2 向指定 group 内 Actor 同步投递消息，并携带来源 Actor 信息。
	CallGroup2(group string, waiter actor.MsgWaiter, msgId uint32, data interface{}, fromTyp uint32, fromId uint64, toTyp uint32, toId uint64, traceId uint32, timeout time.Duration) (interface{}, int32)
	// CallGroupBytes2 向指定 group 内 Actor 同步投递已序列化的 RPC 字节，并携带来源 Actor 信息。
	CallGroupBytes2(group string, waiter actor.MsgWaiter, msgId uint32, data []byte, fromTyp uint32, fromId uint64, toTyp uint32, toId uint64, traceId uint32, timeout time.Duration) (interface{}, int32)
	// IsYearChanged 判断两个时间是否跨自然年。
	IsYearChanged(t1 time.Time, t2 time.Time) bool
	// IsMonthChanged 判断两个时间是否跨自然月。
	IsMonthChanged(t1 time.Time, t2 time.Time) bool
	// IsDayChanged 判断两个时间是否跨自然日。
	IsDayChanged(t1 time.Time, t2 time.Time) bool
	// IsHourChanged 判断两个时间是否跨自然小时。
	IsHourChanged(t1 time.Time, t2 time.Time) bool
	// IsMinuteChanged 判断两个时间是否跨自然分钟。
	IsMinuteChanged(t1 time.Time, t2 time.Time) bool
	// IsSecondChanged 判断两个时间是否跨自然秒。
	IsSecondChanged(t1 time.Time, t2 time.Time) bool
	// IsWeekChanged 判断两个时间是否跨自然周。
	IsWeekChanged(t1 time.Time, t2 time.Time) bool
	// GetWeekStart 返回指定时间所在自然周的起始时间。
	GetWeekStart(t time.Time) time.Time

	// GetLogMgr 返回全局日志管理器。
	GetLogMgr() yx.LogMgr
	// NewLogger 创建受全局日志管理器配置影响的普通日志器。
	NewLogger(tag string) yx.Logger
	// NewIndependentLogger 创建不受全局日志配置影响的独立日志器。
	NewIndependentLogger(tag string) yx.IndependentLogger
	// GetUtilMgr 返回工具方法管理器。
	GetUtilMgr() yx.UtilMgr
	// NewBuffFactory 创建按大小分桶的字节缓冲池工厂。
	NewBuffFactory(minSize uint32, maxSize uint32, queueStep uint32) yx.BuffFactory
	// NewEvent 创建一次性事件信号对象。
	NewEvent() yx.Event
	// NewFastLock 创建自旋互斥锁。
	NewFastLock() yx.FastLock
	// NewSyncDispatcher 创建同步消息派发器。
	NewSyncDispatcher(msgName string) yx.Dispatcher
	// NewAsyncDispatcher 创建异步消息派发器。
	NewAsyncDispatcher(msgName string, maxMsgBuffSize uint16) yx.Dispatcher
	// NewNotifyCenter 创建通知中心。
	NewNotifyCenter() yx.NotifyCenter
	// NewLinkedQueue 创建非线程安全链表队列。
	NewLinkedQueue() yx.Queue
	// NewSyncLinkedQueue 创建线程安全链表队列。
	NewSyncLinkedQueue() yx.Queue
	// NewSyncLimitQueue 创建线程安全限长队列。
	NewSyncLimitQueue(maxSize uint64) yx.Queue
	// NewSet 按指定底层键类型创建集合。
	NewSet(t yx.SetType) yx.Set
	// NewIntSet 创建整数集合。
	NewIntSet() yx.Set
	// NewUintSet 创建无符号整数集合。
	NewUintSet() yx.Set
	// NewFloatSet 创建浮点数集合。
	NewFloatSet() yx.Set
	// NewObjectSet 创建可比较对象集合。
	NewObjectSet() yx.Set
	// NewSimpleBuffer 创建指定容量的简单缓冲区。
	NewSimpleBuffer(capacity int) yx.SimpleBuffer
	// GetSnowflake 返回全局 Snowflake ID 生成器。
	GetSnowflake() yx.Snowflake
}
