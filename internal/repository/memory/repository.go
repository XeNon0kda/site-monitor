package memory

import (
	"context"
	"sync"

	"github.com/google/uuid"
	
	"site-monitor/internal/domain"
)

type Repository struct {
	mu       sync.RWMutex
	sites    map[string]*domain.Site
	statuses map[string]*domain.Status
}

func New() *Repository {
	return &Repository{
		sites:    make(map[string]*domain.Site),
		statuses: make(map[string]*domain.Status),
	}
}

func (r *Repository) AddSite(_ context.Context, site *domain.Site) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if site.ID == "" {
		site.ID = uuid.New().String()
	}
	r.sites[site.ID] = site
	return nil
}

func (r *Repository) RemoveSite(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.sites, id)
	delete(r.statuses, id)
	return nil
}

func (r *Repository) GetSites(_ context.Context) ([]*domain.Site, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	sites := make([]*domain.Site, 0, len(r.sites))
	for _, s := range r.sites {
		sites = append(sites, s)
	}
	return sites, nil
}

func (r *Repository) GetSite(_ context.Context, id string) (*domain.Site, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.sites[id], nil
}

func (r *Repository) UpdateStatus(_ context.Context, status *domain.Status) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.statuses[status.SiteID] = status
	return nil
}

func (r *Repository) GetStatus(_ context.Context, siteID string) (*domain.Status, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.statuses[siteID], nil
}