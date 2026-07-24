package yx

import (
	"crypto"
	"crypto/rsa"
	"crypto/tls"
	"hash"
	"io"
	"net/http"
	"net/url"
	"os/exec"
	"reflect"
	"strings"
	"time"

	"gopkg.in/gomail.v2"
)

// ElliptcType ECDSA 椭圆曲线类型。位长越大安全性越高、签名越慢、占用更大。
type ElliptcType = int

const (
	// ELLIPTIC_TYPE_P521 表示 P-521 曲线，安全级别最高但计算成本也最高。
	ELLIPTIC_TYPE_P521 ElliptcType = iota
	// ELLIPTIC_TYPE_P384 表示 P-384 曲线。
	ELLIPTIC_TYPE_P384
	// ELLIPTIC_TYPE_P256 表示 P-256 曲线。
	ELLIPTIC_TYPE_P256
	// ELLIPTIC_TYPE_P224 表示 P-224 曲线。
	ELLIPTIC_TYPE_P224
)

// EcdsaHash 表示 ECDSA 签名使用的哈希算法。
type EcdsaHash = crypto.Hash

const (
	// EcdsaHashSHA256 表示 ECDSA 使用 SHA-256 哈希算法。
	EcdsaHashSHA256 EcdsaHash = crypto.SHA256
	// EcdsaHashSHA384 表示 ECDSA 使用 SHA-384 哈希算法。
	EcdsaHashSHA384 EcdsaHash = crypto.SHA384
	// EcdsaHashSHA512 表示 ECDSA 使用 SHA-512 哈希算法。
	EcdsaHashSHA512 EcdsaHash = crypto.SHA512
)

// RsaHash 表示 RSA 签名或 OAEP 加解密使用的哈希算法。
type RsaHash = crypto.Hash

const (
	// RsaHashSHA256 表示 RSA 使用 SHA-256 哈希算法。
	RsaHashSHA256 RsaHash = crypto.SHA256
	// RsaHashSHA384 表示 RSA 使用 SHA-384 哈希算法。
	RsaHashSHA384 RsaHash = crypto.SHA384
	// RsaHashSHA512 表示 RSA 使用 SHA-512 哈希算法。
	RsaHashSHA512 RsaHash = crypto.SHA512
)

const (
	// CRC16_MODBUS 表示 CRC-16/MODBUS 多项式。
	CRC16_MODBUS uint16 = 0x8005
	// CRC16_IBM 表示 CRC-16/IBM 多项式。
	CRC16_IBM uint16 = 0x8005
	// CRC16_CCITT 表示 CRC-16/CCITT 多项式。
	CRC16_CCITT uint16 = 0x1021
)

var (
	// Crc16Modbus 是 CRC-16/MODBUS 参数集。
	Crc16Modbus = Crc16Params{
		Name:   "CRC-16/MODBUS",
		Poly:   CRC16_MODBUS,
		Init:   0xffff,
		RefIn:  true,
		RefOut: true,
		XorOut: 0x0000,
		Check:  0x4b37,
	}

	// Crc16ARC 是 CRC-16/ARC 参数集。
	Crc16ARC = Crc16Params{
		Name:   "CRC-16/ARC",
		Poly:   CRC16_IBM,
		Init:   0x0000,
		RefIn:  true,
		RefOut: true,
		XorOut: 0x0000,
		Check:  0xbb3d,
	}

	// Crc16CCITTFalse 是 CRC-16/CCITT-FALSE 参数集。
	Crc16CCITTFalse = Crc16Params{
		Name:   "CRC-16/CCITT-FALSE",
		Poly:   CRC16_CCITT,
		Init:   0xffff,
		RefIn:  false,
		RefOut: false,
		XorOut: 0x0000,
		Check:  0x29b1,
	}

	// Crc16XMODEM 是 CRC-16/XMODEM 参数集。
	Crc16XMODEM = Crc16Params{
		Name:   "CRC-16/XMODEM",
		Poly:   CRC16_CCITT,
		Init:   0x0000,
		RefIn:  false,
		RefOut: false,
		XorOut: 0x0000,
		Check:  0x31c3,
	}

	// Crc16X25 是 CRC-16/X-25 参数集。
	Crc16X25 = Crc16Params{
		Name:   "CRC-16/X-25",
		Poly:   CRC16_CCITT,
		Init:   0xffff,
		RefIn:  true,
		RefOut: true,
		XorOut: 0xffff,
		Check:  0x906e,
	}
)

// Crc16Params 描述一种 CRC-16 算法参数。
type Crc16Params struct {
	// Name 表示 CRC-16 算法名称。
	Name string
	// Poly 表示 CRC-16 多项式。
	Poly uint16
	// Init 表示 CRC-16 初始值。
	Init uint16
	// RefIn 表示输入字节是否按位反射。
	RefIn bool
	// RefOut 表示输出结果是否按位反射。
	RefOut bool
	// XorOut 表示输出结果的异或值。
	XorOut uint16
	// Check 表示标准测试串 "123456789" 的校验值。
	Check uint16
}

// Hash16 是 hash.Hash 风格的 16 位摘要接口。
type Hash16 interface {
	// Hash 让 Hash16 兼容标准库 hash.Hash 接口。
	hash.Hash
	// Sum16 返回当前 16 位摘要值。
	Sum16() uint16
}

// Crc16Table 是 CRC-16 查表缓存。
type Crc16Table [256]uint16

const (
	// BP_MIN_BUFF_SIZE 是缓冲池默认最小桶大小。
	BP_MIN_BUFF_SIZE = 64
	// BP_MAX_BUFF_SIZE 是缓冲池默认最大桶大小。
	BP_MAX_BUFF_SIZE = 4 * 1024
	// BP_QUEUE_STEP 是默认桶大小步进。
	BP_QUEUE_STEP = 64
)

