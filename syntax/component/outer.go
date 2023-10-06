package main

type Inner struct {
}

func (i Inner) DoSomething() {
	println("这是 Inner")
}

func (i Inner) SayHello() {
	println("Hello", i.Name())
}

func (i Inner) Name() string {
	return "Inner"
}

type IIIInner struct {
}

// Outer 一般都会选择用这种
type Outer struct {
	Inner //，嵌套，也就是组合
}

func (o Outer) Name() string {
	return "Outer"
}

type OuterV1 struct {
	Inner
}

func (i OuterV1) DoSomething() {
	println("这是 OuterV1")
}

type OOOOuter struct {
	Inner
	IIIInner
}

type OuterPtr struct {
	*Inner
}

func UseInner() {
	var o Outer
	o.DoSomething()
	o.Inner.DoSomething() //这两个是等价的

	var p *OuterPtr
	p.DoSomething()
}

//组合可以是以下ji几种情况：
//接口组合接口
//结构体组合结构体
//结构体组合结构体指针
//结构体组合接口
//可以组合多个

func main() {
	var o1 OuterV1
	o1.DoSomething()
	o1.Inner.DoSomething()
	//先找自己有没有实现，然后在找组合里面的实现

	var o Outer
	o.SayHello()
}

//组合不是继承，没有多态
