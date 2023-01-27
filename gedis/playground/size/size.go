package main

import "fmt"

// main 切片创建
// go build -gcflags -S size.go
/*func main() {
	arr := []int{1, 2, 3, 4}
	fmt.Println(arr)
}*/

/*
切片的创建过程
0x0018 00024 (./gedis/playground/size/size.go:8)   MOVD    $type.[4]int(SB), R0
0x0020 00032 (./gedis/playground/size/size.go:8)   PCDATA  $1, ZR
0x0020 00032 (./gedis/playground/size/size.go:8)   CALL    runtime.newobject(SB)
0x0024 00036 (./gedis/playground/size/size.go:8)   MOVD    $1, R1
0x0028 00040 (./gedis/playground/size/size.go:8)   MOVD    R1, (R0)
0x002c 00044 (./gedis/playground/size/size.go:8)   MOVD    $2, R2
0x0030 00048 (./gedis/playground/size/size.go:8)   MOVD    R2, 8(R0)
0x0034 00052 (./gedis/playground/size/size.go:8)   MOVD    $3, R2
0x0038 00056 (./gedis/playground/size/size.go:8)   MOVD    R2, 16(R0)
0x003c 00060 (./gedis/playground/size/size.go:8)   MOVD    $4, R2
0x0040 00064 (./gedis/playground/size/size.go:8)   MOVD    R2, 24(R0)
*/

func main() {
	m := make(map[string]int, 10)
	fmt.Println(m)
}

/*
map创建过程
0x0018 00024 (./gedis/playground/size/size.go:28)  MOVD    $type.map[string]int(SB), R0
0x0020 00032 (./gedis/playground/size/size.go:28)  MOVD    $10, R1
0x0024 00036 (./gedis/playground/size/size.go:28)  MOVD    ZR, R2
0x0028 00040 (./gedis/playground/size/size.go:28)  PCDATA  $1, ZR
0x0028 00040 (./gedis/playground/size/size.go:28)  CALL    runtime.makemap(SB)
*/
