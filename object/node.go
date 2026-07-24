package object

// Node 是 Component 的容器，组成树形结构。
//
// 设计要点：
//   - 一个 Node 可挂多个 Component，但 DbComponent / OnlineComponent / TimerComponent /
//     EventComponent 每类只挂一个（语义上的单例桶，便于框架直接派发）
//   - Handle* 系列方法把生命周期事件（加载、登录、上下线、各级时间变化）沿子树自顶向下递归
//   - HotUpdate 沿子树触发组件热更；HotUpdateFinish 是「全部组件已替换完」的二阶段通知
type Node interface {
	// SetId 设置 Node ID。
	SetId(id uint32)
	// GetId 返回 Node ID。
	GetId() uint32
	// GetName 返回 Node 名称。
	GetName() string
	// Init 初始化当前 Node 及其组件。
	Init() error
	// Destroy 销毁当前 Node、组件及子节点。
	Destroy()
	// HotUpdate 在 ObjectBuilder.GetHotVer 高于本节点版本时触发组件替换，返回 true 表示发生了更新。
	HotUpdate() bool
	// HotUpdateFinish 在 HotUpdate 后调用，给 Component 一个重建交叉引用的时机。
	HotUpdateFinish()

	// AddChild 添加子节点。
	AddChild(child Node)
	// RemoveChild 移除指定子节点。
	RemoveChild(child Node)
	// RemoveChildByName 按节点名称移除子节点。
	RemoveChildByName(name string)
	// GetChildren 返回当前 Node 的直接子节点列表。
	GetChildren() []Node
	// GetChild 按节点名称返回直接子节点。
	GetChild(name string) (Node, bool)
	// GetParent 返回父节点。
	GetParent() Node
	// SetParent 设置父节点。
	SetParent(parent Node)

	// SetDbComponent 设置数据组件，Node.LoadData 和 Node.SaveData 会派发到该组件。
	SetDbComponent(name string, component DbComponent)
	// SetOnlineComponent 设置在线状态组件，登录和上下线事件会派发到该组件。
	SetOnlineComponent(name string, component OnlineComponent)
	// SetTimerComponent 设置定时组件，时间变化事件会派发到该组件。
	SetTimerComponent(name string, component TimerComponent)
	// SetEventComponent 设置事件组件，Node.HandleEvent 会派发到该组件。
	SetEventComponent(name string, component Component)
	// AddComponent 通用挂载入口。同名 Component 会被先移除再添加；
	//
	// 若 Component 已属于其他 Node 则返回 false。
	AddComponent(name string, component Component) bool
	// RemoveComponent 按名称移除组件。
	RemoveComponent(name string)
	// GetComponent 按名称返回当前 Node 上挂载的组件。
	GetComponent(name string) (Component, bool)
	// FindComponentByPath 按 "/" 分隔的 Node 名路径找到节点，再返回该节点上的指定组件。
	FindComponentByPath(path string, name string) (Component, bool)

	// EmitEvent 从根节点开始广播事件，即向上找到 root 后调用 HandleEvent。
	EmitEvent(event string, data interface{})
	// HandleEvent 先让当前节点的事件组件处理事件，然后递归派发给子节点。
	HandleEvent(event string, data interface{})

	// LoadData 沿当前子树加载持久化数据。
	LoadData(dbData interface{}, dataVer int64) error
	// SaveData 沿当前子树保存持久化数据。
	SaveData(dbData interface{}) error

	// HandleLoginSuccess 沿当前子树派发登录成功事件。
	HandleLoginSuccess()
	// HandleOnline 沿当前子树派发上线事件。
	HandleOnline()
	// HandleOffline 沿当前子树派发离线事件。
	HandleOffline()

	// HandleUpdate 沿当前子树派发定时更新事件。
	HandleUpdate(diff int32)
	// HandleSecondChange 沿当前子树派发自然秒变化事件。
	HandleSecondChange()
	// HandleMinuteChange 沿当前子树派发自然分钟变化事件。
	HandleMinuteChange()
	// HandleHourChange 沿当前子树派发自然小时变化事件。
	HandleHourChange()
	// HandleDayChange 沿当前子树派发自然日期变化事件。
	HandleDayChange()
	// HandleWeekChange 沿当前子树派发自然周变化事件。
	HandleWeekChange()
	// HandleMonthChange 沿当前子树派发自然月份变化事件。
	HandleMonthChange()
}

// CmdRouter cmd 字符串 -> (Node 路径, Component 名) 的全局路由表。
// 业务启动时统一注册路由，Component.SendCmd 借此找到目标。
type CmdRouter interface {
	// SetCmdPathAndCom 设置命令到 Node 路径和 Component 名称的映射。
	SetCmdPathAndCom(cmd string, path string, name string)
	// GetCmdPathAndCom 返回命令对应的 Node 路径和 Component 名称。
	GetCmdPathAndCom(cmd string) (string, string, bool)
}
