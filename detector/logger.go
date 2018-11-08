package detector

import (
	"github.com/go-kit/kit/log"
	"golang.org/x/net/context"
	"time"
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

func (s logger) AddValueNow(ctx context.Context, key string, value float64) (annomaly, average, hour, day, month float64, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "AddValueNow",
			"key", key,
			"value", value,
			"score", annomaly,
			"average", average,
			"hour", hour,
			"day", day,
			"month", month,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.AddValueNow(ctx, key, value)
}

func (s logger) AddValue(ctx context.Context, key string, value float64, t time.Time) (annomaly, average, hour, day, month float64, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "AddValue",
			"key", key,
			"value", value,
			"score", annomaly,
			"average", average,
			"hour", hour,
			"day", day,
			"month", month,
			"time", t,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.AddValue(ctx, key, value, t)
}
