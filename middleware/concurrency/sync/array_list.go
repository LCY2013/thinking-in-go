package sync

// ArrayList 基于切片实现
type ArrayList[T any] struct {
	values []T
}

func NewArrayList[T any](cap int) *ArrayList[T] {
	panic("implement me")
}

// NewArrayListOf 直接使用 ts，而不会执行复制
func NewArrayListOf[T any](ts []T) *ArrayList[T] {
	return &ArrayList[T]{
		values: ts,
	}
}

func (a *ArrayList[T]) Get(index int) (T, error) {
	// TODO implement me
	panic("implement me")
}

func (a *ArrayList[T]) Append(t T) error {
	// TODO implement me
	panic("implement me")
}

// Add 在ArrayList下标为index的位置插入一个元素
// 当index等于ArrayList长度等同于append
func (a *ArrayList[T]) Add(index int, t T) error {
	if index < 0 || index > len(a.values) {
		return newErrIndexOutOfRange(len(a.values), index)
	}
	a.values = append(a.values, t)
	copy(a.values[index+1:], a.values[index:])
	a.values[index] = t
	return nil
}

func (a *ArrayList[T]) Set(index int, t T) error {
	// TODO implement me
	panic("implement me")
}

func (a *ArrayList[T]) Delete(index int) (T, error) {
	// TODO implement me
	panic("implement me")
}

func (a *ArrayList[T]) Len() int {
	// TODO implement me
	panic("implement me")
}

func (a *ArrayList[T]) Cap() int {
	return cap(a.values)
}

func (a *ArrayList[T]) Range(fn func(index int, t T) error) error {
	for key, value := range a.values {
		e := fn(key, value)
		if e != nil {
			return e
		}
	}
	return nil
}

func (a *ArrayList[T]) AsSlice() []T {
	slice := make([]T, len(a.values))
	copy(slice, a.values)
	return slice
}
