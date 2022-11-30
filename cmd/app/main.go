package main

import (
	"blackwallgroup/config"
	"blackwallgroup/database"
	"blackwallgroup/internal/api"
	"blackwallgroup/internal/repository"
	"blackwallgroup/internal/service"
	"blackwallgroup/queue"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalln("error parsing config")
	}

	pool, err := database.Connect(ctx, cfg)
	if err != nil {
		log.Fatalln("error connecting database")
	}
	defer pool.Close()

	transactionQueue := queue.NewQueue()

	userRepository := repository.NewStorage(pool)
	userService := service.NewUserService(userRepository)

	srv := api.NewServer(pool, transactionQueue, userService)
	go func() {
		if err = srv.Start(cfg.HttpPort); err != nil {
			log.Fatalln(err)
		}
	}()

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	<-ch

	cancel()
	log.Println("shutting down app")
}
