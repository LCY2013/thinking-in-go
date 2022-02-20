package main

/*
 receiver 参数名字要保证唯一外，Go 语言对 receiver 参数的基类型也有约束，那就是 receiver 参数的基类型本身不能为指针类型或接口类型。
 下面的事例都存在问题，不满足receiver 的参数要求。
*/

//type MyInt *int
//func (r MyInt) String() string { // r的基类型为MyInt，编译器报错：invalid receiver type MyInt (MyInt is a pointer type)
//	return fmt.Sprintf("%d", *(*int)(r))
//}
//type MyReader io.Reader
//func (r MyReader) Read(p []byte) (int, error) { // r的基类型为MyReader，编译器报错：invalid receiver type MyReader (MyReader is an interface type)
//	return r.Read(p)
//}

/*
Go 对方法声明的位置也是有约束的，Go 要求，方法声明要与 receiver 参数的基类型声明放在同一个包内。基于这个约束，我们还可以得到两个推论。
第一个推论：我们不能为原生类型（诸如 int、float64、map 等）添加方法。
比如，下面的代码试图为 Go 原生类型 int 增加新方法 Foo，这样做，Go 编译器会报错：
func (i int) Foo() string { // 编译器报错：cannot define new methods on non-local type int
    return fmt.Sprintf("%d", i)
}

第二个推论：不能跨越 Go 包为其他包的类型声明新方法。
比如，下面的代码试图跨越包边界，为 Go 标准库中的 http.Server 类型添加新方法 Foo，这样做，Go 编译器同样会报错：
import "net/http"
func (s http.Server) Foo() { // 编译器报错：cannot define new methods on non-local type http.Server
}
*/

type T struct{}

func (t T) M(n int) {
}

func main() {
	var t T
	t.M(1) // 通过类型T的变量实例调用方法M
	p := &T{}
	p.M(2) // 通过类型*T的变量实例调用方法M

	// 方法 M 是类型 T 的方法，那为什么通过 *T 类型变量也可以调用 M 方法呢？
}
