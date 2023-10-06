package main

import "fmt"

func Slice() {
	s1 := []int{1, 2, 3, 4, 5}
	fmt.Printf("s1: %v,len=%d,cap=%d\n", s1, len(s1), cap(s1))

	s2 := make([]int, 3, 4)
	fmt.Printf("s2: %v,len=%d,cap=%d\n", s2, len(s2), cap(s2))

	//{0,0,0,0}
	s3 := make([]int, 4)
	//s3[0]=1
	fmt.Printf("s3: %v,len=%d,cap=%d\n", s3, len(s3), cap(s3))

	s4 := make([]int, 0, 4)
	s4 = append(s4, 1)
	fmt.Printf("s4: %v,len=%d,cap=%d\n", s4, len(s4), cap(s4))

	//在初始化切片时要预估容量，尽量避免扩容，多了浪费内存，少了要进行扩容，还是一样容易浪费内存
	//遇事不决用切片
	//切片的底层是切片
}

func SubSlice() {
	s1 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	s2 := s1[1:3] //左闭右开原则，数字代表下标
	fmt.Printf("s2: %v,len=%d,cap=%d\n", s2, len(s2), cap(s2))
	//容量是从第一个下标开始算，数到原切片最后一个元素
	//子切片可以利用原切片预留的位置，进行append操作

	s3 := s1[5:]
	fmt.Printf("s3: %v,len=%d,cap=%d\n", s3, len(s3), cap(s3))

	s4 := s1[:4]
	fmt.Printf("s4: %v,len=%d,cap=%d\n", s4, len(s4), cap(s4))

	//子切片和切片究竟会不会互相影响，就抓住一点：它们是不是还共享数组，是否发生了扩容
	//只要子切片或切片任何一个发生了扩容，那它们就不是共享数组，就不会共享内存
}

func ShareSlice() {
	s1 := []int{1, 2, 3, 4}
	s2 := s1[2:]
	fmt.Printf("s1: %v,len=%d,cap=%d\n", s1, len(s1), cap(s1))
	fmt.Printf("s2: %v,len=%d,cap=%d\n", s2, len(s2), cap(s2))

	s2[0] = 99
	fmt.Printf("s1: %v,len=%d,cap=%d\n", s1, len(s1), cap(s1))
	fmt.Printf("s2: %v,len=%d,cap=%d\n", s2, len(s2), cap(s2))

	s2 = append(s2, 199)
	fmt.Printf("s1: %v,len=%d,cap=%d\n", s1, len(s1), cap(s1))
	fmt.Printf("s2: %v,len=%d,cap=%d\n", s2, len(s2), cap(s2))

	s2[1] = 1999
	fmt.Printf("s1: %v,len=%d,cap=%d\n", s1, len(s1), cap(s1))
	fmt.Printf("s2: %v,len=%d,cap=%d\n", s2, len(s2), cap(s2))

	//	只要有一个没有发生扩容，切片和子切片就是共享内存的
}
