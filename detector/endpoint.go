package detector

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

func makeAddDataEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*addDataRequest)

		anomaly, average, reason, err := s.AddValueNow(ctx, req.key, req.value)
		if err != nil {
			return nil, err
		}
		return &AnomalyAddDataResponse{
			AnomalyScore: anomaly,
			Average:      average,
			Explination:  reason,
		}, nil
	}
}

type addDataRequest struct {
	key   string
	value float64
}

type AnomalyAddDataResponse struct {
	AnomalyScore float64 `json:"anomaly_score"`
	Average      float64 `json:"average"`
	Explination  string  `json:"explination"`
}
