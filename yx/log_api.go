package yx

import "time"

// LogLv 表示日志级别。
type LogLv = int

const (
	// LOG_LV_DEBUG 表示调试级别日志，是最详细的日志级别。
	LOG_LV_DEBUG LogLv = 0
	// LOG_LV_INFO 表示普通信息级别日志。
	LOG_LV_INFO LogLv = 1
	// LOG_LV_WARN 表示警告级别日志。
	LOG_LV_WARN LogLv = 2
	// LOG_LV_ERROR 表示错误级别日志，是默认最高严重级别。
	LOG_LV_ERROR LogLv = 3
)

const (
	// LOG_BATCH_DUMP_COUNT 表示日志批量写入的默认条数。
	LOG_BATCH_DUMP_COUNT = 256
	// LOG_MAX_CACHE_SIZE 表示日志异步队列的默认最大缓存条数。
	LOG_MAX_CACHE_SIZE = 8192
	// LOG_DEFAULT_DUMP_SIZE 表示默认日志文件滚动大小，单位字节。
	LOG_DEFAULT_DUMP_SIZE = 64 * 1024 * 1024
	// LOG_DEFAULT_DUMP_THRESHOLD 表示默认触发刷盘的缓存阈值。
	LOG_DEFAULT_DUMP_THRESHOLD = 256
	// LOG_DEFAULT_DUMP_INTV 表示默认定时刷盘间隔，单位毫秒。
	LOG_DEFAULT_DUMP_INTV = 100
	// LOG_DEFAULT_FILE_BUF_SIZE 表示默认文件写缓冲大小，单位字节。
	LOG_DEFAULT_FILE_BUF_SIZE = 256 * 1024
)

// LOG_DEBUG_SWITCH_FILE 表示运行期切换调试日志输出的开关文件名。
const LOG_DEBUG_SWITCH_FILE = "debug_log"

// LOG_CONSOLE_SWITCH_FILE 表示运行期切换控制台日志输出的开关文件名。
const LOG_CONSOLE_SWITCH_FILE = "print_log"

// QueueMode 控制异步日志队列满时的处理方式。
type QueueMode uint8

const (
	// QueueModeBlock 阻塞写日志方，直到后台写协程追上进度。
	QueueModeBlock QueueMode = iota
	// QueueModeDropNewest 在队列满时丢弃新日志。
	QueueModeDropNewest
	// QueueModeDropOldest 在队列满时丢弃一条旧日志，为新日志腾空间。
	QueueModeDropOldest
)

// Format 选择文件和控制台日志行格式。
type Format uint8

const (
	// FormatText 输出人类可读的文本日志。
	FormatText Format = iota
	// FormatJSON 每行输出一个 JSON 对象。
	FormatJSON
)

// RotationPeriod 控制日志文件按本地日历时间轮动，空值禁用时间轮动。
type RotationPeriod string

const (
	RotationNone   RotationPeriod = ""
	RotationHourly RotationPeriod = "hourly"
	RotationDaily  RotationPeriod = "daily"
	RotationWeekly RotationPeriod = "weekly"
)

// LogConf 日志配置。
// Dump 系列字段控制日志落盘的滚动策略（按文件大小或时间间隔切分）。
type LogConf struct {
	// Level 表示日志过滤级别，低于该级别的日志会被忽略。
	Level int `json:"level"`
	// IsShowCaller 表示是否在日志中打印调用方文件和行号。
	IsShowCaller bool `json:"is_show_caller"`
	// DisableFile 表示是否禁用日志文件输出。
	DisableFile bool `json:"disable_file"`
	// IsPrintConsole 表示是否打印到控制台。
	IsPrintConsole bool `json:"is_print_console"`
	// DumpPath 表示日志文件输出路径。
	DumpPath string `json:"dump_path"`
	// DumpFileSize 表示单个日志文件最大字节数，超过后滚动文件。
	DumpFileSize int `json:"dump_file_size"`
	// RotationPeriod 支持 hourly、daily、weekly；空值仅按大小轮动。
	RotationPeriod RotationPeriod `json:"rotation_period"`
	// RotationHour 指定 daily/weekly 模式的本地轮动小时（0–23），默认 0。
	RotationHour int `json:"rotation_hour"`
	// DumpThreshold 表示触发强制刷盘的缓冲条数阈值。
	DumpThreshold int `json:"dump_threshold"`
	// DumpInterval 表示定期刷盘间隔，单位毫秒。
	DumpInterval uint32 `json:"dump_interval"`
}

