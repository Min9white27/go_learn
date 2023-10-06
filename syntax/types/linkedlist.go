package main

//结构体和结构体的字段都遵循大小写控制访问性的原则

type LinkedList struct {
	head *node
	tail *node

	//包外可以访问
	Len int
}

func (l *LinkedList) Add(idx int, val any) error {
	//TODO implement me
	panic("implement me")
}

func (l LinkedList) Append(val any) {
	//TODO implement me
	panic("implement me")
}

func (l LinkedList) Delete(index int) {
	//TODO implement me
	panic("implement me")
}

//func (l LinkedList) Add(idx int, val any) {
//
//}
//
//// 方法接收器，receiver
//
//func (l *LinkedList) AddV1(idx int, val any) {
//
//}

//上面这两者不一样，第一个Add是定义在LinkedList上面，第二个AddV1是定义在*LinkedList上面

type node struct {
	prev *node
	next *node
	//如果在结构体内部还要用结构体自身，也就是要进行递归，那么就只能用指针
	//因为在编译的时候，编译器要确定这个结构体有多大，如果直接写node字段，这样编译器会算不出来，会报错
	//指针在特定的机器上大小是确定的，32位的是4个字节，64位的是8个字节
}
