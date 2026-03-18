package worker

import (
	"context"
	"log"
	"net/http"
	"time"

	"site-monitor/internal/domain"
	"site-monitor/internal/repository"
)

type Checker struct {
	repo      repository.Repository
	broadcast chan<- *domain.Status
	interval  time.Duration
	client    *http.Client
}

func NewChecker(repo repository.Repository, broadcast chan<- *domain.Status, interval time.Duration) *Checker {
	return &Checker{
		repo:      repo,
		broadcast: broadcast,
		interval:  interval,
		client:    &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *Checker) Start(ctx context.Context) {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Println("Проверщик остановлен")
			return
		case <-ticker.C:
			c.checkAll(ctx)
		}
	}
}

func (c *Checker) checkAll(ctx context.Context) {
	sites, err := c.repo.GetSites(ctx)
	if err != nil {
		log.Printf("Не удалось получить список сайтов: %v", err)
		return
	}
	for _, site := range sites {
		select {
		case <-ctx.Done():
			return
		default:
		}
		status := c.checkSite(site)
		if err := c.repo.UpdateStatus(ctx, status); err != nil {
			log.Printf("Ошибка обновления статуса для %s: %v", site.ID, err)
		}
		select {
		case c.broadcast <- status:
		default:
			log.Println("Канал широковещания переполнен")
		}
	}
}

func (c *Checker) checkSite(site *domain.Site) *domain.Status {
	start := time.Now()
	resp, err := c.client.Get(site.URL)
	status := &domain.Status{
		SiteID:    site.ID,
		CheckedAt: start,
	}
	if err != nil {
		status.IsUp = false
		status.Error = err.Error()
		return status
	}
	defer resp.Body.Close()
	status.IsUp = resp.StatusCode >= 200 && resp.StatusCode < 300
	status.StatusCode = resp.StatusCode
	return status
}