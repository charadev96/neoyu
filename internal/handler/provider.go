package handler

import (
	"context"
	"errors"

	v1 "github.com/charadev96/neoyu/gen/neoyu/connection/v1"
	"github.com/charadev96/neoyu/internal/service"

	"connectrpc.com/connect"
	"github.com/google/uuid"
)

type Provider struct {
	svc *service.Provider
}

func NewProvider(s *service.Provider) *Provider {
	return &Provider{s}
}

func (h *Provider) ListProviders(ctx context.Context, req *v1.ListProvidersRequest) (*v1.ListProvidersResponse, error) {
	ps, err := h.svc.List()
	if err != nil {
		return nil, err
	}

	refs := make([]*v1.Provider, len(ps))
	for i, p := range ps {
		refs[i] = &p
	}

	return &v1.ListProvidersResponse{Providers: refs}, nil
}

func (h *Provider) GetProvider(ctx context.Context, req *v1.GetProviderRequest) (*v1.GetProviderResponse, error) {
	p, err := h.svc.Get(uuid.MustParse(req.Id))
	if err != nil {
		if errors.Is(err, service.ErrNotExist) {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		return nil, err
	}

	return &v1.GetProviderResponse{Provider: &p}, nil
}

func (h *Provider) SetProvider(ctx context.Context, req *v1.SetProviderRequest) (*v1.SetProviderResponse, error) {
	if err := h.svc.Set(*req.Provider); err != nil {
		if errors.Is(err, service.ErrNotExist) {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		return nil, err
	}

	return &v1.SetProviderResponse{}, nil
}

func (h *Provider) DeleteProvider(ctx context.Context, req *v1.DeleteProviderRequest) (*v1.DeleteProviderResponse, error) {
	if err := h.svc.Delete(uuid.MustParse(req.Id)); err != nil {
		if errors.Is(err, service.ErrNotExist) {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		return nil, err
	}

	return &v1.DeleteProviderResponse{}, nil
}
