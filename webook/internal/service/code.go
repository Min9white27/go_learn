package service

import (
	"context"
	"fmt"
	"gitee.com/geekbang/basic-go/webook/internal/repository"
	"gitee.com/geekbang/basic-go/webook/internal/service/sms"
	"math/rand"
)

const codeTPlId = "1877556"

var (
	ErrCodeSendTooMany        = repository.ErrCodeSendTooMany
	ErrCodeVerifyTooManyTimes = repository.ErrCodeVerifyTooManyTimes
)

type CodeService interface {
	Send(ctx context.Context, biz, phone string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}

type CodeRedisService struct {
	repo   repository.CodeRepository
	smsSvc sms.Service
	//tplId string
}

func NewCodeRedisService(repo repository.CodeRepository, smsSvc sms.Service) CodeService {
	return &CodeRedisService{
		repo:   repo,
		smsSvc: smsSvc,
	}
}

func (svc *CodeRedisService) Send(ctx context.Context,
	// biz 用于区别业务场景
	biz string, phone string) error {
	//	生成一个验证码
	code := svc.generateCode()
	//	塞进去 Redis
	err := svc.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	//	发送出去
	err = svc.smsSvc.Send(ctx, codeTPlId, []string{code}, phone)
	if err != nil {
		//  这意味着，Redis 有这个验证码
		//  这个 err 可能是超时的 err，所以这个 err 的发送状态并不能确定，有可能已经发送出去了
		//	要重试的话，初始化的时候，传入一个自己就会重试的 smsSvc
		//	但是引入重试的话，用户可能会出现收到多次验证码的问题，还有可能会出现错误扩散的问题，导致系统负载过高
		//	实际情况上，发送失败都不会是系统去处理，不会进行重试，而是让用户自己去重新申请验证码，不过这需要一段时间
	}
	return err

}

// Verify bool 代表验证码通没通过， error 代表系统有没有出错
func (svc *CodeRedisService) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	return svc.repo.Verify(ctx, biz, phone, inputCode)
}

func (svc *CodeRedisService) generateCode() string {
	// 六位数，num 在 0， 999999 之间，包括 0 和 999999
	num := rand.Intn(1000000)
	// 不够六位数的，加上前导 0
	// 比如生成随机数1，就会格式化成 000001
	return fmt.Sprintf("%06d", num)
}

// VerifyV1 非标准写法，可以不用管验证码通过的问题，过于没过都当成系统异常与否的问题
//func (svc *CodeRedisService) VerifyV1(ctx context.Context, biz string, phone string, inputCode string) error {
//
//}
