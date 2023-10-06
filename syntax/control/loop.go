package main

//原则：千万不要对迭代参数取地址！！！因为取出来的地址都会是一样的

func ForLoop() {
	for i := 0; i < 10; i++ {
		println(i)
	}

	for i := 0; i < 10; {
		println(i)
		i++
	}

	i := 0
	for ; i < 10; i++ {
		println(i)
	}
}

func loop1() {
	i := 0
	for i < 10 {
		i++
		println(i)
	}

	//死循环容易出CPU100%的问题
	for true {
		i++
		println(i)
	}

	for {
		i++
		println(i)
	}
}

func ForArr() {
	arr := [3]int{1, 2, 3}
	for index, val := range arr {
		println("下标 ", index, "值 ", val)
	}

	for index := range arr {
		println("下标 ", index, "值 ", arr[index])
	} //结果一样
}

func ForSlice() {
	slice := []int{1, 2, 3}
	for index, val := range slice {
		println("下标 ", index, "值 ", val)
	}

	for index := range slice {
		println("下标 ", index, "值 ", slice[index])
	} //结果一样
}

func ForMap() {
	m := map[string]int{
		"key1": 100,
		"key2": 200,
		"key3": 300,
	}

	for k, v := range m {
		println(k, v)
	}

	for k := range m {
		println(k, m[k])
	}
	//for对map的遍历是完全随机的，输出结果顺序不定
}

func LoopBreak() {
	i := 0
	for true {
		if i > 10 {
			println(i)
			break
		}
		i++
	}
}

func LoopContinue() {
	i := 0
	for i < 10 {
		if i%2 == 1 {
			//println("奇数")
			i++
			continue
		}
		println(i)
		i++
	}
}

func LoopContinueV1() {
	i := 0
	for i < 10 {
		i++
		if i%2 == 1 {
			continue
		}
		println(i)
	}
}
