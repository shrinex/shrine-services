package cache

import "unit/shop/rpc/internal/config"

type Repository struct {
}

func NewRepository(cfg config.Config) *Repository {
	return &Repository{}
}
