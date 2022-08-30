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

	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	l := log.New(os.Stdout, "", log.LUTC)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("test")); err != nil {
			l.Fatal(err)
		}
	})
	srv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	connURL := fmt.Sprintf("postgres://%s:%s@postgresdb:5432/%s",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	pgpool, err := pgxpool.Connect(context.TODO(), connURL)
	if err != nil {
		l.Fatal(err)
	}

	defer pgpool.Close()

	var greeting string
	err = pgpool.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		l.Fatal(err)
	}

	fmt.Println(greeting)

	go func() {
		l.Println("Listening on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	l.Println("shutting down the server, waiting for in-flight connections to finish")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		l.Fatal(err)
	}
	l.Println("server shut down")
}
