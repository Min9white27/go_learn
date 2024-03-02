package repository

import (
	"context"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
	"github.com/ecodeclub/ekit/sqlx"
)

var ErrWaitingSMSNotFound = dao.ErrWaitingSMSNotFound

type AsyncSmsRepository interface {
	// Add 添加一个异步 SMS 记录。
	// 你叫做 Create 或者 Insert 也可以
	Add(ctx context.Context, s domain.AsyncSms) error
	PreemptWaitingSMS(ctx context.Context) (domain.AsyncSms, error)
	ReportScheduleResult(ctx context.Context, id int64, success bool) error
	GetAverageResponseTime(ctx context.Context) int64
}

type asyncSmsRepository struct {
	dao dao.AsyncSmsDAO
}

func NewAsyncSMSRepository(dao dao.AsyncSmsDAO) AsyncSmsRepository {
	return &asyncSmsRepository{
		dao: dao,
	}
}

func (a *asyncSmsRepository) Add(ctx context.Context, s domain.AsyncSms) error {
	return a.dao.Insert(ctx, dao.AsyncSms{
		Config: sqlx.JsonColumn[dao.SmsConfig]{
			Val: dao.SmsConfig{
				TplId:   s.TplId,
				Args:    s.Args,
				Numbers: s.Numbers,
			},
			Valid: true,
		},
		RetryMax: s.RetryMax,
	})
}

func (a *asyncSmsRepository) PreemptWaitingSMS(ctx context.Context) (domain.AsyncSms, error) {
	as, err := a.dao.GetWaitingSMS(ctx)
	if err != nil {
		return domain.AsyncSms{}, err
	}
	return domain.AsyncSms{
		Id:       as.Id,
		TplId:    as.Config.Val.TplId,
		Numbers:  as.Config.Val.Numbers,
		Args:     as.Config.Val.Args,
		RetryMax: as.RetryMax,
	}, nil
}

func (a *asyncSmsRepository) ReportScheduleResult(ctx context.Context, id int64, success bool) error {
	if success {
		return a.dao.MarkSuccess(ctx, id)
	}
	return a.dao.MarkFailed(ctx, id)
}

func (a *asyncSmsRepository) GetAverageResponseTime(ctx context.Context) int64 {
	// 1. 从数据库中查询所有异步短信记录的响应时间
	responseTimes, err := a.dao.GetAllResponseTimes(ctx)
	if err != nil {
		return 0
	}

	// 2. 计算平均响应时间
	var totalResponseTime int64
	for _, responseTime := range responseTimes {
		totalResponseTime += responseTime
	}
	averageResponseTime := totalResponseTime / int64(len(responseTimes))

	return averageResponseTime
}
