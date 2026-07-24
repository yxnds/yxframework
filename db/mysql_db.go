package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

var (
	// ErrMapperIsNil 表示 NamedExec、NamedInsert 等命名参数接口未传入结构体或 map。
	ErrMapperIsNil = errors.New("mapper is nil")
	// ErrAffectRowZero 表示期望影响行数大于 0，但实际影响行数为 0，常见于乐观锁更新失败。
	ErrAffectRowZero = errors.New("affect row is zero")
	// ErrDbInsertFieldNotFind 表示表结构中没有找到可用于插入的字段。
	ErrDbInsertFieldNotFind = errors.New("insert field not find")
	// ErrDbReplaceFieldNotFind 表示表结构中没有找到可用于插入或更新替换的字段。
	ErrDbReplaceFieldNotFind = errors.New("insert field and update field not find")
	// ErrDbKeyFieldNotDefine 表示表结构中没有定义主键字段。
	ErrDbKeyFieldNotDefine = errors.New("key field not define")
	// ErrDbUpdateFieldNotDefine 表示表结构中没有定义可更新字段。
	ErrDbUpdateFieldNotDefine = errors.New("update field not define")
	// ErrDbTableNameEmpty 表示表名为空。
	ErrDbTableNameEmpty = errors.New("table name is empty")
	// ErrDbRowObjCreateFuncIsNil 表示行对象创建函数为空。
	ErrDbRowObjCreateFuncIsNil = errors.New("row object create func is nil")
	// ErrDbRowObjectInvalid 表示行对象必须是结构体或结构体指针。
	ErrDbRowObjectInvalid = errors.New("row object must be struct pointer or struct")
	// ErrDbClientNil 表示数据库客户端为空。
	ErrDbClientNil = errors.New("db client is nil")
	// ErrDbTxNil 表示数据库事务为空。
	ErrDbTxNil = errors.New("db transaction is nil")
)

const (
	// DefaultMysql 表示未指定名称时使用的默认 MySQL 客户端名。
	DefaultMysql = "default"

	// INSERT_TAG 表示结构体字段参与 INSERT 字段生成的标签值。
	INSERT_TAG = "i"
	// SELECT_TAG 表示结构体字段参与 SELECT 字段生成的标签值。
	SELECT_TAG = "s"
	// SELECT_KEY_TAG 表示结构体字段参与 SELECT 条件键生成的标签值。
	SELECT_KEY_TAG = "sk"
	// UPDATE_TAG 表示结构体字段参与 UPDATE 字段生成的标签值。
	UPDATE_TAG = "u"
	// UPDATE_KEY_TAG 表示结构体字段参与 UPDATE 条件键生成的标签值。
	UPDATE_KEY_TAG = "uk"
)

// SqlQueryFieldData 单字段值的原始字节表示，便于业务自定义反序列化
type SqlQueryFieldData = []byte

// SqlQueryRowData 一行查询结果，列名 -> 原始字节
type SqlQueryRowData = map[string]SqlQueryFieldData

// MysqlTable 业务表元信息。
// 通常一张表对应一个固定的行结构 + 一套预生成的 SQL（避免每次拼字符串）。
type MysqlTable interface {
	// GetRowObjectCreateFunc 返回查询时用于创建行对象的工厂函数。
	GetRowObjectCreateFunc() func() interface{}
	// GetInsertSql 返回预生成的 INSERT SQL。
	GetInsertSql() string
	// GetReplaceSql 返回预生成的 INSERT ... ON DUPLICATE KEY UPDATE 风格 SQL。
	GetReplaceSql() string
	// GetSelectSql 返回预生成的 SELECT SQL。
	GetSelectSql() string
	// GetUpdateSql 返回预生成的 UPDATE SQL。
	GetUpdateSql() string
}

// MysqlClient 单个 MySQL 连接池实例。底层是 sqlx.DB（已内置连接池）。
type MysqlClient interface {
	// GetMysql 返回底层 sqlx.DB 连接池。
	GetMysql() *sqlx.DB
	// Ping 检查 MySQL 连接是否可用。
	Ping() error
	// InitTable 注册表元信息；rowObjectCreateFn 用于查询时构造目标结构体
	InitTable(tableName string, rowObjectCreateFn func() interface{}) error
	// GetTable 返回指定表名的表元信息。
	GetTable(tableName string) (MysqlTable, bool)
}

