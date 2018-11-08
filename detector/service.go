package detector

import (
	"golang.org/x/net/context"
	"gonum.org/v1/gonum/stat"
	"log"
	"math"
	"time"
)

type Service interface {
	AddValueNow(ctx context.Context, key string, value float64) (annomaly, average, hour, day, month float64, err error)
	AddValue(ctx context.Context, key string, value float64, t time.Time) (annomaly, average, hour, day, month float64, err error)
}

func NewService(store Store) Service {
	return &anomaly{
		store: store,
	}

}

type anomaly struct {
	store Store
}

func (s anomaly) AddValueNow(ctx context.Context, key string, value float64) (annomaly, average, hour, day, month float64, err error) {
	return s.AddValue(ctx, key, value, time.Now())
}

func (s anomaly) AddValue(ctx context.Context, key string, value float64, t time.Time) (annomaly, average, hour, day, month float64, err error) {
	h, err := s.store.addHourData(key, t.Hour(), t.Minute(), value)
	if err != nil {
		return
	}

	d, err := s.store.addDayOfWeekData(key, t.Hour(), t.Minute(), int(t.Weekday()), value)
	if err != nil {
		return
	}

	m, err := s.store.addDayOfMonthData(key, t.Hour(), t.Minute(), t.Day(), value)
	if err != nil {
		return
	}

	ha, hour := s.calculateScore(key, "hourly", h, value)
	da, day := s.calculateScore(key, "day of week", d, value)
	ma, month := s.calculateScore(key, "day of month", m, value)
	annomaly = (ha + da + ma) / 3
	average = (hour + day + month) / 3

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
