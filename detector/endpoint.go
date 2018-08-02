package detector

import (
	"github.com/go-kit/kit/endpoint"
	"context"
)

func makeAddDataEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*addDataRequest)

		anomaly, reason, err := s.AddValueNow(ctx, req.key, req.value)
		if err != nil {
			return nil, err
		}
		return &AnomalyAddDataResponse{
			AnomalyScore: anomaly,
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
	Explination  string  `json:"explination"`
}
