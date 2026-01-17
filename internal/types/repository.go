package types

import (
	"context"

	. "github.com/crazydw4rf/oil-bank-backend/internal/pkg/result"
)

type BaseRepository[T any] interface {
	Create(ctx context.Context, entity *T) Result[*T]
	Find(ctx context.Context, id int64) Result[*T]
	Update(ctx context.Context, entity *T) Result[*T]
	Delete(ctx context.Context, id int64) Result[bool]
}
