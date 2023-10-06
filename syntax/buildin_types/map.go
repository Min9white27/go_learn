package main

func Map() {
	m1 := map[string]int{
		"key1": 123,
		"key2": 456,
		"key3": 789,
	}
	m1["key3"] = 987
	m1["key4"] = 120

	//make 记得预估容量
	m2 := make(map[string]int, 4) //12是容量，并不是其中的元素
	m2["const1"] = 10086

	val, ok := m1["大明"]
	if ok {
		//有这个键值对
		println(val)
	}
	val = m1["大明"]
	println("小明对应的值是", val) //这里输出的是其对应的零值，也就是0

	println(len(m2))
	for k, v := range m2 {
		println(k, v)
	}

	for k := range m2 {
		println(k, m2[k])
	}

	//map遍历都是随机的

	delete(m2, "const1") //删除
}
