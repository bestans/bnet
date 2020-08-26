package tcp

const (
	//缓存写buffer数量
	cacheWriteBufferNum   = 10
	//缓存写buffer大小（预计90%以上的消息均在1024字节以内）
	cacheWriteBufferSize  = 1024
)
