package object

// ObjectCreateFn 对象工厂函数，由 ObjectBuilder.RegisterObjectCreateFn 注册。
type ObjectCreateFn func() Object

// ObjectHotUpdateFn 热更转换函数：把旧版本对象的状态迁移到新版本对象。
// 实现需要把 obj 的关键字段拷贝/转换到新实例上并返回。
type ObjectHotUpdateFn func(obj Object) Object

// Object 所有可热更对象的基础接口。
// 框架基于 ObjectBuilder 注册的 (cls, createFn, hotUpdateFn, ver) 实现运行时升级。
//
// SetSelf/GetSelf 与 GetBase 是为了支持子类化：子类在调用父类方法时，
// 父类方法需要拿到「真正的实例」（self）才能调到子类重写的钩子。
type Object interface {
	// SetClass 设置对象类名。
	SetClass(typ string)
	// GetClass 返回对象类名。
	GetClass() string

	// SetVer 设置当前实例的热更版本号。
	// 当 ObjectBuilder.GetHotVer 大于此值时，调用 HotUpdate 会触发实际转换。
	SetVer(ver int32)
	// GetVer 返回当前实例的热更版本号。
	GetVer() int32

	// SetSelf 设置真实对象实例引用，用于基础类型回调到子类型实现。
	SetSelf(o Object)
	// GetSelf 返回真实对象实例引用。
	GetSelf() Object

	// SetBase 设置本对象的父类实例，用于子类调用父类方法。
	SetBase(o Object)
	// GetBase 返回本对象的父类实例。
	GetBase() (Object, bool)

	// OnHotUpdate 在版本从 oldVer 升到 newVer 后调用，业务可执行一次性迁移逻辑。
	OnHotUpdate(oldVer int32, newVer int32)
}

// ObjectBuilder 全局对象工厂注册中心 + 热更管理。
//
// 典型用法：
//  1. 业务启动时通过 RegisterObjectCreateFn 注册所有可创建的类
//  2. 发起热更时调用 HotUpdate(cls2HotUpdateFnMap, ver) 注入新版工厂
//  3. 框架后续在 CreateObject / Node.HotUpdate 中按最新 hotVer 应用转换
type ObjectBuilder interface {
	// RegisterObjectCreateFn 注册指定类名的对象创建函数。
	RegisterObjectCreateFn(cls string, createFunc ObjectCreateFn)
	// SetBaseVer 设置初始版本号；新建对象时以此为 ver。
	SetBaseVer(ver int32)
	// GetBaseVer 返回新建对象使用的初始版本号。
	GetBaseVer() int32
	// GetHotVer 返回当前已注入的最新热更版本。
	GetHotVer() int32
	// HotUpdate 注入一批新的热更函数；ver 必须严格大于当前热更版本，否则忽略。
	HotUpdate(cls2HotUpdateFnMap map[string]ObjectHotUpdateFn, ver int32)
	// GetObjectFn 返回对象创建函数、热更函数和目标版本号。
	GetObjectFn(cls string) (ObjectCreateFn, ObjectHotUpdateFn, int32, bool)
}
