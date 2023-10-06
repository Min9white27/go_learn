package main

func Recursive(n int) {
	//使用递归的条件是要有中断的机制，否则会出现栈溢出
	//治标：栈溢出可以考虑通过增加栈的大小进行解决
	//治本：通过修改递归代码进行彻底解决
	//Recursive()
	if n > 10 {
		return
	}
	Recursive(n + 1)
}
