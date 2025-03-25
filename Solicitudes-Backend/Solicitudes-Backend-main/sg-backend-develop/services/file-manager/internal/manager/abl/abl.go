package abl

import "context"

type Client interface {
	ValidateABLData(ctx context.Context, ablNumber string, ablType int) (bool, error)
}
