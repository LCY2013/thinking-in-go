package v2

// alterCounter 用于保存告警计数
type alterCounter int

// New 返回一个未公开的 alterCounter 类型的值
func New(value int) alterCounter {
	return alterCounter(value)
}
