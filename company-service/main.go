package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"github.com/vielendanke/grpc-rest-project/company-service/internal/taskservice"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/vielendanke/grpc-rest-project/company-service/configs"
	"golang.org/x/sync/errgroup"
)

//go:embed configs.json
var fs embed.FS

const configName = "configs.json"

func main() {
	data, readErr := fs.ReadFile(configName)
	if readErr != nil {
		log.Fatal(readErr)
	}
	cfg := configs.NewConfig()
	if unmErr := json.Unmarshal(data, cfg); unmErr != nil {
		log.Fatal(unmErr)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	errCh := make(chan error, 1)
	gr, errCtx := errgroup.WithContext(ctx)
	gr.Go(func() error {
		return task.StartServerGRPS(errCtx, cfg)
	})
	gr.Go(func() error {
		return task.StartServerHTTP(errCtx, cfg)
	})
	gr.Go(func() error {
		return task.StartMetricsServer(cfg)
	})
	go func(ctx context.Context, errCh chan error) {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		errCh <- fmt.Errorf("%v", <-sigCh)
	}(ctx, errCh)
	errCh <- gr.Wait()
	log.Printf("\nService terminated: %v", <-errCh)
}
