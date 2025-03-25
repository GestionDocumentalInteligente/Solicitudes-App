package manager

import (
	"context"

	"github.com/teamcubation/sg-file-manager-api/internal/manager/abl"
)

type ablUseCase struct {
	client abl.Client
}

type ABLUseCase interface {
	ValidateABLData(ctx context.Context, ablNumber string, ablType int) (bool, error)
}

func NewAblUseCase(client abl.Client) ABLUseCase {
	return &ablUseCase{
		client: client,
	}
}

func (a *ablUseCase) ValidateABLData(ctx context.Context, ablNumber string, ablType int) (bool, error) {
	return a.client.ValidateABLData(ctx, ablNumber, ablType)
}
