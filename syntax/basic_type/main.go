package main

func main() {
	var a int = 123
	var b int = 456
	println(a + b)
	println(a - b)
	println(a * b)
	if b != 0 {
		//取余取模时最好都要判断除数是否为0，否则容易报错
		println(a / b)
		println(a % b)
	}

	//只能进行同类型运算，可以进行强制类型转换，不过要小心精度损失
	var c float64 = 78.9
	println(a + int(c)) //float64(a)

	//var d int16 = 789
	//println(a + d)
	//math.Ceil()
	//int uint都有最大最小值，float系列只有最大值和最小正值，没有最小值
	//math.MaxFloat64
	//math.SmallestNonzeroFloat64
	String()
	//Byte()
}
