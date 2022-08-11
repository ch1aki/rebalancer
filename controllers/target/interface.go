package provider

import (
	"context"
)

type GetSetWeighter interface {
	GetWeight(ctx context.Context) (int64, error)
	SetWeight(ctx context.Context, value int64) error
}
