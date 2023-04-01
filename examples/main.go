package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	srv := http.Server{Addr: ":8080"}

	http.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello!")
	})

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		srv.ListenAndServe()
		fmt.Println("Shutting down...bye!")
	}()

	<-shutdown

	err := srv.Shutdown(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	wg.Wait()
}
