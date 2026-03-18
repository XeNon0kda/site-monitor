package service

import (
	"context"
	"log"
	"time"

	"site-monitor/internal/domain"
	"site-monitor/internal/repository"
	"site-monitor/internal/worker"
)

type Monitor struct {
	repo        repository.Repository
	checker     *worker.Checker
	subscribers map[chan *domain.Status]struct{}
	subscribe   chan chan *domain.Status
	unsubscribe chan chan *domain.Status
	broadcast   chan *domain.Status
}

func NewMonitor(repo repository.Repository, interval time.Duration) *Monitor {
	m := &Monitor{
		repo:        repo,
		subscribers: make(map[chan *domain.Status]struct{}),
		subscribe:   make(chan chan *domain.Status),
		unsubscribe: make(chan chan *domain.Status),
		broadcast:   make(chan *domain.Status, 100),
	}
	m.checker = worker.NewChecker(repo, m.broadcast, interval)
	return m
}

func (m *Monitor) Start(ctx context.Context) {
	go m.runHub()
	go m.checker.Start(ctx)
}

func (m *Monitor) runHub() {
	for {
		select {
		case ch := <-m.subscribe:
			m.subscribers[ch] = struct{}{}
		case ch := <-m.unsubscribe:
			delete(m.subscribers, ch)
			close(ch)
		case status := <-m.broadcast:
			for ch := range m.subscribers {
				select {
				case ch <- status:
				default:
					log.Println("Канал подписчика переполнен, пропускаем")
				}
			}
		}
	}
}

func (m *Monitor) Subscribe() chan *domain.Status {
	ch := make(chan *domain.Status, 10)
	m.subscribe <- ch
	return ch
}

func (m *Monitor) Unsubscribe(ch chan *domain.Status) {
	m.unsubscribe <- ch
}

func (m *Monitor) AddSite(ctx context.Context, url string) (*domain.Site, error) {
	site := &domain.Site{
		URL:       url,
		CreatedAt: time.Now(),
	}
	if err := m.repo.AddSite(ctx, site); err != nil {
		return nil, err
	}
	return site, nil
}

func (m *Monitor) RemoveSite(ctx context.Context, id string) error {
	return m.repo.RemoveSite(ctx, id)
}

type SiteWithStatus struct {
	Site   *domain.Site   `json:"site"`
	Status *domain.Status `json:"status"`
}

func (m *Monitor) GetSites(ctx context.Context) ([]SiteWithStatus, error) {
	sites, err := m.repo.GetSites(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]SiteWithStatus, 0, len(sites))
	for _, s := range sites {
		status, _ := m.repo.GetStatus(ctx, s.ID)
		result = append(result, SiteWithStatus{Site: s, Status: status})
	}
	return result, nil
}