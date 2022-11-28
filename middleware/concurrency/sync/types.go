package sync

// List 接口
// 该接口只定义清楚各个方法的行为和表现
type List[T any] interface {
	// Get 返回对应下标的元素，
	// 在下标超出范围的情况下，返回错误
	Get(index int) (T, error)
	// Append 在末尾追加元素
	Append(t T) error
	// Add 在特定下标处增加一个新元素
	// 如果下标超出范围，应该返回错误
	Add(index int, t T) error
	// Set 重置 index 位置的值
	// 如果下标超出范围，应该返回错误
	Set(index int, t T) error
	// Delete 删除目标元素的位置，并且返回该位置的值
	// 如果 index 超出下标，应该返回错误
	Delete(index int) (T, error)
	// Len 返回长度
	Len() int
	// Cap 返回容量
	Cap() int
	// Range 遍历 List 的所有元素
	Range(fn func(index int, t T) error) error
	// AsSlice 将 List 转化为一个切片
	// 不允许返回nil，在没有元素的情况下，
	// 必须返回一个长度和容量都为 0 的切片
	// AsSlice 每次调用都必须返回一个全新的切片
	AsSlice() []T
}
