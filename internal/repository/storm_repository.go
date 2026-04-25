package repository

import (
	"context"
	"fmt"
	"github.com/asdine/storm/v3"
	"github.com/eldius/mineserver-manager/internal/model"
)

type stormRepository struct {
	db *storm.DB
}

func NewStormRepository(dbPath string) (Repository, error) {
	db, err := storm.Open(dbPath)
	if err != nil {
		return nil, fmt.Errorf("opening storm db: %w", err)
	}
	return &stormRepository{db: db}, nil
}

func (r *stormRepository) SaveInstance(ctx context.Context, i *model.Instance) error {
	if err := r.db.Save(i); err != nil {
		return fmt.Errorf("saving instance: %w", err)
	}
	return nil
}

func (r *stormRepository) GetInstance(ctx context.Context, id string) (*model.Instance, error) {
	var i model.Instance
	if err := r.db.One("ID", id, &i); err != nil {
		return nil, fmt.Errorf("getting instance: %w", err)
	}
	return &i, nil
}

func (r *stormRepository) ListInstances(ctx context.Context) ([]model.Instance, error) {
	var instances []model.Instance
	if err := r.db.All(&instances); err != nil {
		return nil, fmt.Errorf("listing instances: %w", err)
	}
	return instances, nil
}

func (r *stormRepository) DeleteInstance(ctx context.Context, id string) error {
	var i model.Instance
	if err := r.db.One("ID", id, &i); err != nil {
		return fmt.Errorf("finding instance to delete: %w", err)
	}
	if err := r.db.DeleteStruct(&i); err != nil {
		return fmt.Errorf("deleting instance: %w", err)
	}
	return nil
}

func (r *stormRepository) Close() error {
	return r.db.Close()
}
