package detector

import (
	"time"
	"github.com/go-kit/kit/log"
	"golang.org/x/net/context"
)

func NewLoggingService(l log.Logger, s Service) Service {
	return &logger{
		l,
		s,
	}
}

type logger struct {
	logger log.Logger
	Service
}

func (s logger) AddValueNow(ctx context.Context, key string, value float64) (annomaly float64, reason string, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "AddValueNow",
			"key", key,
			"value", value,
			"score", annomaly,
			"reason", reason,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.AddValueNow(ctx, key, value)
}

func (s logger) AddValue(ctx context.Context, key string, value float64, t time.Time) (annomaly float64, reason string, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "AddValue",
			"key", key,
			"value", value,
			"score", annomaly,
			"reason", reason,
			"time", t,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.AddValue(ctx, key, value, t)
}
