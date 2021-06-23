package query

import (
	"context"
)

type MemberRead interface {
	GetAll(ctx context.Context) ([][]string, error)
	GetTotal(ctx context.Context) (int64, error)
}