// MysqlMgr 多 MySQL 连接池管理器。
// 通过命名（name）支持同时连接多个库（例如主库 + 日志库）；name 为空时使用 OpenDefault 注册的默认连接。
type MysqlMgr interface {
	// OpenDefault 打开默认 MySQL 连接池。
	OpenDefault(addr string, dbPoolSize int32) error
	// Open 以指定名称打开 MySQL 连接池。
	Open(name string, addr string, dbPoolSize int32) error
	// CloseDefault 关闭默认 MySQL 连接池。
	CloseDefault()
	// Close 关闭指定名称的 MySQL 连接池。
	Close(name string)
	// CloseAll 关闭全部 MySQL 连接池。
	CloseAll()
	// Ping 检查指定名称的 MySQL 连接池是否可用。
	Ping(name string) error
	// GetDefaultClient 返回默认 MySQL 客户端。
	GetDefaultClient() (MysqlClient, bool)
	// GetClient 返回指定名称的 MySQL 客户端。
	GetClient(name string) (MysqlClient, bool)
}

// MysqlHandle 无状态的 SQL 操作集合。
// 设计上把底层 *sqlx.DB 作为第一个参数显式传入，是为了让同一组操作可以复用到不同连接池
// （主从、分库、影子库），避免操作方法和 client 实例耦合。
type MysqlHandle interface {
	// NamedInsert 使用 :name 命名参数执行插入，并返回自增 ID。
	NamedInsert(sqldb *sqlx.DB, query string, mapper interface{}) (int64, error)
	// NamedInsertContext 使用 :name 命名参数执行插入，并返回自增 ID。
	NamedInsertContext(ctx context.Context, sqldb *sqlx.DB, query string, mapper interface{}) (int64, error)
	// NamedExec 使用 :name 命名参数执行写 SQL，并返回影响行数。
	NamedExec(sqldb *sqlx.DB, query string, mapper interface{}) (int64, error)
	// NamedExecContext 使用 :name 命名参数执行写 SQL，并返回影响行数。
	NamedExecContext(ctx context.Context, sqldb *sqlx.DB, query string, mapper interface{}) (int64, error)
	// NamedQuery 使用 :name 命名参数执行查询，并通过 rowObjectCreateFn 构造结果对象。
	NamedQuery(sqldb *sqlx.DB, query string, mapper interface{}, rowObjectCreateFn func() interface{}) ([]interface{}, error)
	// NamedQueryContext 使用 :name 命名参数执行查询，并通过 rowObjectCreateFn 构造结果对象。
	NamedQueryContext(ctx context.Context, sqldb *sqlx.DB, query string, mapper interface{}, rowObjectCreateFn func() interface{}) ([]interface{}, error)
	// Insert 使用 ? 位置参数执行插入，并返回自增 ID。
	Insert(sqldb *sqlx.DB, query string, args ...interface{}) (int64, error)
	// InsertContext 使用 ? 位置参数执行插入，并返回自增 ID。
	InsertContext(ctx context.Context, sqldb *sqlx.DB, query string, args ...interface{}) (int64, error)
	// Exec 使用 ? 位置参数执行写 SQL，并返回影响行数。
	Exec(sqldb *sqlx.DB, query string, args ...interface{}) (int64, error)
	// ExecContext 使用 ? 位置参数执行写 SQL，并返回影响行数。
	ExecContext(ctx context.Context, sqldb *sqlx.DB, query string, args ...interface{}) (int64, error)
	// Query 执行查询并以列名到原始字节的映射返回每行数据。
	Query(sqldb *sqlx.DB, query string) ([]SqlQueryRowData, error)
	// QueryContext 执行查询并以列名到原始字节的映射返回每行数据。
	QueryContext(ctx context.Context, sqldb *sqlx.DB, query string) ([]SqlQueryRowData, error)
	// Select 使用 sqlx 风格将查询结果绑定到 dest，dest 通常是结构体切片指针。
	Select(sqldb *sqlx.DB, dest interface{}, query string, args ...interface{}) error
	// SelectContext 使用 sqlx 风格将查询结果绑定到 dest，dest 通常是结构体切片指针。
	SelectContext(ctx context.Context, sqldb *sqlx.DB, dest interface{}, query string, args ...interface{}) error
	// BeginTransaction 开启事务，调用方必须显式 Commit 或 Rollback。
	BeginTransaction(sqldb *sqlx.DB) (*sql.Tx, error)
	// BeginTransactionContext 开启事务，调用方必须显式 Commit 或 Rollback。
	BeginTransactionContext(ctx context.Context, sqldb *sqlx.DB) (*sql.Tx, error)
	// InsertTransaction 在事务中执行插入，并返回自增 ID。
	InsertTransaction(tx *sql.Tx, query string, args ...interface{}) (int64, error)
	// ExecTransaction 在事务中执行写 SQL。
	ExecTransaction(tx *sql.Tx, query string, args ...interface{}) error
}
