package actor

import "time"

// GlobalTimer 进程级全局时钟。
// 由 Server 在 Start 时按 ServerConf.TimerIntvMs 启动；每次 tick 驱动 Server.OnUpdate，
// 进而向各 Service 投递 MsgIdTick 触发跨秒/分/时/天等回调。
//
// Now 与 GetSystemTime 的区别：
//   - GetSystemTime 返回真实系统时间（不含 offset），用于网络包时间戳等
//   - Now 返回叠加了 SetOffset 偏移后的时间，便于业务做时间穿越调试或同步逻辑时钟
type GlobalTimer interface {
	// SetOffset 设置逻辑时钟偏移（秒），影响 Now、NowUnix 和 NowUnixNano。
	SetOffset(offset int64)
	// Start 启动时钟循环，interval 通常对应 ServerConf.TimerIntvMs。
	Start(interval time.Duration)
	// Stop 停止时钟循环。
	Stop()
	// GetSystemTime 返回真实系统时间，不叠加逻辑偏移。
	GetSystemTime() time.Time
	// GetSystemUnix 返回真实系统 Unix 秒时间戳，不叠加逻辑偏移。
	GetSystemUnix() int64
	// GetSystemUnixNano 返回真实系统 Unix 纳秒时间戳，不叠加逻辑偏移。
	GetSystemUnixNano() int64
	// Now 返回叠加逻辑偏移后的当前时间。
	Now() time.Time
	// NowUnix 返回叠加逻辑偏移后的 Unix 秒时间戳。
	NowUnix() int64
	// NowUnixNano 返回叠加逻辑偏移后的 Unix 纳秒时间戳。
	NowUnixNano() int64
}
