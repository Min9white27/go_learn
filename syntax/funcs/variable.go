package main

func YourName(name string, alias ...string) {
	//alias是一个切片
	if len(alias) > 0 {
		println(alias[0])
	}
}

func YourNameInvoke() {
	YourName("阿橙")
	YourName("阿橙", "橙子")
	YourName("阿橙", "橙子", "清见橙")
}
