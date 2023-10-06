package main

import "fmt"

func Byte() {
	//其对应的时它的ASC LL码值
	//如果想要输出a,需要用到格式化
	var a byte = 'a'
	println(a)
	println(fmt.Sprintf("%c", a))

	//[]byte和string可以互相转换
	var str string = "Hello,Go!"
	var bs []byte = []byte(str) //这里面已经发生了赋值
	bs[0] = 'h'
	println(str, bs)
	//byte占一个字节，即8位
}
