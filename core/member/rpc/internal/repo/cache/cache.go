package cache

import "core/member/rpc/internal/config"

type Repository struct {
}

func NewRepository(cfg config.Config) *Repository {
	return &Repository{}
}
