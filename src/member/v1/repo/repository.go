package repo

import (
	"context"

	"github.com/pianzm/arr/src/member/v1/model"
)

type MemberRepository interface {
	Insert(ctx context.Context, payload []*model.Member) error
}
