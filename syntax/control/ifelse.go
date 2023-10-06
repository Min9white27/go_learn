package main

func IfOnly(age int) string {
	if age >= 18 {
		println("成年人")
	}
	return ""
}

func IfElse(age int) {
	if age >= 18 {
		println("成年人")
	} else {
		println("未成年")
	}
}

func IfElseIf(age int) {
	if age >= 18 {
		println("成年人")
	} else if age >= 12 {
		println("青年人")
	} else {
		println("未成年")
	}
}

func IfNewVariable(start int, end int) string {
	if distance := end - start; distance > 100 {
		println("太远了")
	} else if distance > 60 {
		println("有点远")
	} else {
		println("还好")
	}
	return ""
}
