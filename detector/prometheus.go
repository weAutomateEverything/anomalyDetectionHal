package detector

import (
	"github.com/go-kit/kit/metrics"
	"context"
	"time"
)

func NewPrometheus(counter metrics.Counter, errorCount metrics.Counter, latency metrics.Histogram, s Service) Service {
	return &prom{
		requestCount:   counter,
		requestLatency: latency,
		errorCount:     errorCount,
		Service:        s,
	}
}

type prom struct {
	requestCount   metrics.Counter
	errorCount     metrics.Counter
	requestLatency metrics.Histogram
	Service
}

func (s *prom) AddValueNow(ctx context.Context, key string, value float64) (annomaly float64, reason string, err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "AddValueNow").Add(1)
		s.requestLatency.With("method", "AddValueNow").Observe(time.Since(begin).Seconds())
		if err != nil {
			s.errorCount.With("method", "AddValueNow").Add(1)
		}
	}(time.Now())
	return s.Service.AddValueNow(ctx,key,value)
}
func (s *prom) AddValue(ctx context.Context, key string, value float64, t time.Time) (annomaly float64, reason string, err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "AddValue").Add(1)
		s.requestLatency.With("method", "AddValue").Observe(time.Since(begin).Seconds())
		if err != nil {
			s.errorCount.With("method", "AddValue").Add(1)
		}
	}(time.Now())
	return s.Service.AddValue(ctx, key,value,t)
}
