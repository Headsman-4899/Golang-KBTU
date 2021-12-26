package main

import (
	"context"
	"lectures/hw6/internal/http"
	"lectures/hw6/internal/store/mongodb"
	"log"
	"os"
	"os/signal"
	"syscall"

	lru "github.com/hashicorp/golang-lru"
)

func main() {
	store := mongodb.Init()

	cache, err := lru.New2Q(6)
	if err != nil {
		panic(err)
	}

	srv := http.NewServer(context.Background(),
		http.WithAddress(":8080"),
		http.WithStore(store),
		http.WithCache(cache),
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
