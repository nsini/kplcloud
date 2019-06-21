package namespace

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/nsini/kplcloud/src/repository"
	"time"
)

type loggingService struct {
	logger log.Logger
	Service
}

func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) Detail(ctx context.Context, name string) (resp *repository.Namespace, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "detail",
			"name", name,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Detail(ctx, name)
}
