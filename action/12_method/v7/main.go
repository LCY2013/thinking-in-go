package main

import (
	"fmt"
	"reflect"
)

type Interface interface {
	M1()
	M2()
}
type T struct{}

func (t T) M1() {}

func (t *T) M2() {}

func main() {
	var t T
	var pt *T
	var i Interface
	i = pt

	// T 没有实现 Interface 类型方法列表中的 M2，因此类型 T 的实例 t 不能赋值给 Interface 变量。
	//i = t // cannot use t (type T) as type Interface in assignment: T does not implement Interface (M2 method has pointer receiver)
	i.M2()

	// 为什么 *T 类型的 pt 可以被正常赋值给 Interface 类型变量 i，而 T 类型的 t 就不行呢？
	// 如果说 T 类型是因为只实现了 M1 方法，未实现 M2 方法而不满足 Interface 类型的要求，那么 *T 类型也只是实现了 M2 方法，并没有实现 M1 方法啊？

	// 方法集合也是用来判断一个类型是否实现了某接口类型的唯一手段，可以说，“方法集合决定了接口实现”。

	// 什么是类型的方法集合呢？
	// Go 中任何一个类型都有属于自己的方法集合，或者说方法集合是 Go 类型的一个“属性”。
	// 但不是所有类型都有自己的方法呀，比如 int 类型就没有。
	// 所以，对于没有定义方法的 Go 类型，我们称其拥有空方法集合。

	// 接口类型相对特殊，它只会列出代表接口的方法列表，不会具体定义某个方法，它的方法集合就是它的方法列表中的所有方法，可以一目了然地看到。因此，下面重点讲解的是非接口类型的方法集合。

	// 为了方便查看一个非接口类型的方法集合，这里提供了一个函数 dumpMethodSet，用于输出一个非接口类型的方法集合：
	fmt.Println("-------dumpMethodSet interface--------")
	dumpMethodSet(i)
	fmt.Println("-------dumpMethodSet pointer--------")
	dumpMethodSet(pt)
	fmt.Println("-------dumpMethodSet instance--------")
	dumpMethodSet(t)
	fmt.Println("-------dumpMethodSet int--------")
	var num int
	dumpMethodSet(num)

	// 结论：Go 语言规定，*T 类型的方法集合包含所有以 *T 为 receiver 参数类型的方法，以及所有以 T 为 receiver 参数类型的方法。这就是这个示例中为何 *T 类型的方法集合包含四个方法的原因。

	// 选择 receiver 参数类型的第三个原则
	// 这个原则的选择依据就是 T 类型是否需要实现某个接口。
	// 如果 T 类型需要实现某个接口，那我们就要使用 T 作为 receiver 参数的类型，来满足接口类型方法集合中的所有方法。
	// 如果 T 不需要实现某一接口，但 *T 需要实现该接口，那么根据方法集合概念，*T 的方法集合是包含 T 的方法集合的，这样我们在确定 Go 方法的 receiver 的类型时，参考原则一和原则二就可以了。
}

func dumpMethodSet(i interface{}) {
	dynTyp := reflect.TypeOf(i)
	if dynTyp == nil {
		fmt.Printf("there is no dynamic type\n")
		return
	}
	n := dynTyp.NumMethod()
	if n == 0 {
		fmt.Printf("%s's method set is empty!\n", dynTyp)
		return
	}
	fmt.Printf("%s's method set:\n", dynTyp)
	for j := 0; j < n; j++ {
		fmt.Println("-", dynTyp.Method(j).Name)
	}
	fmt.Printf("\n")
}
