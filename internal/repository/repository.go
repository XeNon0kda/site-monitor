package repository

import (
	"context"
	"site-monitor/internal/domain"
)

type Repository interface {
	AddSite(ctx context.Context, site *domain.Site) error
	RemoveSite(ctx context.Context, id string) error
	GetSites(ctx context.Context) ([]*domain.Site, error)
	GetSite(ctx context.Context, id string) (*domain.Site, error)
	UpdateStatus(ctx context.Context, status *domain.Status) error
	GetStatus(ctx context.Context, siteID string) (*domain.Status, error)
}