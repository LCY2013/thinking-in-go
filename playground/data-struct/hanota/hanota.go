package main

import "fmt"

func main() {
	src := []int{1, 2, 3, 4, 5}
	tmp, dst := make([]int, 0, len(src)), make([]int, 0, len(src))
	hanota(len(src), &src, &tmp, &dst)
	fmt.Println("src", src)
	fmt.Println("tmp", tmp)
	fmt.Println("dst", dst)
}

// hanota 汉诺塔问题
func hanota(n int, src, tmp, dst *[]int) {
	fmt.Println("src", *src)
	fmt.Println("tmp", *tmp)
	fmt.Println("dst", *dst)

	if n == 1 {
		// f(1) 将src上的最后一个盘子移动到dst
		t := make([]int, 0)
		t = append(t, (*src)[0])
		t = append(t, *dst...)
		*dst = t
		*src = (*src)[1:]
		return
	}

	// f(n-1) 子问题, 借助dst将src上的n-1个盘子移动到tmp
	hanota(n-1, src, dst, tmp)
	// f(1) 将src上的最后一个盘子移动到dst
	t := make([]int, 0)
	t = append(t, (*src)[0])
	t = append(t, *dst...)
	*dst = t
	*src = (*src)[1:]
	// f(n-1) 子问题，借助src将tmp上的n-1个盘子移动到dst
	hanota(n-1, tmp, src, dst)
}
