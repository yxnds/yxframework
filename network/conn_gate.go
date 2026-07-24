package network

// SvrStatRsp 通用的 HTTP /status 响应体。
// Code=0 表示健康，其他为业务自定义错误码；Result 携带详细信息（通常是 JSON 字符串）。
type SvrStatRsp struct {
	// Code 表示状态码，0 通常表示成功或健康。
	Code int32 `json:"code"`
	// Result 表示状态详情，通常为 JSON 字符串。
	Result string `json:"result"`
}

// ConnGate 管理面/网关型监听器。
// 与 ConnMgr 不同，ConnGate 通常用于 HTTP 服务（健康检测、运维接口）这类不需要长连接抽象的场景。
type ConnGate interface {
	// Listen 启动网关监听；protoTye 在 HTTP 实现中通常被忽略，address 形如 ":8080"。
	Listen(protoTye string, address string) error
	// Close 优雅关闭网关监听。
	Close() error
}
