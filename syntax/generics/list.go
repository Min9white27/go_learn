package main

import "errors"

// T 类型参数 约束是 any 等于没有约束
//这种做法的好处是，可以在要用的时候，随意指定T的类型
//如果直接指定，T的类型就会固定，用法就大大缩减了
//这样可以写的非常灵活，也不用写过多重复代码

type List[T any] interface {
	Add(adx int, t T)
	Append(t T)
}

func main() {
	//UseList()
	//println(Sum[int](1, 2, 3))
	//println(Sum[float64](1.1, 2.2, 3.3))
	//println(Sum[Integer](1, 2, 3))
	println(Max[int](1, 2, 3, 9, 8, 7))
	println(Min[int](1, 2, 3, 9, 8, -1))
}

func UseList() {
	var l List[int]
	//l.Append("string") 传入字符就会报错，因为指定了是整型
	l.Append(100)
}

//接口和结构体都可以用泛型

type LinkedList[T any] struct {
	head *node[T]
	t    T
}

type node[T any] struct {
	val T
}

// Max 求最大值
func Max[T Number](vals ...T) (T, error) {
	if len(vals) == 0 {
		var t T
		return t, errors.New("你的下标不对")
	}
	res := vals[0]
	for i := 1; i < len(vals); i++ {
		if res < vals[i] {
			res = vals[i]
		}
	}
	return res, nil
}

func Min[T Number](vals ...T) (T, error) {
	if len(vals) == 0 {
		var t T
		return t, errors.New("你的下标不对")
	}
	res := vals[0]
	for i := 1; i < len(vals); i++ {
		if res > vals[i] {
			res = vals[i]
		}
	}
	return res, nil
}
