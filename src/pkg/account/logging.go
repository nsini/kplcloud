package account

import (
	"context"
	"github.com/go-kit/kit/log"
	"time"
)

type loggingService struct {
	logger log.Logger
	Service
}

func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) Detail(ctx context.Context, id int64) (rs map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "detail",
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Detail(ctx, id)
}

func (s *loggingService) List(ctx context.Context, order, by string, pageSize, offset int) (rs []map[string]interface{}, count uint64, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "List",
			"order", order,
			"by", by,
			"pageSize", pageSize,
			"offset", offset,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.List(ctx, order, by, pageSize, offset)
}

func (s *loggingService) Popular(ctx context.Context) (rs []map[string]interface{}, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "Popular",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Popular(ctx)
}
