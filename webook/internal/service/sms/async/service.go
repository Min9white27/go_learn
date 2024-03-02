package async

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/repository"
	"gitee.com/geekbang/basic-go/webook/internal/service/sms"
	"gitee.com/geekbang/basic-go/webook/pkg/logger"
	"time"
)

type Service struct {
	svc  sms.Service
	repo repository.AsyncSmsRepository
	l    logger.LoggerV1
	quit chan struct{}
}

func NewService(svc sms.Service, repo repository.AsyncSmsRepository, l logger.LoggerV1) *Service {
	res := &Service{
		svc:  svc,
		repo: repo,
		l:    l,
		quit: make(chan struct{}),
	}
	go func() {
		res.StartAsyncCycle()
	}()
	return res
}

func (s *Service) StartAsyncCycle() {
	for {
		select {
		case <-s.quit:
			return
		default:
			s.AsyncSend()
		}
	}
}

func (s *Service) AsyncSend() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	as, err := s.repo.PreemptWaitingSMS(ctx)
	cancel()
	switch err {
	case nil:
		// 执行发送
		ctx, cancel = context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		err = s.svc.Send(ctx, as.TplId, as.Args, as.Numbers...)
		if err != nil {
			s.l.Error("执行异步发送短信失败",
				logger.Error(err),
				logger.Int64("id", as.Id))
		}
		res := err == nil
		// 通知 repository 我这一次的执行结果
		err = s.repo.ReportScheduleResult(ctx, as.Id, res)
		if err != nil {
			s.l.Error("执行异步发送短信成功，但是标记数据库失败",
				logger.Error(err),
				logger.Bool("res", res),
				logger.Int64("id", as.Id))
		}
	case repository.ErrWaitingSMSNotFound:
		time.Sleep(time.Second)
	default:
		s.l.Error("抢占异步发送短信任务失败",
			logger.Error(err))
		time.Sleep(time.Second)
	}
}

func (s *Service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	if s.needAsync(ctx) {
		// 需要异步发送，直接转储到数据库
		err := s.repo.Add(ctx, domain.AsyncSms{
			TplId:   tplId,
			Args:    args,
			Numbers: numbers,
			// 设置可以重试三次
			RetryMax: 3,
		})
		return err
	}
	return s.svc.Send(ctx, tplId, args, numbers...)
}

// 1. 基于响应时间的，平均响应时间
//	// 1.1 使用绝对阈值，比如说直接发送的时候，（连续一段时间，或者连续N个请求）响应时间超过了 500ms，然后后续请求转异步

// 什么时候退出异步
// 1. 进入异步 N 分钟后
// 2. 保留 1% 的流量（或者更少），继续同步发送，判定响应时间
func (s *Service) needAsync(ctx context.Context) bool {
	responseTimeThreshold := time.Millisecond * 500
	averageResponseTime := time.Duration(s.repo.GetAverageResponseTime(ctx))
	if averageResponseTime > responseTimeThreshold {
		return true
	}
	return false
}

func (s *Service) quitAsync() {
	close(s.quit)
}
