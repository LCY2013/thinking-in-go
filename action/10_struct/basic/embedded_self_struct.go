package basic

// 不合法
//type T struct {
//	t T
//}

// T 结构体内嵌自己需要是指针类型，指针可以确定一个machine word大小
type T struct {
	t *T
}

// 不合法
//type T1 struct {
//	t2 T2
//}
//
//type T2 struct {
//	t1 T1
//}

type E struct {
	e  *E
	se []E
	me map[string]E
}
