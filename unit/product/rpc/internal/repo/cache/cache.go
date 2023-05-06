package cache

import (
	"unit/product/rpc/internal/config"
)

type Repository struct {
}

func NewRepository(cfg config.Config) *Repository {
	return &Repository{}
}