// FastLockStats 是 FastLock 的观测快照。
type FastLockStats struct {
	// Locked 表示当前是否处于加锁状态。
	Locked bool
	// LockCount 表示成功加锁次数。
	LockCount uint64
	// FailCount 表示尝试加锁失败次数。
	FailCount uint64
	// UnlockCount 表示成功解锁次数。
	UnlockCount uint64
	// SpinCount 表示累计自旋次数。
	SpinCount uint64
}

// QueryParseOptions 控制查询字符串解析兼容性。
type QueryParseOptions struct {
	// Strict 表示是否启用严格解析，遇到无效转义时返回错误。
	Strict bool
	// Semicolon 表示是否允许分号作为查询参数分隔符。
	Semicolon bool
	// PlusAsSpace 表示是否把加号解析为空格。
	PlusAsSpace bool
}

// BuffPool 表示固定大小的字节缓冲池。
type BuffPool interface {
	// Get 从缓冲池中取出一个指定大小的字节切片引用。
	Get() *[]byte
	// Put 将字节切片引用归还给缓冲池。
	Put(buffRef *[]byte)
}

// BuffFactory 表示按大小分桶的字节缓冲池工厂。
//
// 内部按最小大小、最大大小和步长分桶，CreateBuff 取最匹配的桶，ReuseBuff 归还到对应桶。
type BuffFactory interface {
	// CreateBuff 按请求大小创建或复用一个合适的字节缓冲。
	CreateBuff(size uint32) *[]byte
	// ReuseBuff 将 CreateBuff 返回的字节缓冲归还到对应池中。
	ReuseBuff(buffRef *[]byte)
}

// Event 一次性事件信号。Broadcast 后所有当前/未来 Wait 的 goroutine 立即返回。
// 典型用法：用 evtStop 作为 goroutine 退出信号。
// Close 后所有 Wait 返回 ErrEvtClosed。
type Event interface {
	// Broadcast 广播事件通知，使当前和后续等待方立即返回。
	Broadcast() error
	// Wait 阻塞等待事件被广播或关闭。
	Wait() error
	// GetChan 返回事件内部通知通道。
	GetChan() chan byte
	// Close 关闭事件，唤醒所有等待方。
	Close()
	// IsClose 返回事件是否已经关闭。
	IsClose() bool
	// WaitUntilTimeout 等待事件广播，超时返回 ErrEvtWaitTimeout。
	WaitUntilTimeout(timeoutMSec uint32) error
	// GetBroadcastCount 返回累计广播次数。
	GetBroadcastCount() uint64
	// GetWaiterCount 返回当前等待方数量。
	GetWaiterCount() int64
	// GetTimeoutCount 返回累计等待超时次数。
	GetTimeoutCount() uint64
	// GetWaitSuccCount 返回累计成功等待次数。
	GetWaitSuccCount() uint64
	// GetWaitClosedCount 返回累计因事件关闭而结束等待的次数。
	GetWaitClosedCount() uint64
}

// FastLock 自旋互斥锁，适合临界区极短的场景；持锁时间长的场景仍应使用 sync.Mutex。
// TryLock 的 maxCnt 是自旋上限，0 表示无限自旋直到拿到锁。
type FastLock interface {
	// TryLock 尝试自旋加锁，maxCnt 为最大自旋次数，0 表示一直自旋直到成功。
	TryLock(maxCnt uint32) error
	// Unlock 解锁，未加锁时行为由具体实现决定。
	Unlock()
	// TryUnlock 尝试解锁并返回是否成功。
	TryUnlock() bool
	// IsLocked 返回当前是否处于加锁状态。
	IsLocked() bool
	// TryLockWithTimeout 在指定超时时间内尝试加锁。
	TryLockWithTimeout(timeout time.Duration) error
	// Stats 返回 FastLock 运行统计快照。
	Stats() FastLockStats
}

// Mail SMTP 邮件发送参数。
type Mail struct {
	// User 表示 SMTP 登录用户或发件人账号。
	User string
	// AliasName 表示发件人显示名称。
	AliasName string
	// Pwd 表示 SMTP 登录密码或授权码。
	Pwd string
	// Host 表示 SMTP 服务器主机。
	Host string
	// Port 表示 SMTP 服务器端口。
	Port int
	// To 表示收件人地址列表。
	To []string
	// Subject 表示邮件主题。
	Subject string
	// Body 表示邮件正文。
	Body string
}

// MailOptions 扩展 Mail 的可选发送参数。
type MailOptions struct {
	// FromName 表示覆盖 Mail.AliasName 的发件人显示名称。
	FromName string
	// Cc 表示抄送人地址列表。
	Cc []string
	// Bcc 表示密送人地址列表。
	Bcc []string
	// ReplyTo 表示回复地址列表。
	ReplyTo []string
	// TextBody 表示纯文本正文。
	TextBody string
	// HTMLBody 表示 HTML 正文。
	HTMLBody string
	// Headers 表示自定义邮件头。
	Headers map[string][]string
	// Attachments 表示附件列表。
	Attachments []MailAttachment
	// Timeout 表示发送邮件的超时时间。
	Timeout time.Duration
	// TLSConfig 表示 SMTP TLS 配置。
	TLSConfig *tls.Config
	// SSL 表示是否使用隐式 TLS 连接。
	SSL bool
	// LocalName 表示 SMTP 握手时使用的本地主机名。
	LocalName string
	// DisableAuth 表示是否禁用 SMTP 登录认证。
	DisableAuth bool
	// AllowEmptySubject 表示是否允许邮件主题为空。
	AllowEmptySubject bool
}

// MailAttachment 描述一份邮件附件。
type MailAttachment struct {
	// FilePath 表示附件文件路径。
	FilePath string
	// Name 表示附件在邮件中的展示名称。
	Name string
}

