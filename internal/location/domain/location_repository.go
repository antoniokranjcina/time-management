package domain

import "context"

type LocationRepository interface {
	Create(ctx context.Context, location *Location) (*Location, error)
	GetAll(ctx context.Context) ([]Location, error)
	GetById(ctx context.Context, id string) (*Location, error)
	Update(ctx context.Context, id, name string) (*Location, error)
	Delete(ctx context.Context, id string) error
}
