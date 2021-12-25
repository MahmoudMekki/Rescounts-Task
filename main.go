package main

import (
	"context"
	"github.com/MahmoudMekki/Rescounts-Task/cmd/auth-service/handler"
	"github.com/MahmoudMekki/Rescounts-Task/config"
	"github.com/MahmoudMekki/Rescounts-Task/pkg/repo"
	"github.com/MahmoudMekki/Rescounts-Task/pkg/token/jwt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var cfg config.Config

func init() {
	cfg = config.LoadConfig()
}
func main() {
	l := log.New(os.Stdout, "[Rescounts-Task] ", log.LstdFlags)
	db := cfg.DataBase.OpenDB()
	tknService := jwt.New(cfg.JWT.Secret)
	userRepo := repo.NewUserAccountRepo(db)

	userSignupHandler := handler.NewCreateUserAccountHandler(l, userRepo, tknService)
	adminSignupHandler := handler.NewCreateAdminAccountHandler(l, userRepo, tknService)

	mux := http.NewServeMux()
	mux.Handle("/user/signup", userSignupHandler)
	mux.Handle("/admin/signup", adminSignupHandler)
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
	db.Close()
	httpServer.Shutdown(tc)
}
