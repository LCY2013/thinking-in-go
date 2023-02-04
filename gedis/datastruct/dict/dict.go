package dict

// Consumer is used to traversal dict, if it returns false the traversal will be break
type Consumer func(key string, val any) bool

// Dict is interface of a key-value data structure
type Dict interface {
	Get(key string) (val any, exists bool)
	Len() int
	Put(key string, val any) (result int)
	PutIfAbsent(key string, val any) (result int)
	PutIfExists(key string, val any) (result int)
	Remove(key string) (result int)
	ForEach(consumer Consumer)
	Keys() []string
	RandomKeys(limit int) []string
	RandomDistinctKeys(limit int) []string
	Clear()
}
