package db

import "go.mongodb.org/mongo-driver/mongo"

const (
	// DefaultMongo 是未指定名称时使用的默认 Mongo 客户端名。
	DefaultMongo = "default"

	// MongoUserFieldId 是 Mongo 文档主键字段。
	MongoUserFieldId = "_id"
	// MongoUserFieldModName 是 Mongo KV 模式下的模块名字段。
	MongoUserFieldModName = "mod_name"
	// MongoUserFieldData 是 Mongo 文档中保存主体业务数据的字段。
	MongoUserFieldData = "data"
	// MongoUserFieldDataVer 是 Mongo 文档中用于乐观锁的数据版本字段。
	MongoUserFieldDataVer = "data_ver"
	// MongoUserFieldAreaId 是 Mongo 文档中保存大区 ID 的字段。
	MongoUserFieldAreaId = "area_id"
	// MongoUserFieldSfId 是 Mongo 文档中保存区服或实例 ID 的字段。
	MongoUserFieldSfId = "sf_id"
	// MongoUserFieldCreateTime 是 Mongo 文档创建时间字段。
	MongoUserFieldCreateTime = "create_time"
	// MongoUserFieldUpdateTime 是 Mongo 文档更新时间字段。
	MongoUserFieldUpdateTime = "update_time"
	// MongoUserFieldDeleteTime 是 Mongo 文档删除时间字段。
	MongoUserFieldDeleteTime = "delete_time"
)

// MongoClient 单个 MongoDB 连接实例（一个 mongo.Client 内含连接池）。
// 通过 mainDb 指定的默认库快速取 collection；也支持跨库访问。
type MongoClient interface {
	// GetMongo 返回底层 mongo.Client。
	GetMongo() *mongo.Client
	// Ping 检查 MongoDB 连接是否可用。
	Ping() error
	// GetMainDbCollection 从 OpenDefault 或 Open 指定的主库中获取集合。
	GetMainDbCollection(name string) (*mongo.Collection, bool)
	// GetCollectionFromDb 从指定数据库中获取集合。
	GetCollectionFromDb(dbName string, collectionName string) (*mongo.Collection, bool)
}

// MongoMgr 多 MongoDB 连接管理。命名风格与 MysqlMgr 一致。
type MongoMgr interface {
	// OpenDefault 打开默认 MongoDB 客户端。
	OpenDefault(addr string, mainDb string) error
	// Open 以指定名称打开 MongoDB 客户端。
	Open(name string, addr string, mainDb string) error
	// CloseDefault 关闭默认 MongoDB 客户端。
	CloseDefault()
	// Close 关闭指定名称的 MongoDB 客户端。
	Close(name string)
	// CloseAll 关闭全部 MongoDB 客户端。
	CloseAll()
	// Ping 检查指定名称的 MongoDB 客户端是否可用。
	Ping(name string) error
	// GetDefaultClient 返回默认 MongoDB 客户端。
	GetDefaultClient() (MongoClient, bool)
	// GetClient 返回指定名称的 MongoDB 客户端。
	GetClient(name string) (MongoClient, bool)
}

// MongoHandle 提供两套存储模式的读写：
//   - 关系型（Relational）：一个 docId 对应一条文档，data 字段存整体业务数据
//   - KV：一个 (modName, docId) 对应一条文档，按模块名隔离，便于按模块独立扩展
//
// 两套模式都支持 dataVer 乐观锁：读时取出版本号，写时带上原版本号；
// 服务端版本不一致则写入失败，业务侧重读重试，避免并发覆盖。
type MongoHandle interface {
	// SetMongoDataVersion 设置全局默认数据版本，新建文档时使用该版本。
	SetMongoDataVersion(ver int64)
	// GetRelationalMongoData 读取关系型文档。
	//
	// fields 为指定字段；allFields=true 读全部；binFilterKeys 用于过滤掉某些公共字段。
	// 返回值：data 主体、字段映射、dataVer、是否成功。
	GetRelationalMongoData(collection *mongo.Collection, docId uint64, fields []string, allFields bool, binFilterKeys map[string]bool) ([]byte, map[string][]byte, int64, bool)
	// SaveRelationalMongoData 保存关系型文档，dataVer 用于乐观锁校验。
	SaveRelationalMongoData(collection *mongo.Collection, areaId uint32, docId uint64, data []byte, field2DataMap map[string][]byte, dataVer int64) bool
	// GetKVMongoData 按模块名和文档 ID 读取 KV 模式文档。
	GetKVMongoData(collection *mongo.Collection, modName string, docId uint64, fields []string, allFields bool) ([]byte, map[string][]byte, int64, bool)
	// SaveKVMongoData 保存 KV 模式文档，areaOrFightId 用于区分大区或战斗实例。
	SaveKVMongoData(collection *mongo.Collection, modName string, areaOrFightId uint32, docId uint64, data []byte, field2DataMap map[string][]byte, dataVer int64) bool
}
