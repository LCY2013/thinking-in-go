package main

import "fmt"

/*
方法与函数等价关系如下：
func (t T) M1() <=> F1(t T)
func (t *T) M2() <=> F2(t *T)

首先，当 receiver 参数的类型为 T 时：
当我们选择以 T 作为 receiver 参数类型时，M1 方法等价转换为F1(t T)。我们知道，Go 函数的参数采用的是值拷贝传递，也就是说，F1 函数体中的 t 是 T 类型实例的一个副本。这样，我们在 F1 函数的实现中对参数 t 做任何修改，都只会影响副本，而不会影响到原 T 类型实例。
据此我们可以得出结论：当我们的方法 M1 采用类型为 T 的 receiver 参数时，代表 T 类型实例的 receiver 参数以值传递方式传递到 M1 方法体中的，实际上是 T 类型实例的副本，M1 方法体中对副本的任何修改操作，都不会影响到原 T 类型实例。

第二，当 receiver 参数的类型为 *T 时：
当我们选择以 *T 作为 receiver 参数类型时，M2 方法等价转换为F2(t *T)。同上面分析，我们传递给 F2 函数的 t 是 T 类型实例的地址，这样 F2 函数体中对参数 t 做的任何修改，都会反映到原 T 类型实例上。
*/
type change struct {
	name string
}

func (c *change) change() {
	c.name = fmt.Sprintf("%s-%s", "changed", c.name)
}

type unChange struct {
	name string
}

func (u unChange) unChange() {
	u.name = fmt.Sprintf("%s-%s", "unChanged", u.name)
}

func main() {
	c := change{"1"}
	u := unChange{"2"}
	c.change()
	fmt.Println(c.name)

	u.unChange()
	fmt.Println(u.name)
	unChange.unChange(u)
}
