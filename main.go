package main

import (
	"context"
	auth "github.com/MahmoudMekki/Rescounts-Task/cmd/auth-service/handler"
	admin "github.com/MahmoudMekki/Rescounts-Task/cmd/product-service/handler"
	"github.com/MahmoudMekki/Rescounts-Task/config"
	"github.com/MahmoudMekki/Rescounts-Task/pkg/repo"
	"github.com/MahmoudMekki/Rescounts-Task/pkg/stripe"
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
	prodRepo := repo.NewProductsRepo(db)
	stripeRepo := repo.NewStripeRepo(db)
	stripeClient := stripe.NewStripe(cfg.JWT.Secret)

	userSignupHandler := auth.NewCreateUserAccountHandler(l, userRepo, tknService)
	adminSignupHandler := auth.NewCreateAdminAccountHandler(l, userRepo, tknService)
	loginHandler := auth.NewLoginHandler(l, userRepo, tknService)
	GetAddProductHandler := admin.NewGetAddProductHandler(l, prodRepo, stripeClient, stripeRepo)
	deleteUpdateProductHandler := admin.NewDeleteUpdateProductHandler(l, prodRepo, stripeClient, stripeRepo)
	mw := auth.NewMiddleWare(l, tknService)

	mux := http.NewServeMux()
	mux.Handle("/auth/user/signup", userSignupHandler)
	mux.Handle("/auth/admin/signup", adminSignupHandler)
	mux.Handle("/auth/login", loginHandler)
	mux.Handle("/products", mw.MW(GetAddProductHandler))
	mux.Handle("/products/", mw.MW(deleteUpdateProductHandler))

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
