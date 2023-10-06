package main

import "fmt"

func NewUser() {
	//初始化结构体
	u := User{}
	fmt.Printf("%v\n", u)
	fmt.Printf("%+v\n", u) //最好用这个

	//up是一个指针
	up := &User{}
	fmt.Printf("up %+v\n", up) //

	up2 := new(User)
	fmt.Printf("up2 %+v\n", up2)
	//	上面两个效果一样

	u4 := User{Name: "Tom", FirstName: "Lee", Age: 18}
	//永远要选择用这种初始化，因为在生产环境中，结构体是很复杂的，用下面的初始化，很容易会搞混，编译器也不会报错
	fmt.Printf("%+v\n", u4)
	u5 := User{"Tom", "Lee", 18} //用这种，同类型的字段容易搞混
	fmt.Printf("%+v\n", u5)

	//赋值后是可以改的
	u4.Name = "Jerry"
	u5.Age = 20

	var up3 *User
	println(up3) //如果声明了一个指针，没有赋值，编译器会默认为nil，其地址为0x0
	//	在空指针上面想要访问它的值或者是方法会报错

}

type User struct {
	Name      string
	FirstName string
	Age       uint8
}

//声明成结构体的时候，想要改变字段的值，会发现改不了

func (u User) ChangeName(name string) {
	fmt.Printf("change name 中 u 的地址 %p\n", &u)
	u.Name = name
	//方法调用本身是值传递的
}

//声明成地址的时候，才可以改变字段里面的值

func (u *User) ChangeAge(age uint8) {
	fmt.Printf("change name 中 u 的地址 %p\n", u)
	u.Age = age
	//这里传的是地址，即复制的也是地址，所以改的值是原来地址的值
	//遇事不决用指针，提高指针优先级
	//指针可以调用结构体方法，反过来也可以
}

func ChangeUser() {
	u1 := User{Name: "Tom", Age: 30}
	fmt.Printf("u1 的地址 %p\n", &u1)
	//(&u1).ChangeName("Jerry")
	u1.ChangeName("Jerry") //改的是其复制的值
	u1.ChangeAge(35)
	fmt.Printf("%+v", u1)

	up1 := &User{Name: "Tom", Age: 30}
	fmt.Printf("u1 的地址 %p\n", &up1)
	//(*up1).ChangeName("Jerry")
	up1.ChangeName("Jerry") //也不会生效d
	up1.ChangeAge(35)
	fmt.Printf("%+v", up1)
}

type Integer int

func UseInt() {
	i1 := 10
	i2 := Integer(i1)
	var i3 Integer = 11
	println(i2, i3)
}

type Fish struct {
	Name *Fish
}

func (f Fish) Swim() {
	println("fish 在游")
}

type FakeFish Fish

func UseFish() {
	var f1 Fish = Fish{}
	f1.Swim()
	f2 := FakeFish{&f1} //F2.Swim()
	//println(f1.Name)
	println(f2.Name)
	//搞不懂
}
