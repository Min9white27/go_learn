package main

func Switch(status int) {
	//GO里面switch的优点是不用在每个段落都写一个break
	switch status {
	case 0:
		println("初始化")
	case 1:
		println("调试中")
	case 2:
		println("运行中")
	case 3:
		println("终止")
	default:
		println("未知状态")
		//default有没有都可以
	}
}

func SwitchBool(age int) {
	//可以通过判断语句建立分支
	switch {
	case age >= 18:
		println("成年人")
	case age >= 40:
		println("中年人")
	case age >= 70:
		println("老年人")
	case age < 18:
		println("小孩子")
	}
	//这种写法很少见，这里面的条件要尽量做到是互斥的，提升代码可读性
}
