package main

import "fmt"

func main() {
	//inline()
	//outline()
	//deferOrder()
	m5()
}

func deferOrder() {
	defer m7()
	defer m8()
}

func outline() {
	m5()
	m7()
	m6()
}

func inline() {
	m1()
	m2()
	m4()
	m3()
}

func m1() {
	defer func() {
		panic("m1")
	}()
}

func m2() {
	defer func() {
		panic("m2")
	}()
}

func m3() {
	defer func() {
		panic("m3")
	}()
}

func m4() {
	fmt.Println("m4")
}

func m5() {
	defer func() {
		fmt.Println("before")
	}()
	defer func() {
		if r := recover(); r != nil {
		}
	}()
	defer func() {
		fmt.Println("after")
	}()
	panic("m5")
}

func m6() {
	defer func() {
		if r := recover(); r != nil {
		}
	}()
	panic("m6")
}

func m7() {
	fmt.Println("m7")
}

func m8() {
	fmt.Println("m8")
}
