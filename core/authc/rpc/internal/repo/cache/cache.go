package cache

import (
	"core/authc/rpc/internal/config"
)

type Repository struct {
}

func NewRepository(cfg config.Config) *Repository {
	return &Repository{}
}
