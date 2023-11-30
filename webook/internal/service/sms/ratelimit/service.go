package ratelimit

import (
	"context"
	"fmt"
	"gitee.com/geekbang/basic-go/webook/internal/service/sms"
	"gitee.com/geekbang/basic-go/webook/pkg/ratelimit"
)

var errLimited = fmt.Errorf("触发了限流")

type RatelimitSMSService struct {
	svc     sms.Service
	limiter ratelimit.Limiter
}

func NewRatelimitSMSService(svc sms.Service, limiter ratelimit.Limiter) sms.Service {
	return &RatelimitSMSService{
		svc:     svc,
		limiter: limiter,
	}
}

func (r *RatelimitSMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	limiter, err := r.limiter.Limit(ctx, "sms:tencent")
	if err != nil {
		// 系统错误
		// 可以限流，在下游很坑的时候，保守限流
		// 可以不限，在下游很强的时候，业务可用性要求很高的时候，尽量容错策略，
		return fmt.Errorf("短信服务判断是否限流出现问题, %w", err)
	}
	if limiter {
		return errLimited
	}
	// 加代码，加新特性
	err = r.svc.Send(ctx, tpl, args, numbers...)
	// 加代码，加新特性
	return err
}
