package mapstruct

// Ptr 转指针
func Ptr[T any](t T) *T {
	return &t
}
