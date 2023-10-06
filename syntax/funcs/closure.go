package main

import "fmt"

// Closure 闭包，方法+它绑定的运行上下文（看它的定义是用在外面的）
func Closure(name string) func() string {
	return func() string {
		return "Hello, " + name
	}
}

func Closure1() func() string {
	name := "阿橙"
	age := 18
	return func() string {
		return fmt.Sprintf("Hello,%s,%d", name, age)
	}
}

// 闭包如果使用不当，可能会引起内存泄露的问题

func Closure2() func() int {
	age := 0
	return func() int {
		age++
		return age
	}
}
