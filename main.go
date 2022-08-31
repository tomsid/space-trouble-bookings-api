package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"space-trouble-bookings-api/api"
	"space-trouble-bookings-api/config"
	"space-trouble-bookings-api/db"
	"space-trouble-bookings-api/spacex"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	zapLog, err := zap.NewProduction()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer zapLog.Sync()
	l := zapLog.Sugar()

	cfg := config.Config{}
	if err = env.Parse(&cfg); err != nil {
		l.Fatal(err)
	}

	connURL := fmt.Sprintf("postgres://%s:%s@postgresdb:5432/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBName)

	pgpool, err := pgxpool.Connect(context.TODO(), connURL)
	if err != nil {
		l.Fatal(err)
	}
	defer pgpool.Close()

	m, err := migrate.New("file://db_migrations", connURL)
	if err != nil {
		l.Fatal(err)
	}
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			l.Info("DB version is the latest")
		} else {
			l.Fatal(err)
		}
	}

	spacexClient := spacex.NewClient(&http.Client{Timeout: 15 * time.Second})
	handlers := api.NewAPI(spacexClient, db.NewPGStorage(pgpool), l)
	r := chi.NewRouter()
	r.Get("/booking", handlers.Bookings)
	r.Post("/booking", handlers.BookFlight)

	srv := http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		l.Info("Listening on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	l.Info("shutting down the server, waiting for in-flight connections to finish")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		l.Fatal(err)
	}
	l.Info("server shut down")
}
