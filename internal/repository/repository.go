package repository

import (
	"context"
	"github.com/eldius/mineserver-manager/internal/model"
)

type Repository interface {
	SaveInstance(ctx context.Context, i *model.Instance) error
	GetInstance(ctx context.Context, id string) (*model.Instance, error)
	ListInstances(ctx context.Context) ([]model.Instance, error)
	DeleteInstance(ctx context.Context, id string) error
	Close() error
}
