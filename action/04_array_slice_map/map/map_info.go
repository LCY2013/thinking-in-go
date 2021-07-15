package maps

import "fmt"

// MapFunc 字典函数，用于演示字典函数的示例
func MapFunc() {
	// 使用内置函数 make 创建字典映射, 键是 string，值是 int
	dict := make(map[string]int)
	fmt.Println(dict)

	// 使用字面量创建 字典映射，键、值都是string
	dict1 := map[string]string{"Red": "#da1337", "Orange": "#e95a22"}
	fmt.Println(dict1)

	// 使用字面量创建 一个映射，使用 字符切片作为键
	// # fufeng.org/sample04/map
	//./map_info.go:17:11: invalid map key type []string
	// dict2 := map[[]string]int{}

	// 使用字面量创建 一个映射，使用 字符切片作为值
	dict3 := map[int][]string{}
	fmt.Println(dict3)

	// 使用字面量声明一个空的 映射，用来存储颜色以及颜色对应的十六进制代码
	dict4 := map[string]string{}
	// 定义一个新的颜色映射信息
	dict4["Red"] = "#da1337"

	// 通过声明创建一个nil的映射
	var colors map[string]string
	// 将 Red 加入到映射中
	// panic: assignment to entry in nil map
	// colors["Red"] = "#da1337"

	// 判断映射中某个键是否存在
	// 获取键 Blue 对应的值
	value, exists := colors["Blue"]
	if exists {
		fmt.Println(value)
	}

	// 通过返回的值是不是默认值来判断是不是不存在
	v := colors["Blue"]
	if v != "" {
		fmt.Println(v)
	}

	// 使用range来迭代 映射的数据，返回的是键值
	color := map[string]string{
		"AliceBlue":   "#f0f8ff",
		"Coral":       "#ff7F50",
		"DarkGray":    "#a9a9a9",
		"ForestGreen": "#228b22",
	}
	for key, value := range color {
		fmt.Printf("key: %s, value: %s\n", key, value)
	}

	// 从映射中删除一项
	delete(color, "Coral")
	for key, value := range color {
		fmt.Printf("delete -> key: %s, value: %s\n", key, value)
	}
}

// FuncToFuncMap 在函数间传递映射
// 映射在函数间传递不会制造出该映射的一个副本，也就是在传递的函数里面执行修改映射的操作会改动传入的那个函数中映射的数据
func FuncToFuncMap() {
	// 使用字面量声明一个空的 映射，用来存储颜色以及颜色对应的十六进制代码
	colors := map[string]string{
		"AliceBlue":   "#f0f8ff",
		"Coral":       "#ff7F50",
		"DarkGray":    "#a9a9a9",
		"ForestGreen": "#228b22",
	}

	// 显示映射里面的所有键值对信息
	for key, value := range colors {
		fmt.Printf("key: %s, value: %s\n", key, value)
	}

	fmt.Println()

	// 调用函数移除指定的键
	removeColors(colors, "Coral")

	// 显示映射里面所有键值对信息
	for key, value := range colors {
		fmt.Printf("key: %s, value: %s\n", key, value)
	}

}

// removeColors 移除映射中指定的颜色信息
func removeColors(colors map[string]string, removeKey string) {
	delete(colors, removeKey)
}
