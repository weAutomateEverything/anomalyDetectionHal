package detector

import (
	"fmt"
	"gopkg.in/mgo.v2"
)

const hourKey = "%v_%00d_%00d"
const weekKey = "%v_%00d_%00d_dow_%d"
const monthKey = "%v_%00d_%00d_dom_%d"

func NewDataStore(db *mgo.Database) Store {
	return &store{
		db: db,
	}
}

type Store interface {
	getHourData(key string, hour, minute int) ([]float64, error)
	getDayOfWeekData(key string, hour, minute, dayOfWeek int) ([]float64, error)
	getDayOfMonthData(key string, hour, minute, dayOfMonth int) ([]float64, error)

	addHourData(key string, hour, minute int, value float64) ([]float64, error)
	addDayOfWeekData(key string, hour, minute, dayOfWeek int, value float64) ([]float64, error)
	addDayOfMonthData(key string, hour, minute, dayOfMonth int, value float64) ([]float64, error)
}

type store struct {
	db *mgo.Database
}

func (s *store) addHourData(key string, hour, minute int, value float64) ([]float64, error) {
	k := fmt.Sprintf(hourKey, key, hour, minute)
	return s.update(k, value)

}

func (s *store) addDayOfWeekData(key string, hour, minute, dayOfWeek int, value float64) ([]float64, error) {
	k := fmt.Sprintf(weekKey, key, hour, minute, dayOfWeek)
	return s.update(k, value)
}

func (s *store) addDayOfMonthData(key string, hour, minute, dayOfMonth int, value float64) ([]float64, error) {
	k := fmt.Sprintf(monthKey, key, hour, minute, dayOfMonth)
	return s.update(k, value)

}

func (s *store) getHourData(key string, hour, minute int) ([]float64, error) {
	k := fmt.Sprintf(hourKey, key, hour, minute)
	return s.get(k)
}

func (s *store) getDayOfWeekData(key string, hour, minute, dayOfWeek int) ([]float64, error) {
	k := fmt.Sprintf(weekKey, key, hour, minute, dayOfWeek)
	return s.get(k)

}

func (s *store) getDayOfMonthData(key string, hour, minute, dayOfMonth int) ([]float64, error) {
	k := fmt.Sprintf(monthKey, key, hour, minute, dayOfMonth)
	return s.get(k)

}

func (s *store) get(key string) ([]float64, error) {
	q, err := s.find(key)
	if err != nil {
		return nil, err
	}
	return q.Data, nil
}

func (s *store) find(key string) (hourData, error) {
	c := s.db.C("anomaly_data")
	q := c.FindId(key)
	h := &hourData{}

	count, err := q.Count()
	if err != nil {
		return *h, err
	}

	if count == 0 {
		h.Key = key
		h.Data = make([]float64, 0)
		return *h, nil
	}

	err = q.One(&h)
	return *h, err

}

func (s *store) update(key string, value float64) ([]float64, error) {
	q, err := s.find(key)
	if err != nil {
		return nil, err
	}

	q.Data = append(q.Data, value)
	c := s.db.C("anomaly_data")
	if len(q.Data) == 1{
		err = c.Insert(q)
	} else {
		err = c.UpdateId(key,q)
	}
	return q.Data, err
}

type hourData struct {
	Key  string `json:"id" bson:"_id,omitempty"`
	Data []float64
}
