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

	"github.com/prometheus/client_golang/prometheus/promhttp"

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

	if err := database.RunMigrations(cfg.DB.MigrateDSN()); err != nil {
		log.Fatalf("Ошибка миграций: %v", err)
	}

	redisClient, err := database.NewRedis(&cfg.Redis)
	if err != nil {
		log.Fatalf("Ошибка подключения к Redis: %v", err)
	}
	defer redisClient.Close()

	loginLimiter := middleware.NewRateLimiter(
		cfg.LoginMaxAttempts,
		time.Duration(cfg.LoginWindowMin)*time.Minute,
		time.Duration(cfg.LoginBlockMin)*time.Minute,
	)
	handler.SetLoginRateLimiter(loginLimiter)
	loginRateLimitMiddleware := middleware.LoginRateLimitMiddleware(loginLimiter)

	userRepo := repository.NewUserRepository(db)
	linkRepo := repository.NewLinkRepository(db)

	authService := service.NewAuthService(userRepo, cfg.JWT.Secret, cfg.JWT.TTL)
	linkService := service.NewLinkService(linkRepo, redisClient)

	authHandler := handler.NewAuthHandler(authService)
	linkHandler := handler.NewLinkHandler(linkService)
	redirectHandler := handler.NewRedirectHandler(linkService)

	authMiddleware := middleware.AuthMiddleware([]byte(cfg.JWT.Secret))

	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	mux.HandleFunc("GET /r/{code}", redirectHandler.RedirectToLongURL)

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			http.ServeFile(w, r, "web/pages/index.html")
		case "/login":
			http.ServeFile(w, r, "web/pages/auth.html")
		case "/cabinet":
			http.ServeFile(w, r, "web/pages/cabinet.html")
		default:
			http.NotFound(w, r)
		}
	})

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	mux.Handle("GET /metrics", promhttp.Handler())

	mux.HandleFunc("POST /api/auth/register", authHandler.Register)
	mux.Handle("POST /api/auth/login", loginRateLimitMiddleware(http.HandlerFunc(authHandler.Login)))
	mux.Handle("POST /api/auth/logout", authMiddleware(http.HandlerFunc(authHandler.Logout)))
	mux.Handle("GET /api/auth/me", authMiddleware(http.HandlerFunc(authHandler.Me)))

	mux.HandleFunc("POST /api/links", middleware.OptionalAuthMiddleware([]byte(cfg.JWT.Secret))(linkHandler.CreateShortLink))
	mux.Handle("GET /api/links", authMiddleware(http.HandlerFunc(linkHandler.GetLinks)))
	mux.Handle("DELETE /api/links/{id}", authMiddleware(http.HandlerFunc(linkHandler.DeleteLink)))
	mux.Handle("GET /api/links/{id}/stats", authMiddleware(http.HandlerFunc(linkHandler.GetStats)))

	handlerWithMetrics := middleware.MetricsMiddleware(mux)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      handlerWithMetrics,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		users, _ := userRepo.Count(context.Background())
		middleware.SetActiveUsers(int64(users))

		links, _ := linkRepo.Count(context.Background())
		middleware.SetActiveLinks(int64(links))

		ticker := time.NewTicker(5 * time.Minute)
		go func() {
			for range ticker.C {
				users, _ := userRepo.Count(context.Background())
				middleware.SetActiveUsers(int64(users))
				links, _ := linkRepo.Count(context.Background())
				middleware.SetActiveLinks(int64(links))
			}
		}()
	}()

	stop := make(chan os.Signal, 2)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Сервер запущен на http://localhost%s", srv.Addr)
		log.Printf("Prometheus метрики: http://localhost%s/metrics", srv.Addr)
		log.Printf("Grafana: http://localhost:3000 (admin/admin)")
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
