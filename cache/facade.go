package cache

// Cache 缓存接口.
type Cache interface {
	Get(key string) // 尝试获取，阻塞
	TryGet()        // 非阻塞
}