// MailSender 是底层 SMTP 发送器接口，便于测试替换。
type MailSender interface {
	// Send 发送一封已经构建好的邮件消息。
	Send(from string, to []string, msg []byte) error
}

// MailDialer 是 gomail 拨号发送接口，便于注入自定义发送实现。
type MailDialer interface {
	// DialAndSend 建立连接并发送一封或多封 gomail 消息。
	DialAndSend(...*gomail.Message) error
}

// Observer 观察者：响应 Dispatcher.Notify 推送的消息。
// 注意：在异步 Dispatcher 中 OnNotify 在专用 goroutine 中调用，业务自行保证线程安全。
type Observer interface {
	// OnNotify 在 Dispatcher 派发消息时调用。
	OnNotify(name string, params ...interface{})
}

// NotifyMsg 是异步派发队列中的通知消息。
type NotifyMsg struct {
	// Name 表示消息名称。
	Name string
	// Params 表示消息参数列表。
	Params []interface{}
}

// ObserverPanicHandler 用于处理观察者回调中的 panic。
type ObserverPanicHandler func(msgName string, observer Observer, recovered interface{})

// Dispatcher 单一消息名的发布订阅。
// 同步实现立刻在 Notify 调用方 goroutine 中分发；异步实现把消息入队、由专用 goroutine 顺序投递。
type Dispatcher interface {
	// AddObserver 添加一个观察者。
	AddObserver(o Observer)

	// RemoveObserver 移除一个观察者。
	RemoveObserver(o Observer)

	// Notify 向所有观察者派发消息参数。
	Notify(params ...interface{})

	// GetMsgName 返回 Dispatcher 负责的消息名。
	GetMsgName() string
	// SetPanicHandler 设置观察者回调 panic 时的处理器。
	SetPanicHandler(handler ObserverPanicHandler)
	// GetObserverCount 返回当前观察者数量。
	GetObserverCount() int
	// ClearObservers 移除全部观察者。
	ClearObservers()
	// GetNotifyCount 返回累计派发次数。
	GetNotifyCount() uint64
	// GetPanicCount 返回观察者回调累计 panic 次数。
	GetPanicCount() uint64
}

// AsyncDispatcher 是带后台队列的异步派发器。
type AsyncDispatcher interface {
	// Dispatcher 让异步派发器具备普通派发器能力。
	Dispatcher

	// TryNotify 尝试将消息入队，队列满时返回 false。
	TryNotify(params ...interface{}) bool
	// NotifyWithTimeout 在指定超时时间内尝试将消息入队。
	NotifyWithTimeout(timeout time.Duration, params ...interface{}) bool
	// Start 启动后台派发循环。
	Start()
	// Stop 停止后台派发循环。
	Stop()
	// StopWithTimeout 在指定超时时间内等待后台派发循环停止。
	StopWithTimeout(timeout time.Duration) bool
	// IsStarted 返回异步派发器是否已经启动。
	IsStarted() bool
	// IsStopped 返回异步派发器是否已经停止。
	IsStopped() bool
	// GetQueueLen 返回当前异步消息队列长度。
	GetQueueLen() int
	// GetQueueCap 返回异步消息队列容量。
	GetQueueCap() int
	// GetEnqueueCount 返回累计成功入队次数。
	GetEnqueueCount() uint64
	// GetDropCount 返回累计丢弃消息次数。
	GetDropCount() uint64
}

// NotifyCenter 多消息名的发布订阅中心，内部维护 msgName -> Dispatcher 的映射。
type NotifyCenter interface {
	// AddObserver 为指定消息名添加观察者。
	AddObserver(msgName string, o Observer)
	// RemoveObserver 从指定消息名移除观察者。
	RemoveObserver(msgName string, o Observer)
	// Notify 向指定消息名的观察者派发消息。
	Notify(msgName string, params ...interface{})
	// AddDispatcher 注册指定消息名的派发器。
	AddDispatcher(msgName string, dispatcher Dispatcher)
	// RemoveDispatcher 移除指定消息名的派发器。
	RemoveDispatcher(msgName string)
	// GetDispatcher 返回指定消息名的派发器。
	GetDispatcher(msgName string) (Dispatcher, bool)
	// GetDispatcherCount 返回当前派发器数量。
	GetDispatcherCount() int
	// ClearDispatchers 移除全部派发器。
	ClearDispatchers()
}

// Queue 通用队列抽象。
// 不同实现的并发与容量语义：
//   - LinkedQueue        非线程安全，无容量限制
//   - SyncLinkedQueue    线程安全，无容量限制
//   - SyncLimitQueue     线程安全，固定容量，满则 Enqueue 返回 ErrQueFullQueue
type Queue interface {
	// Enqueue 将元素加入队列。
	Enqueue(item interface{}) error

	// Dequeue 从队列取出一个元素。
	Dequeue() (interface{}, error)

	// GetSize 返回当前队列元素数量。
	GetSize() uint64
	// Peek 查看队首元素但不移除。
	Peek() (interface{}, error)
	// TryDequeue 尝试取出一个元素，队列为空时返回 false。
	TryDequeue() (interface{}, bool)
	// IsEmpty 返回队列是否为空。
	IsEmpty() bool
	// Clear 清空队列并返回清理的元素数量。
	Clear() uint64
	// Drain 最多取出 max 个元素；max 为 0 时取出全部元素。
	Drain(max uint64) []interface{}
	// ToSlice 以切片形式返回当前队列快照。
	ToSlice() []interface{}
}

// LimitQueue 是带容量限制的线程安全队列。
type LimitQueue interface {
	// Queue 让限长队列具备普通队列能力。
	Queue

	// TryEnqueue 尝试将元素加入队列，队列满时返回 false。
	TryEnqueue(item interface{}) bool
	// GetMaxSize 返回队列最大容量。
	GetMaxSize() uint64
	// SetMaxSize 设置队列最大容量。
	SetMaxSize(maxSize uint64)
	// IsFull 返回队列是否已满。
	IsFull() bool
	// Available 返回队列剩余可用容量。
	Available() uint64
}

