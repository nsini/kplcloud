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

func (s *loggingService) Get(ctx context.Context, name string) (resp *repository.Namespace, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "get",
			"name", name,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Get(ctx, name)
}

func (s *loggingService) Post(ctx context.Context, name, displayName string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "post",
			"name", name,
			"displayName", displayName,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Post(ctx, name, displayName)
}

func (s *loggingService) Update(ctx context.Context, name, displayName string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "update",
			"name", name,
			"displayName", displayName,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Update(ctx, name, displayName)
}

func (s *loggingService) Sync(ctx context.Context) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "sync",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.Sync(ctx)
}
