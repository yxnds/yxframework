package object

// CmdHandleFn cmd 请求-响应处理器。
// cmd 是字符串路由，data 为业务自定义；返回值传给调用方，bool 标识处理是否成功。
type CmdHandleFn func(cmd string, data interface{}) (interface{}, bool)

// EventHandleFn 事件广播处理器（与 actor 包的 EventHandleFn 不同，这里是 object 域内事件）。
type EventHandleFn func(event string, data interface{})

// Component 挂载到 Node 上的功能单元。
// 一个 Node 可挂多个 Component，每个 Component 负责一类功能（DB、登录态、定时、命令、事件）。
//
// 命令 vs 事件：
//   - Cmd 是定向请求，有返回值，通过 CmdRouter 路由到具体 (path, comName)
//   - Event 是广播，无返回值，沿 Node 树自顶向下分发
type Component interface {
	// Object 让 Component 具备热更对象的基础能力。
	Object

	// SetName 设置 Component 名称。
	SetName(name string)
	// GetName 返回 Component 名称。
	GetName() string

	// SetNode 设置 Component 所属的 Node。
	SetNode(node Node)
	// GetNode 返回 Component 所属的 Node。
	GetNode() Node
	// GetNodeId 返回所属 Node 的 ID；Node 为空时返回 0。
	GetNodeId() uint32

	// RegisterCmdHandler 注册命令处理器。
	RegisterCmdHandler(cmd string, handler CmdHandleFn)
	// HandleCmd 执行指定命令并返回处理结果。
	HandleCmd(cmd string, data interface{}) (interface{}, bool)
	// SendCmd 经 CmdRouter 路由后调用目标 Component 的 HandleCmd。
	SendCmd(cmd string, data interface{}) (interface{}, bool)

	// RegisterEventHandler 注册事件处理器。
	RegisterEventHandler(event string, handler EventHandleFn)
	// HandleEvent 执行指定事件处理器。
	HandleEvent(event string, data interface{})
	// SendEvent 等价于 GetNode().EmitEvent，从根节点开始广播事件。
	SendEvent(event string, data interface{})

	// OnInit 在 Component 初始化时调用。
	OnInit() error
	// OnDestroy 在 Component 销毁时调用。
	OnDestroy()
	// OnNodeHotUp 在所属 Node 完成热更（包含子 Component 替换）后被调用，
	//
	// 业务可在此处重建对其他 Component 的引用。
	OnNodeHotUp()
}

// DbComponent 负责持久化的 Component 细分类型。
// 由 Node.SetDbComponent 设置后，Node.LoadData/SaveData 会自顶向下沿树调用 OnLoad/OnSave。
type DbComponent interface {
	// Component 让 DbComponent 具备普通 Component 的能力。
	Component

	// OnLoad 加载持久化数据并初始化 Component。
	OnLoad(dbData interface{}, dataVer int64) error
	// OnSave 将 Component 状态保存到持久化数据对象。
	OnSave(dbData interface{}) error
}

// OnlineComponent 在线状态 Component。
// Node.HandleLoginSuccess/HandleOnline/HandleOffline 会派发到所有 OnlineComponent。
type OnlineComponent interface {
	// Component 让 OnlineComponent 具备普通 Component 的能力。
	Component

	// OnLoginSuccess 在登录成功后调用。
	OnLoginSuccess()
	// OnOnline 在对象上线时调用。
	OnOnline()
	// OnOffline 在对象离线时调用。
	OnOffline()
}

// TimerComponent 时间驱动 Component。
// Node.HandleUpdate / HandleSecondChange / ... / HandleMonthChange 会派发到所有 TimerComponent。
// diff 为距上次 OnUpdate 的毫秒数。
type TimerComponent interface {
	// Component 让 TimerComponent 具备普通 Component 的能力。
	Component

	// OnUpdate 在定时更新时调用，diff 为距离上次调用的毫秒数。
	OnUpdate(diff int32)
	// OnSecondChange 在自然秒变化时调用。
	OnSecondChange()
	// OnMinuteChange 在自然分钟变化时调用。
	OnMinuteChange()
	// OnHourChange 在自然小时变化时调用。
	OnHourChange()
	// OnDayChange 在自然日期变化时调用。
	OnDayChange()
	// OnWeekChange 在自然周变化时调用。
	OnWeekChange()
	// OnMonthChange 在自然月份变化时调用。
	OnMonthChange()
}
