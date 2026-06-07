package service

import (
	"fmt"
	"slices"

	v1 "github.com/charadev96/neoyu/gen/neoyu/connection/v1"
	"github.com/charadev96/neoyu/internal/db"

	"github.com/google/uuid"
)

type ProviderSchema struct {
	Providers []v1.Provider `json:"providers"`
}

type Provider struct {
	db *db.FileStore[ProviderSchema]
}

func NewProvider(s *db.FileStore[ProviderSchema]) *Provider {
	return &Provider{s}
}

func (s *Provider) List() ([]v1.Provider, error) {
	data, err := s.db.Load()
	if err != nil {
		return nil, fmt.Errorf("load: %w", err)
	}

	return data.Providers, nil
}

func (s *Provider) Get(id uuid.UUID) (v1.Provider, error) {
	data, err := s.db.Load()
	if err != nil {
		return v1.Provider{}, fmt.Errorf("load: %w", err)
	}

	i := s.index(data.Providers, id)
	if i == -1 {
		return v1.Provider{}, ErrNotExist
	}

	return data.Providers[i], nil
}

func (s *Provider) Set(id uuid.UUID, p v1.Provider) error {
	data, err := s.db.Load()
	if err != nil {
		return fmt.Errorf("load: %w", err)
	}

	i := s.index(data.Providers, id)
	if i == -1 {
		data.Providers = append(data.Providers, p)
	} else {
		data.Providers[i] = p
	}

	if err := s.db.Save(data); err != nil {
		return fmt.Errorf("save: %w", err)
	}
	return nil
}

func (s *Provider) Delete(id uuid.UUID) error {
	data, err := s.db.Load()
	if err != nil {
		return fmt.Errorf("load: %w", err)
	}

	i := s.index(data.Providers, id)
	if i == -1 {
		return ErrNotExist
	}
	data.Providers = slices.Delete(data.Providers, i, i)

	if err := s.db.Save(data); err != nil {
		return fmt.Errorf("save: %w", err)
	}
	return nil
}

func (s Provider) index(ps []v1.Provider, id uuid.UUID) int {
	return slices.IndexFunc(ps, func(p v1.Provider) bool {
		return p.Id == id.String()
	})
}
