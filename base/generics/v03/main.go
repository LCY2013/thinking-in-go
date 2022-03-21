package main

import (
	"fmt"
	"github.com/samber/lo"
)

/*
在 Golang 支持泛型之前，实现像 lodash.js 这样一套适配多种数据类型的完整的工具库是非常不容易的。有一些开源库通过其他方式实现了部分功能，大致有三种方案：

纯手撸 - 毫无疑问，这种方式是最不优雅的，需要对每种类型进行开发，需要做很多无聊的工作。
代码生成 - 通过脚本辅助生成针对不同类型的工具函数，比如 go-dash/slice。
使用反射 - 这种方式可以实现目的，但是反射会带来较大复杂度和造成运行时性能的下降。go-funk 和robpike/filter都是通过该种方式实现的工具库。
somber/lo 是一个基于 Golang 泛型实现的的 lodash 风格工具库，比较好的避免了上面的问题。

somber/lo 包含了非常多的方法，主要可以划分为以下几类：
1、slice 辅助方法
2、map 辅助方法
3、tuples 辅助方法
4、多个集合之间计算辅助方法
5、搜索查询辅助方法
6、其他函数式编程辅助方法等

*/

// uniqSlice 以切片去重举例
func uniqSlice() {
	names1 := lo.Uniq[string]([]string{"Samuel", "Marc", "Samuel"})
	// 调用非常简单，并且在大多数情况下，可以省略类型的指定:
	names2 := lo.Uniq([]string{"Samuel", "Marc", "Samuel"})
	// []string{"Samuel", "Marc"}
	fmt.Println(names1, names2)
}

// filterSlice 再比如过滤掉切片中不符合规则的元素
func filterSlice() {
	even := lo.Filter([]int{1, 2, 3, 4}, func(v int, i int) bool {
		return v%2 == 0
	})
	fmt.Println(even)
}

// main somber/lo 基于泛型包装了非常多的工具方法，可以大大节省我们的开发时间，避免重复开发，提升效率。
// 但是该库开源至今才两周，可能会有一些问题缺陷存在其中，线上使用还需要谨慎一些。
func main() {
	uniqSlice()
}
