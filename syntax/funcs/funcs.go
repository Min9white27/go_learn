package main

// 同一个包里面，方法名字不能重名，GO没有重载特性
//main方法没有参数

// Func1 没有参数
func Func1() {

}

// Func2 有一个参数
func Func2(a int) {

}

// Func3 有多个参数
func Func3(a int, b float32) {

}

// Func4 多个参数，一个类型
func Func4(a, c string) {

}

//可以没有返回值，也可以一个，或者多个

// Func5 有一个返回值
func Func5(a, c string) string {
	//	有返回值，要保证一定返回
	return "Hello go"
}

// Func6 有多个返回值
func Func6(a, c string) (string, string, int) {
	//	有返回值，要保证一定返回
	return "Hello go", "DaMing", 1345
}

func Func7(a, c string) (str string, name string, age int) {
	//返回值可以带名字,不过要做到同意，有一个带名字就必须全都要

	return "Hello go", "DaMing", 1345
}

func Func8(a, c string) (str string, name string, age int) {
	//返回值可以带名字,不过要做到同意，有一个带名字就必须全都要
	str = "Hello go"
	name = "DaMing"
	age = 18
	return
}

func Func9() (str string, name string, age int) {
	//带名字的返回值可以直接返回
	//也可以不进行赋值，其等价于"" , "" , 0,就是对应类型的零值
	return
}

func Func10() (name string, err error) {
	//虽然带名字，但可以不用赋值
	return "DaMing", nil
}

//带名字的返回值的优点是可以明确的知道其返回值是个什么东西，但缺点是会无形中放大了返回值的作用域，使其从方法代码块最开始就生效
