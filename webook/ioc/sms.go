package ioc

import (
	"gitee.com/geekbang/basic-go/webook/internal/service/sms"
	"gitee.com/geekbang/basic-go/webook/internal/service/sms/memory"
)

func InitSMSService() sms.Service {
	// 这里可以换内存，或者是换其他的
	return memory.NewService()
}