// LogOptions 保存单个日志器的运行时配置。
//
// 零值会在 Configure 中规范化，调用方只需要设置自己关心的字段。
type LogOptions struct {
	// Level 表示日志过滤级别。
	Level LogLv
	// ShowCaller 表示是否输出调用方文件和行号。
	ShowCaller bool
	// Console 表示是否输出到控制台。
	Console bool
	// PrintFunc 表示自定义日志输出函数。
	PrintFunc func(lv LogLv, logStr string)
	// DisableFile 表示禁用文件输出；DumpPath 非空时默认仍会输出到文件。
	DisableFile bool

	// DumpPath 表示日志文件输出路径。
	DumpPath string
	// MaxFileSize 表示单个日志文件最大字节数，超过后滚动文件。
	MaxFileSize int64
	// RotationPeriod 按本地时间轮动：整点、每天或每周一 RotationHour 点。
	// 到期后下一次写入时轮动，与 MaxFileSize 同时生效。
	RotationPeriod RotationPeriod
	// RotationHour 用于 daily/weekly 模式，取值 0–23，默认 0。
	RotationHour int
	// FlushThreshold 表示触发刷盘的缓存条数阈值。
	FlushThreshold int
	// FlushInterval 表示定期刷盘间隔。
	FlushInterval time.Duration
	// QueueSize 表示异步日志队列容量。
	QueueSize int
	// QueueMode 表示异步日志队列满时的处理方式。
	QueueMode QueueMode
	// BatchSize 表示后台批量写入的最大条数。
	BatchSize int
	// FileBufferSize 表示日志文件写缓冲大小。
	FileBufferSize int
	// Format 表示日志行格式。
	Format Format

	// DebugSwitchFile 表示运行期切换调试日志输出的开关文件名。
	DebugSwitchFile string
	// ConsoleSwitchFile 表示运行期切换控制台输出的开关文件名。
	ConsoleSwitchFile string

	// EnableStackTrace 表示是否在高等级日志中附带调用栈。
	EnableStackTrace bool
	// StacktraceLevel 表示触发调用栈输出的最低日志级别。
	StacktraceLevel LogLv
	// SyncOnFlush 表示 Flush 时是否同步落盘。
	SyncOnFlush bool
}

// Options 是 LogOptions 的兼容别名。
type Options = LogOptions

// LogStats 是日志计数器和队列状态的瞬时快照。
type LogStats struct {
	// Enqueued 表示已入队的日志条数。
	Enqueued uint64
	// Dropped 表示因队列策略丢弃的日志条数。
	Dropped uint64
	// Written 表示已写出的日志条数。
	Written uint64

	// WriteErrors 表示日志写出失败次数。
	WriteErrors uint64
	// Rotations 表示日志文件滚动次数。
	Rotations uint64
	// Flushes 表示刷盘次数。
	Flushes uint64
	// PrintPanics 表示自定义输出函数发生 panic 的次数。
	PrintPanics uint64

	// QueueLen 表示当前异步队列长度。
	QueueLen int
	// QueueCap 表示当前异步队列容量。
	QueueCap int
	// Started 表示日志器是否已经启动。
	Started bool
	// Stopped 表示日志器是否已经停止。
	Stopped bool
}

// Stats 是 LogStats 的兼容别名。
type Stats = LogStats

