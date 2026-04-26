package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Wild-sergunys/shrtic/internal/config"
	"github.com/Wild-sergunys/shrtic/internal/database"
	"github.com/Wild-sergunys/shrtic/internal/handler"
	"github.com/Wild-sergunys/shrtic/internal/middleware"
	"github.com/Wild-sergunys/shrtic/internal/repository"
	"github.com/Wild-sergunys/shrtic/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфига: %v", err)
	}

	db, err := database.NewPostgres(&cfg.DB)
	if err != nil {
		log.Fatalf("Ошибка подключения к PostgreSQL: %v", err)
	}
	defer db.Close()

	// Миграции
	if err := database.RunMigrations(cfg.DB.MigrateDSN()); err != nil {
		log.Fatalf("Ошибка миграций: %v", err)
	}

	redisClient, err := database.NewRedis(&cfg.Redis)
	if err != nil {
		log.Fatalf("Ошибка подключения к Redis: %v", err)
	}
	defer redisClient.Close()

	// Репозитории
	userRepo := repository.NewUserRepository(db)
	linkRepo := repository.NewLinkRepository(db)

	// Сервисы
	authService := service.NewAuthService(userRepo, cfg.JWT.Secret, cfg.JWT.TTL)
	linkService := service.NewLinkService(linkRepo, redisClient)

	// Хендлеры
	authHandler := handler.NewAuthHandler(authService)
	linkHandler := handler.NewLinkHandler(linkService)
	redirectHandler := handler.NewRedirectHandler(linkService)

	// Middleware
	authMiddleware := middleware.AuthMiddleware([]byte(cfg.JWT.Secret))

	// Роутер
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Auth
	mux.HandleFunc("POST /api/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/auth/login", authHandler.Login)
	mux.Handle("POST /api/auth/logout", authMiddleware(http.HandlerFunc(authHandler.Logout)))
	mux.Handle("GET /api/auth/me", authMiddleware(http.HandlerFunc(authHandler.Me)))

	// Links
	mux.Handle("POST /api/links", authMiddleware(http.HandlerFunc(linkHandler.CreateShortLink)))
	mux.Handle("GET /api/links", authMiddleware(http.HandlerFunc(linkHandler.GetLinks)))
	mux.Handle("DELETE /api/links/{id}", authMiddleware(http.HandlerFunc(linkHandler.DeleteLink)))

	// Redirect (публичный)
	mux.HandleFunc("GET /{code}", redirectHandler.RedirectToLongURL)

	// Сервер
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	stop := make(chan os.Signal, 2)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Сервер запущен на http://localhost%s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	<-stop
	log.Println("Получен сигнал остановки")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("Ожидание завершения активных запросов...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Сервер завершён принудительно: %v", err)
	} else {
		log.Println("Все запросы завершены корректно")
	}
}
