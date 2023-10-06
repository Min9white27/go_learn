package main

// 可以把方法赋值给变量，然后变量就可以直接发起调用。例如：你希望用户传入的数据可以调用一些常用的函数

func Functional1() {
	println("hello functional")
}

var Abc = func() string {
	return "hello"
}

func Functional2(a int) {

}

func UseFunctional1() {
	MyFunc := Functional1
	MyFunc()
	//想要通过赋值给变量进行函数的重载，就不能改变原来的参数列表
	//Abc = func(a int) string {
	//
	//}
	//有参数的重载在调用的时候还是需要写入参数
	MyFunc2 := Functional2
	MyFunc2(18)
}

func Functional3() {
	//新定义了一个方法，赋值给了 fn，或者说是定义了一个方法名为 fn 的方法
	//这个方法称之为局部方法，作用域只在Functional3这个方法内
	//在实际中，这个做法是为了不让其他作用域用到这个方法
	fn := func() string {
		return "hello"
	}
	fn()
}

// Functional4 方法作为返回值，它的意思是返回一个，返回 string的无参数方法
func Functional4() func() string {
	return func() string {
		return "hello"
	}
}

// Functional5 匿名方法立刻发起调用
func Functional5() {
	fn := func() string {
		return "hello"
	}()
	println(fn)
}
