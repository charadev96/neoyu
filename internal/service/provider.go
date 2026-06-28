package service

import (
	v1 "github.com/charadev96/neoyu/gen/neoyu/connection/v1"
	"github.com/charadev96/neoyu/internal/db"

	"github.com/google/uuid"
)

type ProviderSchema struct {
	Providers *[]*v1.Provider `json:"providers"`
}

type Provider struct {
	db *db.Lens[ProviderSchema, *v1.Provider]
}

func NewProvider(s *db.File[ProviderSchema]) *Provider {
	return &Provider{
		db.NewLens[ProviderSchema, *v1.Provider](s,
			func(s ProviderSchema) *[]*v1.Provider {
				return s.Providers
			},
		),
	}
}

func (s *Provider) List() ([]*v1.Provider, error) {
	ps, err := s.db.List()
	return ps, err
}

func (s *Provider) Get(id uuid.UUID) (*v1.Provider, error) {
	p, err := s.db.Load(id)
	return p, err
}

func (s *Provider) Set(p *v1.Provider) error {
	return s.db.Save(p)
}

func (s *Provider) Delete(id uuid.UUID) error {
	return s.db.Delete(id)
}
