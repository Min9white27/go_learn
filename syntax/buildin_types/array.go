package main

import "fmt"

func Array() {
	a1 := [3]int{9, 8, 7}
	fmt.Printf("a1: %v,len=%d,cap=%d\n", a1, len(a1), cap(a1)) //len 长度 cap 容量

	a2 := [3]int{9, 8}
	fmt.Printf("a2: %v,len=%d,cap=%d\n", a2, len(a2), cap(a2))

	var a3 [3]int
	fmt.Printf("a3: %v,len=%d,cap=%d\n", a3, len(a3), cap(a3))

	//数组不支持切片操作
	//a3 = append(a3,1)

	//按下标索引，如果编译器能判断出来下标越界，那么会编译错误
	//如果不能，那么运行时会报错，出现 panic
	fmt.Printf("a1[1]: %d", a1[1])

	//使用 for range 循环
	//基本不用数组
}
