package main

func main() {
	//str, name, age := Func9()
	//println(str, name, age)
	//
	//Func4("Deng", "Ming")
	//name, err := Func10()
	//println(name, err)

	//在多个返回值的时候，如果需要忽略某个返回值，可以使用 _
	//name1, _ := Func10()
	//println(name1)

	//使用:=的前提是，等号左边至少要有一个新变量
	//name1, _ := Func10()
	//println(name1)
	//	核心原则：同一个作用域内，如果左边出现了新变量，那么就需要用:=来接受返回值
	//Recursive(10)
	//UseFunctional1()

	//fn := Closure("阿橙")
	// fn 其实已经从 Closure 中返回了
	//但是我 fn 还要用到 "阿橙"
	//println(fn())
	//GetAge := Closure1()
	//println(GetAge())
	//GetAge := Closure2()
	//println(GetAge())
	//println(GetAge())
	//println(GetAge())
	//println(GetAge())
	//println(GetAge())
	//GetAge = Closure2()
	//println(GetAge())
	//println(GetAge())
	//println(GetAge())
	//println(GetAge())
	//println(GetAge())

	//Defer()
	//DeferClosure()
	//DeferClosureV1()

	//println(DeferReturn())
	//println(DeferReturnV1())

	DeferClosureLoopV1()
	DeferClosureLoopV2()
	DeferClosureLoopV3()
}
