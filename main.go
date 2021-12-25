package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	l := log.New(os.Stdout, "[Rescounts-Task] ", log.LstdFlags)
	mux := http.NewServeMux()

	httpServer := &http.Server{
		Addr:         ":9090",
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}

	//Go routine for gracefully shutdown the server
	go func() {
		err := httpServer.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Listening to interrupt or kill signal from the OS
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Received Terminate, Gracefully shutdown the server", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	l.Print("goodbye \n")
	httpServer.Shutdown(tc)
}