// SetType Set 的底层键类型。不同 SetType 内部使用不同 map，避免 interface{} 比较开销。
type SetType = uint8

const (
	// SET_TYPE_INT 表示以 int 作为键类型的集合。
	SET_TYPE_INT SetType = iota
	// SET_TYPE_UINT 表示以 uint 作为键类型的集合。
	SET_TYPE_UINT
	// SET_TYPE_FLOAT 表示以 float64 作为键类型的集合。
	SET_TYPE_FLOAT
	// SET_TYPE_OBJ 表示以可比较对象作为键类型的集合。
	SET_TYPE_OBJ
)

// Set 集合接口。Add 重复元素不报错（幂等），Pop 顺序由实现决定（不保证 FIFO/LIFO）。
type Set interface {
	// Add 向集合添加元素。
	Add(data interface{}) error
	// Remove 从集合移除元素。
	Remove(data interface{}) error
	// Exist 判断元素是否存在，并返回类型校验错误。
	Exist(data interface{}) (bool, error)
	// GetElements 返回集合元素快照。
	GetElements() []interface{}
	// GetSize 返回集合元素数量。
	GetSize() int
	// Pop 弹出一个元素，顺序由实现决定。
	Pop() (interface{}, bool)
	// AddMany 批量添加元素。
	AddMany(values ...interface{}) error
	// RemoveMany 批量移除元素。
	RemoveMany(values ...interface{}) error
	// Contains 判断元素是否存在，错误场景返回 false。
	Contains(data interface{}) bool
	// IsEmpty 返回集合是否为空。
	IsEmpty() bool
	// Clear 清空集合并返回被清理元素数量。
	Clear() int
}

// SimpleBuffer 单缓冲区读写双指针抽象，用于网络拆包等场景。
// 典型流程：读 socket 到 GetWriteBuff -> UpdateWriteOffset -> 从 GetData 取数据消费 ->
// 数据消费完调用 MoveDataToBegin 整理空间。
type SimpleBuffer interface {
	// GetWriteBuffSize 返回当前可写缓冲区大小。
	GetWriteBuffSize() int
	// GetWriteBuff 返回可直接写入的尾部缓冲区。
	GetWriteBuff() []byte
	// GetDataLen 返回当前可读数据长度。
	GetDataLen() int
	// GetData 返回当前可读数据切片。
	GetData() []byte
	// MoveDataToBegin 把未消费数据搬到缓冲区起始位置，腾出尾部写空间。
	MoveDataToBegin()
	// UpdateWriteOffset 在外部写入 GetWriteBuff 后推进写偏移。
	UpdateWriteOffset(addLen uint32)
	// Skip 跳过指定长度的可读数据。
	Skip(skipLen uint32)
	// SimulateRead 拷贝数据但不前移读指针，用于预读。
	SimulateRead(b []byte) (n int, err error)
	// Read 读取数据并推进读指针。
	Read(b []byte) (n int, err error)
	// Write 写入数据。
	Write(data []byte) (int, error)
	// WriteAll 写入完整数据，空间不足时返回错误。
	WriteAll(data []byte) error
	// WriteAutoGrow 写入数据，空间不足时自动扩容。
	WriteAutoGrow(data []byte) (int, error)
	// ReadByte 读取一个字节。
	ReadByte() (byte, error)
	// Peek 查看前 n 字节但不推进读指针。
	Peek(n int) ([]byte, error)
	// ReadAll 读取全部可读数据并清空读区。
	ReadAll() []byte
	// Clear 清空缓冲区数据。
	Clear()
	// Capacity 返回缓冲区容量。
	Capacity() int
	// IsEmpty 返回缓冲区是否没有可读数据。
	IsEmpty() bool
	// IsFull 返回缓冲区是否没有可写空间。
	IsFull() bool
	// Reset 重置缓冲区容量并清空数据。
	Reset(capacity int) error
	// Grow 增加缓冲区容量。
	Grow(additional int) error
	// ReadFrom 从 io.Reader 读取数据并写入缓冲区。
	ReadFrom(r io.Reader) (int64, error)
}

// SnowflakeStats 是 Snowflake 生成器的运行快照。
type SnowflakeStats struct {
	// WorkerId 表示当前生成器使用的工作节点 ID。
	WorkerId int64
	// Initialized 表示生成器是否已经初始化。
	Initialized bool
	// LastTimestamp 表示最近一次生成 ID 使用的时间戳。
	LastTimestamp int64
	// Sequence 表示当前毫秒内的序列号。
	Sequence int64
	// GeneratedCount 表示累计生成 ID 数量。
	GeneratedCount uint64
	// WaitCount 表示因同毫秒序列耗尽而等待下一毫秒的次数。
	WaitCount uint64
	// ClockBackCount 表示检测到系统时钟回拨的次数。
	ClockBackCount uint64
}

// Snowflake 分布式唯一 ID 生成器（标准 Twitter Snowflake 算法）。
// 64 位 ID = 时间戳(41) + workerId(10) + sequence(12)，每毫秒每节点最多 4096 个 ID。
// 初始化时必须为每个进程分配唯一 workerId，否则会产生重复 ID。
type Snowflake interface {
	// Init 使用指定 workerId 初始化生成器。
	Init(workerId int64) error
	// NextId 生成一个 ID，失败时返回 false。
	NextId() (int64, bool)
	// NextIdErr 生成一个 ID，失败时返回具体错误。
	NextIdErr() (int64, error)
	// NextIds 批量生成 count 个 ID，失败时返回 false。
	NextIds(count int) ([]int64, bool)
	// NextIdsErr 批量生成 count 个 ID，失败时返回具体错误。
	NextIdsErr(count int) ([]int64, error)
	// UnmarshalId 解出毫秒时间戳、workerId 和序列号，便于排查 ID 归属。
	UnmarshalId(id int64) (int64, int64, int64)
	// SetMaxClockBackwardMs 设置允许等待的最大时钟回拨毫秒数。
	SetMaxClockBackwardMs(maxBackwardMs int64)
	// Stats 返回 Snowflake 生成器运行统计快照。
	Stats() SnowflakeStats
}

