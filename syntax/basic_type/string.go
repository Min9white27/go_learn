package main

import (
	"fmt"
	"unicode/utf8"
)

func String() {
	//在字符串里面的”“（双引号）需要用\来进行转义
	//不需要手写转义，可以事先写好，然后再复制到Goland,IDE会自动完成转义
	println("He said:\"Hello , GO !\"")
	println(`我可以换行
这里是新的
但是这里不能再有反引号
同时也不能用\来进行转义`)
	//拼接,GO不可以进行别的类型拼接
	println("Hello " + "Go")
	//println("Hello" + string(123)) 这个123会转换成其对应的ASC LL码
	//如果一定要拼接123，要用到格式化
	println(fmt.Sprintf("Hello %d", 123))
	println(len("abc"))
	println(len("你好"))
	println(utf8.RuneCountInString("你好"))
	//字节长度：和编码无关，用len(str)获取
	//字符长度：和编码有关，用编码库来计算，默认情况下使用utf8库

	//byte bytes包
}
