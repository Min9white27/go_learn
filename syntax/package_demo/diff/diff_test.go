package diff_test

import "gitee.com/geekbang/basic-go/syntax/package_demo/diff"

func UseHello() {
	//要调用Hello()，它会认为它在两个不同的包里面，所以需要在前面加包的前缀
	//正常情况不会用这种做法，只会在集成测试里面使用
	diff.Hello()
}