// UtilMgr 工具方法的上帝接口：把全部无状态工具函数收拢在一起，避免业务包散落引用。
type UtilMgr interface {
	// AesEncrypt 使用默认 AES 模式加密原始数据。
	AesEncrypt(origData []byte, key []byte) ([]byte, error)
	// AesDecrypt 使用默认 AES 模式解密密文数据。
	AesDecrypt(encrypted []byte, key []byte) ([]byte, error)
	// AesEncryptBase64 使用 Base64 编码密钥加密数据，并返回 Base64 编码密文。
	AesEncryptBase64(origData []byte, keyBase64 string) (string, error)
	// AesDecryptBase64 使用 Base64 编码密钥解密 Base64 编码密文。
	AesDecryptBase64(encryptedBase64 string, keyBase64 string) ([]byte, error)
	// AesEncryptCBC 使用 AES-CBC 和 PKCS7 填充加密数据。
	AesEncryptCBC(origData []byte, key []byte) ([]byte, error)
	// AesDecryptCBC 使用 AES-CBC 和 PKCS7 填充解密数据。
	AesDecryptCBC(encrypted []byte, key []byte) ([]byte, error)
	// AesEncryptCBCBase64 使用 AES-CBC 加密数据，并返回 Base64 编码密文。
	AesEncryptCBCBase64(origData []byte, keyBase64 string) (string, error)
	// AesDecryptCBCBase64 使用 AES-CBC 解密 Base64 编码密文。
	AesDecryptCBCBase64(encryptedBase64 string, keyBase64 string) ([]byte, error)
	// AesEncryptGCM 使用 AES-GCM 加密数据。
	AesEncryptGCM(origData []byte, key []byte) ([]byte, error)
	// AesDecryptGCM 使用 AES-GCM 解密数据。
	AesDecryptGCM(encrypted []byte, key []byte) ([]byte, error)
	// AesEncryptGCMBase64 使用 AES-GCM 加密数据，并返回 Base64 编码密文。
	AesEncryptGCMBase64(origData []byte, keyBase64 string) (string, error)
	// AesDecryptGCMBase64 使用 AES-GCM 解密 Base64 编码密文。
	AesDecryptGCMBase64(encryptedBase64 string, keyBase64 string) ([]byte, error)

	// Crc16CheckSum 使用默认 CRC16-CCITT 算法计算校验值。
	Crc16CheckSum(data []byte) (uint16, error)
	// Crc16ModbusCheckSum 使用 CRC-16/MODBUS 参数计算校验值。
	Crc16ModbusCheckSum(data []byte) (uint16, error)
	// Crc16ARCCheckSum 使用 CRC-16/ARC 参数计算校验值。
	Crc16ARCCheckSum(data []byte) (uint16, error)
	// Crc16CCITTFalseCheckSum 使用 CRC-16/CCITT-FALSE 参数计算校验值。
	Crc16CCITTFalseCheckSum(data []byte) (uint16, error)
	// Crc16XMODEMCheckSum 使用 CRC-16/XMODEM 参数计算校验值。
	Crc16XMODEMCheckSum(data []byte) (uint16, error)
	// Crc16X25CheckSum 使用 CRC-16/X-25 参数计算校验值。
	Crc16X25CheckSum(data []byte) (uint16, error)
	// Crc16CheckSumWithParams 使用指定 CRC16 参数计算校验值。
	Crc16CheckSumWithParams(data []byte, params Crc16Params) (uint16, error)
	// Checksum 使用指定 CRC16 参数计算校验值，输入为空时不返回错误。
	Checksum(data []byte, params Crc16Params) uint16
	// NewCrc16Hash 创建 hash.Hash 风格的 CRC16 计算器。
	NewCrc16Hash(params Crc16Params) Hash16
	// MakeCrc16Table 创建指定 CRC16 参数的查表缓存。
	MakeCrc16Table(params Crc16Params) *Crc16Table

	// EcdsaGenerateKey 按指定曲线生成 ECDSA 公私钥。
	EcdsaGenerateKey(ellType ElliptcType) (pubKey []byte, priKey []byte, err error)
	// EcdsaGenerateKeyDer 按指定曲线生成 DER 编码的 ECDSA 公私钥。
	EcdsaGenerateKeyDer(ellType ElliptcType) (pubKey []byte, priKey []byte, err error)
	// EcdsaGenerateKeyPem 按指定曲线生成 PEM 文本格式的 ECDSA 公私钥。
	EcdsaGenerateKeyPem(ellType ElliptcType) (pemPubKey string, pemPriKey string, err error)
	// EcdsaGenerateKeyPKCS8 按指定曲线生成 PKCS8 私钥格式的 ECDSA 公私钥。
	EcdsaGenerateKeyPKCS8(ellType ElliptcType) (pubKey []byte, priKey []byte, err error)
	// EcdsaSign 使用默认哈希算法和私钥对数据签名，返回 r 和 s。
	EcdsaSign(origData []byte, priKey []byte) ([]byte, []byte, error)
	// EcdsaVerify 使用默认哈希算法和公钥校验 r、s 签名。
	EcdsaVerify(origData []byte, rText []byte, sText []byte, pubKey []byte) (bool, error)
	// EcdsaSignSHA256 使用 SHA-256 和私钥对数据签名，返回 r 和 s。
	EcdsaSignSHA256(origData []byte, priKey []byte) ([]byte, []byte, error)
	// EcdsaSignSHA384 使用 SHA-384 和私钥对数据签名，返回 r 和 s。
	EcdsaSignSHA384(origData []byte, priKey []byte) ([]byte, []byte, error)
	// EcdsaSignSHA512 使用 SHA-512 和私钥对数据签名，返回 r 和 s。
	EcdsaSignSHA512(origData []byte, priKey []byte) ([]byte, []byte, error)
	// EcdsaVerifySHA256 使用 SHA-256 和公钥校验 r、s 签名。
	EcdsaVerifySHA256(origData []byte, rText []byte, sText []byte, pubKey []byte) (bool, error)
	// EcdsaVerifySHA384 使用 SHA-384 和公钥校验 r、s 签名。
	EcdsaVerifySHA384(origData []byte, rText []byte, sText []byte, pubKey []byte) (bool, error)
	// EcdsaVerifySHA512 使用 SHA-512 和公钥校验 r、s 签名。
	EcdsaVerifySHA512(origData []byte, rText []byte, sText []byte, pubKey []byte) (bool, error)
	// EcdsaSignPem 使用 PEM 私钥对数据签名，返回 r 和 s。
	EcdsaSignPem(origData []byte, pemPriKey string) ([]byte, []byte, error)
	// EcdsaVerifyPem 使用 PEM 公钥校验 r、s 签名。
	EcdsaVerifyPem(origData []byte, rText []byte, sText []byte, pemPubKey string) (bool, error)
	// EcdsaSignWithHash 使用指定哈希算法和私钥对数据签名，返回 r 和 s。
	EcdsaSignWithHash(origData []byte, priKey []byte, hashType EcdsaHash) ([]byte, []byte, error)
	// EcdsaVerifyWithHash 使用指定哈希算法和公钥校验 r、s 签名。
	EcdsaVerifyWithHash(origData []byte, rText []byte, sText []byte, pubKey []byte, hashType EcdsaHash) (bool, error)
	// EcdsaSignASN1 使用默认哈希算法和私钥生成 ASN.1 编码签名。
	EcdsaSignASN1(origData []byte, priKey []byte) ([]byte, error)
	// EcdsaVerifyASN1 使用默认哈希算法和公钥校验 ASN.1 编码签名。
	EcdsaVerifyASN1(origData []byte, signData []byte, pubKey []byte) (bool, error)
	// EcdsaSignASN1WithHash 使用指定哈希算法和私钥生成 ASN.1 编码签名。
	EcdsaSignASN1WithHash(origData []byte, priKey []byte, hashType EcdsaHash) ([]byte, error)
	// EcdsaVerifyASN1WithHash 使用指定哈希算法和公钥校验 ASN.1 编码签名。
	EcdsaVerifyASN1WithHash(origData []byte, signData []byte, pubKey []byte, hashType EcdsaHash) (bool, error)

	// SendMail 使用默认选项同步发送邮件。
	SendMail(m *Mail) error
	// SendMailWithOptions 使用扩展选项同步发送邮件。
	SendMailWithOptions(m *Mail, opts *MailOptions) error
	// SendMailWithDialer 使用自定义拨号器同步发送邮件。
	SendMailWithDialer(m *Mail, opts *MailOptions, dialer MailDialer) error
	// SendMailWithSender 使用自定义底层发送器发送已构建邮件。
	SendMailWithSender(m *Mail, opts *MailOptions, sender MailSender) error
	// BuildMailMessage 根据邮件参数和选项构建 gomail 消息。
	BuildMailMessage(m *Mail, opts *MailOptions) (*gomail.Message, error)
	// NewMailDialer 根据邮件参数和选项创建邮件拨号器。
	NewMailDialer(m *Mail, opts *MailOptions) (MailDialer, error)
	// ValidateMail 校验邮件参数和选项。
	ValidateMail(m *Mail, opts *MailOptions) error

	// RsaGenerateKey 生成指定位数的 PKCS1 RSA 公私钥。
	RsaGenerateKey(bits int) (pubKey []byte, priKey []byte, err error)
	// RsaGenerateKeyPem 生成指定位数的 PEM 文本格式 PKCS1 RSA 公私钥。
	RsaGenerateKeyPem(bits int) (pemPubKey string, pemPriKey string, err error)
	// RsaGenerateKeyPKCS8 生成指定位数的 PKCS8 RSA 公私钥。
	RsaGenerateKeyPKCS8(bits int) (pubKey []byte, priKey []byte, err error)
	// RsaGenerateKeyPemPKCS8 生成指定位数的 PEM 文本格式 PKCS8 RSA 公私钥。
	RsaGenerateKeyPemPKCS8(bits int) (pemPubKey string, pemPriKey string, err error)
	// RsaEncrypt 使用 RSA 公钥加密数据。
	RsaEncrypt(origData []byte, pubKey []byte) ([]byte, error)
	// RsaDecrypt 使用 RSA 私钥解密数据。
	RsaDecrypt(encrypted []byte, priKey []byte) ([]byte, error)
	// RsaEncryptPem 使用 PEM 公钥加密数据。
	RsaEncryptPem(origData []byte, pemPubKey string) ([]byte, error)
	// RsaDecryptPem 使用 PEM 私钥解密数据。
	RsaDecryptPem(encrypted []byte, pemPriKey string) ([]byte, error)
	// RsaSign 使用默认哈希算法和 RSA 私钥签名数据。
	RsaSign(origData []byte, priKey []byte) ([]byte, error)
	// RsaVerify 使用默认哈希算法和 RSA 公钥校验签名。
	RsaVerify(origData []byte, signData []byte, pubKey []byte) (bool, error)
	// RsaSignPem 使用 PEM 私钥签名数据。
	RsaSignPem(origData []byte, pemPriKey string) ([]byte, error)
	// RsaVerifyPem 使用 PEM 公钥校验签名。
	RsaVerifyPem(origData []byte, signData []byte, pemPubKey string) (bool, error)
	// RsaEncryptOAEP 使用 RSA-OAEP 和默认哈希算法加密数据。
	RsaEncryptOAEP(origData []byte, pubKey []byte) ([]byte, error)
	// RsaDecryptOAEP 使用 RSA-OAEP 和默认哈希算法解密数据。
	RsaDecryptOAEP(encrypted []byte, priKey []byte) ([]byte, error)
	// RsaEncryptOAEPPem 使用 PEM 公钥进行 RSA-OAEP 加密。
	RsaEncryptOAEPPem(origData []byte, pemPubKey string) ([]byte, error)
	// RsaDecryptOAEPPem 使用 PEM 私钥进行 RSA-OAEP 解密。
	RsaDecryptOAEPPem(encrypted []byte, pemPriKey string) ([]byte, error)
	// RsaEncryptOAEPWithHash 使用指定哈希算法和标签进行 RSA-OAEP 加密。
	RsaEncryptOAEPWithHash(origData []byte, pubKey []byte, hashType RsaHash, label []byte) ([]byte, error)
	// RsaDecryptOAEPWithHash 使用指定哈希算法和标签进行 RSA-OAEP 解密。
	RsaDecryptOAEPWithHash(encrypted []byte, priKey []byte, hashType RsaHash, label []byte) ([]byte, error)
	// RsaSignWithHash 使用指定哈希算法和 RSA 私钥签名数据。
	RsaSignWithHash(origData []byte, priKey []byte, hashType RsaHash) ([]byte, error)
	// RsaVerifyWithHash 使用指定哈希算法和 RSA 公钥校验签名。
	RsaVerifyWithHash(origData []byte, signData []byte, pubKey []byte, hashType RsaHash) (bool, error)
	// RsaSignPSS 使用 RSA-PSS 和默认哈希算法签名数据。
	RsaSignPSS(origData []byte, priKey []byte) ([]byte, error)
	// RsaVerifyPSS 使用 RSA-PSS 和默认哈希算法校验签名。
	RsaVerifyPSS(origData []byte, signData []byte, pubKey []byte) (bool, error)
	// RsaSignPSSWithHash 使用指定哈希算法进行 RSA-PSS 签名。
	RsaSignPSSWithHash(origData []byte, priKey []byte, hashType RsaHash) ([]byte, error)
	// RsaVerifyPSSWithHash 使用指定哈希算法校验 RSA-PSS 签名。
	RsaVerifyPSSWithHash(origData []byte, signData []byte, pubKey []byte, hashType RsaHash) (bool, error)
	// ParseRsaPublicKey 解析 DER 编码的 RSA 公钥。
	ParseRsaPublicKey(pubKey []byte) (*rsa.PublicKey, error)
	// ParseRsaPrivateKey 解析 DER 编码的 RSA 私钥。
	ParseRsaPrivateKey(priKey []byte) (*rsa.PrivateKey, error)

	// ParseQuery 解析查询字符串。
	ParseQuery(query string) (url.Values, error)
	// ParseQueryStrict 以严格模式解析查询字符串。
	ParseQueryStrict(query string) (url.Values, error)
	// ParseQueryForm 按表单兼容规则解析查询字符串。
	ParseQueryForm(query string) (url.Values, error)
	// ParseQueryWithOptions 按指定选项解析查询字符串。
	ParseQueryWithOptions(query string, opts QueryParseOptions) (url.Values, error)
	// EncodeURIComponent 按前端 encodeURIComponent 兼容规则编码字符串。
	EncodeURIComponent(s string) string
	// DecodeURIComponent 按前端 decodeURIComponent 兼容规则解码字符串。
	DecodeURIComponent(s string) (string, error)
	// EncodePathSegment 编码 URL 路径片段。
	EncodePathSegment(s string) string
	// DecodePathSegment 解码 URL 路径片段。
	DecodePathSegment(s string) (string, error)
	// EncodeQuery 将查询参数编码为查询字符串。
	EncodeQuery(values url.Values) string
	// EncodeQuerySorted 按键排序后将查询参数编码为查询字符串。
	EncodeQuerySorted(values url.Values) string
	// AddQuery 将查询参数追加到原始 URL。
	AddQuery(rawURL string, values url.Values) (string, error)
	// WriteSvrStatRsp 向 http.ResponseWriter 输出标准状态响应。
	WriteSvrStatRsp(writer http.ResponseWriter, code int32, result string) error
	// WriteJSONResponse 向 http.ResponseWriter 输出 JSON 响应并设置状态码。
	WriteJSONResponse(writer http.ResponseWriter, statusCode int, v interface{}) error

	// FormatTimeStr 按指定格式把时间写入 builder 并返回字符串。
	FormatTimeStr(format string, t time.Time, builder *strings.Builder) (string, error)
	// FormatInt 将整数写入 builder，maxLength 和 bFillZero 控制定宽补零。
	FormatInt(num int64, maxLength uint32, bFillZero bool, builder *strings.Builder) string
	// FormatUint 将无符号整数写入 builder，maxLength 和 bFillZero 控制定宽补零。
	FormatUint(num uint64, maxLength uint32, bFillZero bool, builder *strings.Builder) string
	// PowerOfTen 返回 10 的 power 次方。
	PowerOfTen(power uint32) uint64

	// ProtectRun 执行函数并自动恢复 panic。
	ProtectRun(entry func())
	// ProtectRunWithHandler 执行函数并在 panic 时调用自定义处理器。
	ProtectRunWithHandler(entry func(), handler func(interface{}, []byte))
	// RunDangerCode 执行可能 panic 的函数，调试模式下不恢复 panic。
	RunDangerCode(entry func(), bDebugMode bool)

	// GetFileSize 返回指定文件大小。
	GetFileSize(path string) (int64, error)
	// IsFileExist 判断指定路径是否存在文件或目录。
	IsFileExist(path string) (bool, error)
	// LoadJsonConf 读取 JSON 配置文件并反序列化到 v，可在解析前通过 decodeCb 解码数据。
	LoadJsonConf(v interface{}, path string, decodeCb func(data []byte) ([]byte, error)) error
	// LoadJsonConfStrict 读取 JSON 配置文件并使用严格模式反序列化到 v。
	LoadJsonConfStrict(v interface{}, path string, decodeCb func(data []byte) ([]byte, error)) error

	// GetClassReflectName 返回对象的反射类型全名。
	GetClassReflectName(obj interface{}) (string, error)
	// GetReflectNameByType 返回 reflect.Type 对应的类型全名。
	GetReflectNameByType(t reflect.Type) string
	// GetFullPackageName 从类型全名中提取完整包名。
	GetFullPackageName(classReflectName string) string
	// GetFilePackageName 从完整包名中提取文件级包名。
	GetFilePackageName(fullPackName string) string
	// GetClassName 从类型全名中提取类名。
	GetClassName(classReflectName string) string
	// GetFilePackageClassName 从类型全名中提取文件级包名和类名。
	GetFilePackageClassName(classReflectName string) string

	// Daemon 以子进程方式拉起程序，进程退出后按 restartDelay 秒重启，chanExit 关闭后停止守护。
	Daemon(program string, args []string, restartDelay uint16, chanExit chan bool, onStartSucc func(cmd *exec.Cmd)) error

	// Random 返回 [min, max] 闭区间内的随机整数。
	Random(min int, max int) int
	// RandomIntn 返回 [0, n) 区间内的随机整数。
	RandomIntn(n int) int
	// SetRandomSeed 设置随机数种子。
	SetRandomSeed(seed int64)

	// KvGetInt 从字符串映射中读取 int，缺失或转换失败时返回默认值。
	KvGetInt(kv map[string]string, key string, defaultVal int) int
	// KvGetUint 从字符串映射中读取 uint，缺失或转换失败时返回默认值。
	KvGetUint(kv map[string]string, key string, defaultVal uint) uint
	// KvGetUint32 从字符串映射中读取 uint32，缺失或转换失败时返回默认值。
	KvGetUint32(kv map[string]string, key string, defaultVal uint32) uint32
	// KvGetInt32 从字符串映射中读取 int32，缺失或转换失败时返回默认值。
	KvGetInt32(kv map[string]string, key string, defaultVal int32) int32
	// KvGetUint64 从字符串映射中读取 uint64，缺失或转换失败时返回默认值。
	KvGetUint64(kv map[string]string, key string, defaultVal uint64) uint64
	// KvGetInt64 从字符串映射中读取 int64，缺失或转换失败时返回默认值。
	KvGetInt64(kv map[string]string, key string, defaultVal int64) int64
	// KvGetBool 从字符串映射中读取 bool，缺失或转换失败时返回默认值。
	KvGetBool(kv map[string]string, key string, defaultVal bool) bool
	// KvGetString 从字符串映射中读取 string，缺失时返回默认值。
	KvGetString(kv map[string]string, key string, defaultVal string) string

	// NewBuffPool 创建固定大小的字节缓冲池。
	NewBuffPool(buffSize uint32) BuffPool
	// NewBuffFactory 创建按大小分桶的字节缓冲池工厂。
	NewBuffFactory(minSize uint32, maxSize uint32, queueStep uint32) BuffFactory
	// NewEvent 创建一次性事件信号对象。
	NewEvent() Event
	// NewFastLock 创建自旋互斥锁。
	NewFastLock() FastLock
	// NewBaseDispatcher 创建基础同步派发器。
	NewBaseDispatcher(msgName string) Dispatcher
	// NewSyncDispatcher 创建同步派发器。
	NewSyncDispatcher(msgName string) Dispatcher
	// NewAsyncDispatcher 创建异步派发器。
	NewAsyncDispatcher(msgName string, maxMsgBuffSize uint16) AsyncDispatcher
	// NewNotifyCenter 创建通知中心。
	NewNotifyCenter() NotifyCenter
	// NewLinkedQueue 创建非线程安全链表队列。
	NewLinkedQueue() Queue
	// NewSyncLinkedQueue 创建线程安全链表队列。
	NewSyncLinkedQueue() Queue
	// NewSyncLimitQueue 创建线程安全限长队列。
	NewSyncLimitQueue(maxSize uint64) LimitQueue
	// NewSet 按指定底层键类型创建集合。
	NewSet(t SetType) Set
	// NewIntSet 创建整数集合。
	NewIntSet() Set
	// NewUintSet 创建无符号整数集合。
	NewUintSet() Set
	// NewFloatSet 创建浮点数集合。
	NewFloatSet() Set
	// NewObjectSet 创建可比较对象集合。
	NewObjectSet() Set
	// NewSimpleBuffer 创建指定容量的简单缓冲区。
	NewSimpleBuffer(capacity int) SimpleBuffer
	// NewSnowflake 创建未初始化的 Snowflake ID 生成器。
	NewSnowflake() Snowflake
	// NewSnowflakeWithWorker 创建并初始化指定 workerId 的 Snowflake ID 生成器。
	NewSnowflakeWithWorker(workerId int64) (Snowflake, error)
	// ParseSnowflakeId 解析 Snowflake ID 中的时间戳、workerId 和序列号。
	ParseSnowflakeId(id int64) (int64, int64, int64)
	// MaxSnowflakeWorkerId 返回 Snowflake 允许的最大 workerId。
	MaxSnowflakeWorkerId() int64
	// MaxSnowflakeSequence 返回 Snowflake 每毫秒允许的最大序列号。
	MaxSnowflakeSequence() int64
	// SnowflakeBeginTimestamp 返回 Snowflake 算法使用的起始毫秒时间戳。
	SnowflakeBeginTimestamp() int64
}
