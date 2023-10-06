package main

//接口是一组行为的抽象
//即便是业务开发，也应该面向接口编程
//当你怀疑要不要定义接口的时候，加上去！！！

type List interface {
	Add(idx int, val any) error //可以有返回值
	Append(val any)
	Delete(index int)
}
