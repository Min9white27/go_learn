package lock

import "sync"

// LockDemo
// 优先使用 RWMutex，优先加读锁
// 读锁就是你加了读锁，别人同时也可以加读锁，但是不能同时加写锁，只有你释放了读锁后才可以。别人还可以继续读，但是不能写
// 写锁就是你加了写锁，谁都不能加读写锁。别人不能读写，只有你能。
// 常用地并发优化手段，用读写锁来优化读锁
type LockDemo struct {
	lock sync.Mutex
}

func NewLockDemo() *LockDemo {
	return &LockDemo{}
}

func (l *LockDemo) PanicDemo() {
	l.lock.Lock()
	//	如果在中间 panic 了，就无法释放锁
	panic("abc")
	l.lock.Unlock()
}

// DeferDemo 一般情况上，这样写最合理
func (l *LockDemo) DeferDemo() {
	l.lock.Lock()
	defer l.lock.Unlock()
}

// NoPointerDemo LockDemo不加指针，在它调起NoPointerDemo这个方法的时候，会引起复制，锁也会复制，这样你就会有了两把锁，编辑器就会报错。
//func (l LockDemo) NoPointerDemo() {
//	l.lock.Lock()
//	defer l.lock.Unlock()
//}

// LockDemoV1 新人推荐用这种
type LockDemoV1 struct {
	lock *sync.Mutex
}

// NewLockDemoV1 在 sync.Mutex 是指针的时候，就需要显示的初始化它，反之则不用
func NewLockDemoV1() LockDemoV1 {
	return LockDemoV1{
		// 如果不初始化，lock 就是 nil，你一用就会 panic
		lock: &sync.Mutex{},
	}
}

func (l LockDemoV1) NoPointerDemo() {
	l.lock.Lock()
	defer l.lock.Unlock()
}
