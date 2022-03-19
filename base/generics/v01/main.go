package main

import "fmt"

type testInt int64

// Number 在这个示例中，可以看到我们将testInt作为自定义类型传入了泛型方法中，在这种情况下，如果不给Number中的int64加~，这里就会报错。加上~之后代表以int64为基本类型的自定义类型也可以通过泛型约束。
type Number interface {
	~int64
}

func SumIntsOrFloats[K comparable, V Number](m map[K]V) V {
	var s V
	for _, v := range m {
		s += v
	}
	return s
}

func main() {
	ints := map[string]testInt{
		"first":  34,
		"second": 12,
	}

	fmt.Printf("Genrics Sums: %v\n", SumIntsOrFloats(ints))
}
