package db

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	// DefaultRedis 是未指定名称时使用的默认 Redis 客户端名。
	DefaultRedis = "default"
)

// RedisSubscribeCb 订阅回调。
// channel 表示消息所在频道；pattern 表示匹配的模式；第三个参数表示消息体。
type RedisSubscribeCb func(channel string, pattern string, payload string)

// RedisClient 单个 Redis 实例的封装。
type RedisClient interface {
	// GetRedis 返回底层 redis.Client。
	GetRedis() *redis.Client
	// Ping 检查 Redis 连接是否可用。
	Ping() error
}

// RedisMgr 多 Redis 连接管理。命名风格与 MysqlMgr/MongoMgr 保持一致。
type RedisMgr interface {
	// OpenDefault 打开默认 Redis 客户端。
	OpenDefault(addr string, pwd string, db int) error
	// Open 以指定名称打开 Redis 客户端。
	Open(name string, addr string, pwd string, db int) error
	// CloseDefault 关闭默认 Redis 客户端。
	CloseDefault()
	// Close 关闭指定名称的 Redis 客户端。
	Close(name string)
	// CloseAll 关闭全部 Redis 客户端。
	CloseAll()
	// Ping 检查指定名称的 Redis 客户端是否可用。
	Ping(name string) error
	// GetDefaultClient 返回默认 Redis 客户端。
	GetDefaultClient() (RedisClient, bool)
	// GetClient 返回指定名称的 Redis 客户端。
	GetClient(name string) (RedisClient, bool)
}

// RedisHandle 无状态的 Redis 命令集合。
// 与 MysqlHandle 一致，底层 *redis.Client 作为第一个参数显式传入以便跨实例复用。
type RedisHandle interface {
	// Set 设置字符串或可序列化值。
	Set(cli *redis.Client, key string, value interface{}) error
	// SetContext 设置字符串或可序列化值。
	SetContext(ctx context.Context, cli *redis.Client, key string, value interface{}) error
	// SetWithExpire 设置值并同时设置过期时间。
	SetWithExpire(cli *redis.Client, key string, value interface{}, expiration time.Duration) error
	// SetWithExpireContext 设置值并同时设置过期时间。
	SetWithExpireContext(ctx context.Context, cli *redis.Client, key string, value interface{}, expiration time.Duration) error
	// SetNX 在 key 不存在时设置值。
	SetNX(cli *redis.Client, key string, value interface{}) (bool, error)
	// SetNXContext 在 key 不存在时设置值。
	SetNXContext(ctx context.Context, cli *redis.Client, key string, value interface{}) (bool, error)
	// SetNXWithExpire 在 key 不存在时设置值并同时设置过期时间。
	SetNXWithExpire(cli *redis.Client, key string, value interface{}, expiration time.Duration) (bool, error)
	// SetNXWithExpireContext 在 key 不存在时设置值并同时设置过期时间。
	SetNXWithExpireContext(ctx context.Context, cli *redis.Client, key string, value interface{}, expiration time.Duration) (bool, error)
	// MSet 批量设置多个 key-value。
	MSet(cli *redis.Client, pairs ...interface{}) error
	// Incr 将 key 对应整数值加一并返回新值。
	Incr(cli *redis.Client, key string) (int64, error)
	// Get 读取字符串值。
	Get(cli *redis.Client, key string) (string, error)
	// GetContext 读取字符串值。
	GetContext(ctx context.Context, cli *redis.Client, key string) (string, error)
	// MGet 批量读取多个 key 的值。
	MGet(cli *redis.Client, keys ...string) ([]interface{}, error)
	// Del 删除一个或多个 key，并返回删除数量。
	Del(cli *redis.Client, keys ...string) (int64, error)
	// DelContext 删除一个或多个 key，并返回删除数量。
	DelContext(ctx context.Context, cli *redis.Client, keys ...string) (int64, error)
	// Keys 按 pattern 查询 key 列表。
	Keys(cli *redis.Client, pattern string) ([]string, error)
	// Exists 返回存在的 key 数量。
	Exists(cli *redis.Client, keys ...string) (int64, error)

	// Expire 设置 key 的相对过期时间。
	Expire(cli *redis.Client, key string, expiration time.Duration) (bool, error)
	// ExpireContext 设置 key 的相对过期时间。
	ExpireContext(ctx context.Context, cli *redis.Client, key string, expiration time.Duration) (bool, error)
	// ExpireAt 设置 key 在指定时间过期。
	ExpireAt(cli *redis.Client, key string, tm time.Time) (bool, error)
	// Persist 移除 key 的过期时间，使其持久保存。
	Persist(cli *redis.Client, key string) (bool, error)
	// TTL 返回 key 的剩余生存时间。
	TTL(cli *redis.Client, key string) (time.Duration, error)

	// HSet 设置 Hash 字段值，并返回新增字段数量。
	HSet(cli *redis.Client, key string, field string, value interface{}) (int64, error)
	// HSetContext 设置 Hash 字段值，并返回新增字段数量。
	HSetContext(ctx context.Context, cli *redis.Client, key string, field string, value interface{}) (int64, error)
	// HGet 读取 Hash 指定字段的字符串值。
	HGet(cli *redis.Client, key string, field string) (string, error)
	// HGetContext 读取 Hash 指定字段的字符串值。
	HGetContext(ctx context.Context, cli *redis.Client, key string, field string) (string, error)
	// HMSet 批量设置 Hash 字段。
	HMSet(cli *redis.Client, key string, fields map[string]interface{}) error
	// HMGet 批量读取 Hash 字段。
	HMGet(cli *redis.Client, key string, fields ...string) ([]interface{}, error)
	// HGetAll 读取 Hash 的全部字段和值。
	HGetAll(cli *redis.Client, key string) (map[string]string, error)
	// HExists 判断 Hash 字段是否存在。
	HExists(cli *redis.Client, key string, field string) (bool, error)
	// HSetNX 在 Hash 字段不存在时设置字段值。
	HSetNX(cli *redis.Client, key string, field string, value interface{}) (bool, error)
	// HDel 删除 Hash 字段并返回删除数量。
	HDel(cli *redis.Client, key string, field string) (int64, error)

	// RPush 将元素追加到 List 右侧。
	RPush(cli *redis.Client, key string, values ...interface{}) error
	// LSet 按下标设置 List 元素。
	LSet(cli *redis.Client, key string, index int64, value interface{}) error
	// BLPop 阻塞式左弹出；timeout 为 0 表示无限等待。
	BLPop(cli *redis.Client, timeout time.Duration, keys ...string) ([]string, error)
	// LRem 从 List 中删除指定数量的匹配元素。
	LRem(cli *redis.Client, key string, count int64, value interface{}) (int64, error)
	// LLen 返回 List 长度。
	LLen(cli *redis.Client, key string) (int64, error)
	// LRange 返回 List 指定范围内的元素。
	LRange(cli *redis.Client, key string, start int64, stop int64) ([]string, error)

	// SubscribeContext 阻塞订阅指定频道；ctx 取消时退出，退出后关闭 exitChan 通知调用方。
	SubscribeContext(ctx context.Context, cli *redis.Client, exitChan chan byte, cb RedisSubscribeCb, channels ...string)
	// Publish 向指定频道发布消息。
	Publish(cli *redis.Client, channel string, message interface{}) error
}
