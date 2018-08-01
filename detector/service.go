package detector

import (
	"time"
	"gonum.org/v1/gonum/stat"
	"log"
	"math"
	"fmt"
	"golang.org/x/net/context"
)

type Service interface {
	AddValueNow(ctx context.Context, key string, value float64) (annomaly float64, reason string, err error)
	AddValue(ctx context.Context, key string, value float64, t time.Time) (annomaly float64, reason string, err error)
}

func NewService(store Store) Service {
	return &anomaly{
		store: store,
	}

}

type anomaly struct {
	store Store
}

func (s anomaly) AddValueNow(ctx context.Context, key string, value float64) (annomaly float64, reason string, err error) {
	return s.AddValue(ctx, key, value, time.Now())
}

func (s anomaly) AddValue(ctx context.Context, key string, value float64, t time.Time) (annomaly float64, reason string, err error) {
	h, err := s.store.addHourData(key, t.Hour(), t.Minute(), value)
	if err != nil {
		return 0, "", err
	}

	d, err := s.store.addDayOfWeekData(key, t.Hour(), t.Minute(), int(t.Weekday()), value)
	if err != nil {
		return 0, "", err
	}

	m, err := s.store.addDayOfMonthData(key, t.Hour(), t.Minute(), t.Day(), value)
	if err != nil {
		return 0, "", err
	}

	ha, hm := s.calculateScore(key, "hourly", h, value)
	da, dm := s.calculateScore(key, "day of week", d, value)
	ma, mm := s.calculateScore(key, "day of month", m, value)

	annomaly = (ha + da + ma) / 3

	reason = fmt.Sprintf("The current value is %v. Using previous data for hour %v, minute %v. The normal value for this time is %v. Normally at this time on a %v the expected value would be %v. Normally at this time on the %v day of the month we would expect the value to be %v",
		value, t.Hour(), t.Minute(), hm, t.Weekday(), dm, t.Day(), mm,
	)

	return

}

func (s anomaly) calculateScore(key, recordtype string, input []float64, value float64) (annomaly, mean float64) {
	mean, stddev := stat.MeanStdDev(input, nil)
	stderr := stat.StdErr(stddev, float64(len(input)))

	meandelta := math.Abs(value - mean)
	if meandelta > stddev {
		annomaly = meandelta / stddev
	}

	log.Printf("%v:%v - value %v, mean %v, standard deviation %v, error: %v, annomaly %v", key, recordtype, value, mean, stddev, stderr, annomaly)
	return

}