// Logger 普通日志接口。
// D/I/W/E 分别对应 Debug/Info/Warn/Error；参数风格类似 fmt.Println，会自动加空格分隔。
// Detail 用于批量输出多行日志（一次 flush，减少锁竞争）。
type Logger interface {
	// D 输出 Debug 级别日志。
	D(a ...interface{})
	// I 输出 Info 级别日志。
	I(a ...interface{})
	// W 输出 Warn 级别日志。
	W(a ...interface{})
	// E 输出 Error 级别日志。
	E(a ...interface{})
	// Debugf 按格式化字符串输出 Debug 级别日志。
	Debugf(format string, a ...interface{})
	// Infof 按格式化字符串输出 Info 级别日志。
	Infof(format string, a ...interface{})
	// Warnf 按格式化字符串输出 Warn 级别日志。
	Warnf(format string, a ...interface{})
	// Errorf 按格式化字符串输出 Error 级别日志。
	Errorf(format string, a ...interface{})
	// Detail 批量输出指定级别的多行日志。
	Detail(lv LogLv, logs [][]interface{})
	// Ln 输出空行。
	Ln()
}

// IndependentLogger 独立 logger，不受全局 LogMgr 设置影响。
// 用于需要单独输出到特定文件 / 不同级别的场景（如行为日志、计费日志）。
type IndependentLogger interface {
	// Logger 让 IndependentLogger 具备普通日志输出能力。
	Logger
	// StartLogger 启动独立日志器。
	StartLogger()
	// StopLogger 停止独立日志器。
	StopLogger()
	// Stop 停止独立日志器并返回停止过程中的错误。
	Stop() error
	// Flush 立即刷新独立日志器缓存。
	Flush() error
	// SetLogLevel 设置独立日志器过滤级别。
	SetLogLevel(lv LogLv)
	// SetShowCaller 设置是否输出调用方文件和行号。
	SetShowCaller(bShowCaller bool)
	// SetPrintFunc 替换底层输出函数，可改为推送到日志收集系统。
	SetPrintFunc(printFunc func(lv LogLv, logStr string))
	// SetConsole 设置是否输出到控制台。
	SetConsole(enabled bool)
	// SetFormat 设置日志行格式。
	SetFormat(format Format)
	// Configure 使用完整选项配置独立日志器。
	Configure(opts Options) error
	// Options 返回独立日志器当前配置快照。
	Options() Options
	// Stats 返回独立日志器运行统计快照。
	Stats() Stats
	// ConfigLogger 使用兼容配置初始化独立日志器。
	ConfigLogger(cfg *LogConf, printFunc func(lv LogLv, logStr string)) error
}

// LogMgr 全局日志管理器。
// 通过 NewLogger(tag) 创建的所有 Logger 共享此处的配置。
type LogMgr interface {
	// StartLogger 启动全局日志管理器。
	StartLogger()
	// StopLogger 停止全局日志管理器。
	StopLogger()
	// Stop 停止全局日志管理器并返回停止过程中的错误。
	Stop() error
	// Flush 立即刷新全局日志缓存。
	Flush() error
	// SetLogLevel 设置全局日志过滤级别。
	SetLogLevel(lv LogLv)
	// SetShowCaller 设置是否输出调用方文件和行号。
	SetShowCaller(bShowCaller bool)
	// LogArgs 将变长参数转成切片，便于在 Detail 等接口中拼接。
	LogArgs(a ...interface{}) []interface{}
	// SetPrintFunc 替换全局日志底层输出函数。
	SetPrintFunc(printFunc func(lv LogLv, logStr string))
	// SetConsole 设置全局日志是否输出到控制台。
	SetConsole(enabled bool)
	// SetFormat 设置全局日志行格式。
	SetFormat(format Format)
	// Configure 使用完整选项配置全局日志管理器。
	Configure(opts Options) error
	// Options 返回全局日志管理器当前配置快照。
	Options() Options
	// Stats 返回全局日志管理器运行统计快照。
	Stats() Stats
	// ConfigLogger 使用兼容配置初始化全局日志管理器。
	ConfigLogger(cfg *LogConf, printFunc func(lv LogLv, logStr string)) error
}
