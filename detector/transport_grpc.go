package detector

import "golang.org/x/net/context"

func NewGrpcBinding(s Service) AnomalyDetectorServer{
	return &grpcBinding{
		s:s,
	}
}

type grpcBinding struct {
	s Service
}

func (s grpcBinding) AnalyseData(ctx context.Context, i *InputData) (*AnomalyScore, error) {
	score, msg, err := s.s.AddValueNow(ctx,i.Key,i.Value)
	return &AnomalyScore{
		Explanation: msg,
		Score: score,
	}, err
}

