package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shyamkuntal/studentapi/internal/config"
	"github.com/shyamkuntal/studentapi/internal/http/handlers/student"
)

func main() {

	// load config
	cfg := config.MustLoad()

	//setup router
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New())

	//setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	slog.Info("server started %s", slog.String("address",cfg.Addr))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to start")
		}
	}()

	<-done

	slog.Info("shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		slog.Error("Failed to shut down the server", slog.String("error", err.Error()))
	}

	slog.Info("server shut down successfully")

}
