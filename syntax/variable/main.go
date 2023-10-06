package main

// 通过首字母控制包内包外，首字母大写包外可以访问
//var Global = "全局变量"
//
//// 首字母小写，只能是包内里面可以使用，其子包不能使用
//var internal = "包变量"

var (
	First  string = "abc"
	second uint32 = 123
	aa            = "hello"
)

func main() {
	//a,b,c为局部变量，私有变量
	//局部变量声明了必须要使用，否则会报错
	var a int = 123
	var b = 456 //省略类型，go可以进行类型推断，不需要写类型
	var c uint = 789
	println(a, b, c)
	e := 678 //只能用在局部变量，Golang会进行类型推断
	println(e)
	//同作用域下同个变量只能声明一次，但可以更改它的值
	//var a int = 456
	a = 456
	var aa int = 123
	println(aa)
}
