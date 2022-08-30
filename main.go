package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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
