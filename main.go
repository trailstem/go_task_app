package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/trailstem/go_task_app/config"
)

// func run(ctx context.Context, l net.Listener) error {
func run(ctx context.Context) error {

	cfg, err := config.New()

	if err != nil {
		return err
	}
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen port %d: %v", cfg.Port, err)
	}
	url := fmt.Sprintf("http://%s", l.Addr().String())
	log.Printf("start with: %v", url)

	mux := NewMux()
	s := NewServer(l, mux)
	return s.Run(ctx)
}

func main() {

	if err := run(context.Background()); err != nil {
		log.Printf("failed to listen terminate server : %v", err)
		os.Exit(1)
	}
}
