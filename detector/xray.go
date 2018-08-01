package detector

import (
	"time"
	"github.com/aws/aws-xray-sdk-go/xray"
	"context"
)

func NewXray(s Service) Service {
	return &x{
		s,
	}
}

type x struct {
	service Service
}

func (x *x) AddValueNow(ctx context.Context,key string, value float64) (annomaly float64, reason string, err error) {
	xray.Capture(ctx, "detector.AddValueNow", func(ctx context.Context) error {
		annomaly, reason, err = x.service.AddValueNow(ctx, key, value)
		return err
	})
	return
}

func (x *x) AddValue(ctx context.Context,key string, value float64, t time.Time) (annomaly float64, reason string, err error) {
	xray.Capture(ctx, "detector.AddValue", func(ctx context.Context) error {
		annomaly, reason, err = x.service.AddValue(ctx, key, value, t)
		return err
	})
	return
}
