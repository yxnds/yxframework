package yx

import "errors"

var (
	// ErrAesDataLenZero 表示 AES 输入数据长度为 0。
	ErrAesDataLenZero = errors.New("data len is 0")
	// ErrAesKeyLenZero 表示 AES 密钥长度为 0。
	ErrAesKeyLenZero = errors.New("key len is 0")
	// ErrAesUnpaddingTooMuch 表示 AES 解密后填充长度超过数据长度，通常是密钥错误或密文被篡改。
	ErrAesUnpaddingTooMuch = errors.New("unpadding count more than data len")
	// ErrAesCiphertextInvalid 表示 AES 密文格式无效。
	ErrAesCiphertextInvalid = errors.New("aes ciphertext invalid")
	// ErrAesPaddingInvalid 表示 AES 填充数据无效。
	ErrAesPaddingInvalid = errors.New("aes padding invalid")

	// ErrCrc16DataLenZero 表示 CRC16 输入数据长度为 0。
	ErrCrc16DataLenZero = errors.New("data len is 0")

	// ErrEcdsaDataLenZero 表示 ECDSA 输入数据长度为 0。
	ErrEcdsaDataLenZero = errors.New("data len is 0")
	// ErrEcdsaKeyLenZero 表示 ECDSA 密钥长度为 0。
	ErrEcdsaKeyLenZero = errors.New("key len is 0")
	// ErrEcdsaParsePubKeyFail 表示 ECDSA 公钥解析失败。
	ErrEcdsaParsePubKeyFail = errors.New("parse public key failed")
	// ErrEcdsaInvalidCurve 表示 ECDSA 曲线类型无效。
	ErrEcdsaInvalidCurve = errors.New("invalid ecdsa curve")
	// ErrEcdsaInvalidPem 表示 ECDSA PEM 文本无效。
	ErrEcdsaInvalidPem = errors.New("invalid ecdsa pem")
	// ErrEcdsaInvalidSignature 表示 ECDSA 签名数据无效。
	ErrEcdsaInvalidSignature = errors.New("invalid ecdsa signature")
	// ErrEcdsaUnsupportedHash 表示 ECDSA 不支持指定哈希算法。
	ErrEcdsaUnsupportedHash = errors.New("unsupported ecdsa hash")

	// ErrEvtWaitTimeout 表示等待事件通知超时。
	ErrEvtWaitTimeout = errors.New("event wait time out")
	// ErrEvtClosed 表示事件已经关闭后仍调用 Wait。
	ErrEvtClosed = errors.New("event closed")

	// ErrLckTryLockFail 表示 FastLock 自旋超过上限仍未拿到锁。
	ErrLckTryLockFail = errors.New("try lock failed")

	// ErrLoggerStopped 表示日志器已经停止。
	ErrLoggerStopped = errors.New("logv2: logger is stopped")
	// ErrInvalidConfig 表示日志配置无效。
	ErrInvalidConfig = errors.New("logv2: invalid config")

	// ErrMailIsNil 表示邮件参数对象为空。
	ErrMailIsNil = errors.New("mail is nil")
	// ErrMailInvalidFrom 表示发件人信息无效。
	ErrMailInvalidFrom = errors.New("mail from is invalid")
	// ErrMailInvalidTo 表示收件人列表无效。
	ErrMailInvalidTo = errors.New("mail to is invalid")
	// ErrMailInvalidHost 表示 SMTP 主机地址无效。
	ErrMailInvalidHost = errors.New("mail host is invalid")
	// ErrMailInvalidPort 表示 SMTP 端口无效。
	ErrMailInvalidPort = errors.New("mail port is invalid")
	// ErrMailInvalidAddress 表示邮件地址格式无效。
	ErrMailInvalidAddress = errors.New("mail address is invalid")
	// ErrMailInvalidAttach 表示邮件附件配置无效。
	ErrMailInvalidAttach = errors.New("mail attachment is invalid")
	// ErrMailInvalidTimeout 表示邮件发送超时时间无效。
	ErrMailInvalidTimeout = errors.New("mail timeout is invalid")
	// ErrMailDialerIsNil 表示邮件拨号器为空。
	ErrMailDialerIsNil = errors.New("mail dialer is nil")
	// ErrMailSenderIsNil 表示邮件发送器为空。
	ErrMailSenderIsNil = errors.New("mail sender is nil")
	// ErrMailMessageBuildFail 表示邮件消息构建失败。
	ErrMailMessageBuildFail = errors.New("mail message build failed")

	// ErrObsMsgEmpty 表示观察者消息名为空。
	ErrObsMsgEmpty = errors.New("msgName is empty")

	// ErrQueEmptyQueue 表示在空队列上执行 Dequeue。
	ErrQueEmptyQueue = errors.New("empty queue")
	// ErrQueFullQueue 表示限长队列容量已满。
	ErrQueFullQueue = errors.New("full queue")
	// ErrQueIsNil 表示队列对象为空。
	ErrQueIsNil = errors.New("queue is nil")
	// ErrQueSizeOverflow 表示队列大小溢出。
	ErrQueSizeOverflow = errors.New("queue size overflow")

	// ErrRsaDataLenZero 表示 RSA 输入数据长度为 0。
	ErrRsaDataLenZero = errors.New("data len is 0")
	// ErrRsaKeyLenZero 表示 RSA 密钥长度为 0。
	ErrRsaKeyLenZero = errors.New("key len is 0")
	// ErrRsaInvalidPem 表示 RSA PEM 文本无效。
	ErrRsaInvalidPem = errors.New("invalid rsa pem")
	// ErrRsaInvalidKey 表示 RSA 密钥无效。
	ErrRsaInvalidKey = errors.New("invalid rsa key")
	// ErrRsaInvalidCiphertext 表示 RSA 密文无效。
	ErrRsaInvalidCiphertext = errors.New("invalid rsa ciphertext")
	// ErrRsaMessageTooLong 表示 RSA 明文长度超过当前密钥和填充方式允许的上限。
	ErrRsaMessageTooLong = errors.New("rsa message too long")
	// ErrRsaUnsupportedHash 表示 RSA 不支持指定哈希算法。
	ErrRsaUnsupportedHash = errors.New("unsupported rsa hash")

	// ErrSetKeyIsNil 表示集合键为空。
	ErrSetKeyIsNil = errors.New("key is nil")
	// ErrSetNotIntKey 表示向 IntSet 传入了非整数键。
	ErrSetNotIntKey = errors.New("not int key")
	// ErrSetNotUintKey 表示向 UintSet 传入了非无符号整数键。
	ErrSetNotUintKey = errors.New("not uint key")
	// ErrSetNotFloatKey 表示向 FloatSet 传入了非浮点键。
	ErrSetNotFloatKey = errors.New("not float key")
	// ErrSetIsNil 表示集合对象为空。
	ErrSetIsNil = errors.New("set is nil")
	// ErrSetKeyUncomparable 表示集合键不可比较，不能作为 map 键。
	ErrSetKeyUncomparable = errors.New("key is uncomparable")

	// ErrSbNoSpaceToRead 表示缓冲区剩余写空间不足以容纳新数据。
	ErrSbNoSpaceToRead = errors.New("no space to read data")
	// ErrSbBuffEmpty 表示缓冲区为空。
	ErrSbBuffEmpty = errors.New("buffer empty")
	// ErrSbReachEnd 表示缓冲区读写位置已经到达末尾。
	ErrSbReachEnd = errors.New("buffer reach end")
	// ErrSbIsNil 表示 SimpleBuffer 对象为空。
	ErrSbIsNil = errors.New("simple buffer is nil")
	// ErrSbInvalidCapacity 表示 SimpleBuffer 容量无效。
	ErrSbInvalidCapacity = errors.New("simple buffer capacity is invalid")

	// ErrSfInvalidWorkerId 表示 workerId 超出允许范围。
	ErrSfInvalidWorkerId = errors.New("invalid worker id")
	// ErrSfNotInitialized 表示 Snowflake 尚未初始化。
	ErrSfNotInitialized = errors.New("snowflake is not initialized")
	// ErrSfClockBackwards 表示系统时钟发生回拨且超过可等待范围。
	ErrSfClockBackwards = errors.New("snowflake clock moved backwards")

	// ErrUriInvalidEscape 表示 URI 转义序列无效。
	ErrUriInvalidEscape = errors.New("invalid uri escape")
	// ErrUriWriterNil 表示 HTTP 响应写入器为空。
	ErrUriWriterNil = errors.New("http response writer is nil")
)
