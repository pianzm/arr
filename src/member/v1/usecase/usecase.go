package usecase

import (
	"context"

	"github.com/pianzm/arr/src/member/v1/model"
)

type MemberUsecase interface {
	Publish(ctx context.Context, params *model.QueueStatus) error
	GetStatus(ctx context.Context, requestID string) (*model.QueueStatus, error)
	SetStatus(ctx context.Context, params *model.QueueStatus) error
}
