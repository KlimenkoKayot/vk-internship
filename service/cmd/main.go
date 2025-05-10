package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/klimenkokayot/vk-internship/libs/logger"
	"github.com/klimenkokayot/vk-internship/service/config"
	"github.com/klimenkokayot/vk-internship/service/internal/app"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	log, err := logger.NewAdapter(&logger.Config{
		Adapter: logger.AdapterZap,
		Level:   logger.LevelDebug,
	})
	if err != nil {
		panic(err)
	}

	server := app.NewServer(log)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := server.Start(cfg.GRPC.Address); err != nil {
			log.Error("failed to start gRPC server", logger.Field{
				Key:   "error",
				Value: err,
			})
			cancel()
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		log.Info("shutting down server...")
	case <-ctx.Done():
		log.Info("server context done")
	}

	if err := server.Stop(ctx); err != nil {
		log.Error("failed to stop server", logger.Field{
			Key:   "error",
			Value: err,
		})
	}

	log.Info("server stopped")
}
