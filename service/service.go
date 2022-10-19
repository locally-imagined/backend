package helpers

import (
	"context"

	"github.com/locally-imagined/goa/gen/calc"
)

type Service struct{}

func (s *Service) Multiply(ctx context.Context, p *calc.MultiplyPayload) (int, error) {
	return *p.A + *p.B, nil
}

func (s *Service) Add(ctx context.Context, p *calc.AddPayload) (int, error) {
	return *p.A * *p.B, nil
}
