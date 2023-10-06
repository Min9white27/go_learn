package main

//defer 的作用主要是用来释放资源

// Defer 也叫延迟调用
func Defer() {
	//defer 是后定义的先出，先定后出
	defer func() {
		println("第一个 defer")
	}()

	defer func() {
		println("第二个 defer")
	}()
}

func DeferClosure() {
	i := 0
	defer func() {
		println(i)
	}()
	i = 1
	//输出的是 1
}

func DeferClosureV1() {
	i := 0
	defer func(val int) {
		println(val)
	}(i)
	i = 1
	//输出的是 0
	//这里 i 的值 0 ，事先已经传了过去，已经固定了，所以后面再更改也没用了
}

//确定值的原则：1.作为参数传入时，定义defer的时候就确定了；2.作为闭包引入的，执行defer对应方法的时候才确定

//defer 修改返回值
//返回值带名字的可以修改，没有名字的修改不了

func DeferReturn() int {
	a := 0
	defer func() {
		a = 1
	}()
	return a
}

func DeferReturnV1() (a int) {
	a = 0
	defer func() {
		a = 1
	}()
	return a
}

//指针的修改是可以生效的

func DeferReturnV2() *MyStruct {
	a := &MyStruct{
		"阿橙",
	}
	defer func() {
		a.name = "橙子"
	}()
	return a
}

type MyStruct struct {
	name string
}

// 自测题
func DeferClosureLoopV1() {
	for i := 0; i < 10; i++ {
		defer func() {
			println(i)
		}()
		//10个10
	}
}

func DeferClosureLoopV2() {
	for i := 0; i < 10; i++ {
		defer func(val int) {
			println(val)
			//9-0
		}(i)

	}
}

func DeferClosureLoopV3() {
	for i := 0; i < 10; i++ {
		j := i
		defer func() {
			println(j)
			//9-0
		}()
	}
}
