package main

import (
	"backend-service-template/config"
	user_domain "backend-service-template/internal/domain/user"
	user_handler "backend-service-template/internal/handler/user"
	user_service "backend-service-template/internal/service/user"
	"backend-service-template/pkg/database"
	"backend-service-template/pkg/logger"
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	envPath := flag.String("envpath", ".env", "Path to .env file")

	flag.Parse()

	cfg := config.LoadConfig(*envPath)

	log := logger.New(cfg.Level)

	ctx := context.Background()
	db, err := database.NewPostgres(ctx, cfg)
	if err != nil {
		log.Fatal("database connection failed", zap.Error(err))
	}
	defer db.Close()

	userRepo := user_domain.NewRepository(db)
	userSvc := user_service.NewService(userRepo)
	userHandler := user_handler.NewHandler(userSvc, log)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /users", userHandler.Create)
	mux.HandleFunc("GET /users/{id}", userHandler.Get)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.HttpPort),
		Handler: mux,
	}

	go func() {
		log.Info("server starting", zap.String("addr", srv.Addr))
		err := srv.ListenAndServe()
		if err != nil {
			log.Fatal("server start failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = srv.Shutdown(ctx)
	if err != nil {
		log.Fatal("server forced to shutdown", zap.Error(err))
	}

	log.Info("server exited properly")
}
