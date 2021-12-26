package main

import (
	"context"
	"lectures/hw6/internal/http"
	"lectures/hw6/internal/message_broker"
	"lectures/hw6/internal/message_broker/kafka"
	"lectures/hw6/internal/store/mongodb"
	"log"
	"os"
	"os/signal"
	"syscall"

	lru "github.com/hashicorp/golang-lru"
)

func main() {
	store := mongodb.Init()
	ctx, cancel := context.WithCancel(context.Background())
	go CatchTermination(cancel)

	twoQueueCache, err := lru.New2Q(6)
	if err != nil {
		panic(err)
	}

	brokers := []string{"localhost:29092"}
	broker := kafka.NewBroker(brokers, twoQueueCache, "peer2")
	if err := broker.Connect(ctx); err != nil {
		panic(err)
	}
	defer func(broker message_broker.MessageBroker) {
		err := broker.Close()
		if err != nil {

		}
	}(broker)

	srv := http.NewServer(
		ctx,
		http.WithAddress(":8080"),
		http.WithStore(store),
		http.WithCache(twoQueueCache),
		http.WithBroker(broker),
	)
	if err := srv.Run(); err != nil {
		log.Println(err)
	}

	srv.WaitForGracefulTermination()
}

func CatchTermination(cancel context.CancelFunc) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Print("[WARN] caught termination signal")
	cancel()
}
