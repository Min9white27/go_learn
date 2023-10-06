package main

const External = "包外"
const internal = "包内"
const (
	i = 123 //一样支持类型推断
)

// iota 枚举
// 方便地控制常量初始化
const (
	statusA = iota
	statusB
	statusC
	statusD

	day = iota << 1
)

const (
	a = iota + 11*11
	b

	c = 100
	d
)

func main() {
	const a = 123
	//a = 456 常量不能改变
}
