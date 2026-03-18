package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"site-monitor/internal/config"
	"site-monitor/internal/handler"
	"site-monitor/internal/repository/memory"
	"site-monitor/internal/service"
)

func main() {
	cfg := config.Load()

	repo := memory.New()
	monitor := service.NewMonitor(repo, cfg.CheckInterval)
	h := handler.NewHandler(monitor)

	r := mux.NewRouter()
	h.RegisterRoutes(r)

	srv := &http.Server{
		Addr:         cfg.Port,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	monitor.Start(ctx)

	go func() {
		log.Printf("Сервер запущен на %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка сервера: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Завершение работы...")

	cancel() 

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()
	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatal("Принудительное завершение:", err)
	}
	log.Println("Сервер остановлен")
}