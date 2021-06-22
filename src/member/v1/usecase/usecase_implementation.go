package usecase

import (
	"context"
	"encoding/json"

	redisConf "github.com/pianzm/arr/config/redis"
	"github.com/pianzm/arr/helper"
	"github.com/pianzm/arr/src/member/v1/model"
)

type UsecaseImpl struct {
	Publisher redisConf.Client
}

func NewUsecase(redisClient redisConf.Client) MemberUsecase {
	return &UsecaseImpl{
		Publisher: redisClient,
	}
}
func (a *UsecaseImpl) Publish(ctx context.Context, params *model.QueueStatus) error {
	if err := a.Publisher.Publish(ctx, helper.DownloadChannel, params); err != nil {
		return err
	}
	if _, err := a.Publisher.Set(ctx, params.RequestID, params, 0); err != nil {
		return err
	}

	return nil
}

func (a *UsecaseImpl) GetStatus(ctx context.Context, requestID string) (*model.QueueStatus, error) {
	res, err := a.Publisher.Get(ctx, requestID)
	if err != nil {
		return nil, err
	}
	model := model.QueueStatus{}
	if err := json.Unmarshal([]byte(res), &model); err != nil {
		return nil, err
	}
	return &model, nil
}

func (a *UsecaseImpl) SetStatus(ctx context.Context, params *model.QueueStatus) error {
	params.Completed = true
	if _, err := a.Publisher.Set(ctx, params.RequestID, params, 0); err != nil {
		return err
	}
	return nil
}
