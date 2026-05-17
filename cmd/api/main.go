package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mikaelcaua/welcome-university-api/internal/infra/config"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/database"
	"github.com/mikaelcaua/welcome-university-api/internal/infra/bootstrap"
)

func main() {
	ctx := context.Background()
	appConfig := config.Load()
	if err := appConfig.Validate(); err != nil {
		log.Fatalf("invalid configuration: %v", err)
	}

	pool, err := database.Connect(ctx, appConfig)
	if err != nil {
		log.Fatalf("falha ao conectar no banco: %v", err)
	}
	defer pool.Close()

	if err := database.Migrate(ctx, pool); err != nil {
		log.Fatalf("falha ao migrar banco: %v", err)
	}

	router, shutdown, err := bootstrap.Create(ctx, appConfig, pool)
	if err != nil {
		log.Fatalf("falha ao criar rotas: %v", err)
	}
	defer shutdown()

	server := &http.Server{
		Addr:              ":" + appConfig.ServerPort,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("api ouvindo na porta %s", appConfig.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("falha ao iniciar servidor: %v", err)
		}
	}()

	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, syscall.SIGINT, syscall.SIGTERM)
	<-stopSignal

	shutdownContext, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownContext); err != nil {
		log.Printf("falha ao encerrar servidor: %v", err)
	}
}
