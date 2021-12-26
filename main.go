package main

import (
	"context"
	auth "github.com/MahmoudMekki/Rescounts-Task/cmd/auth-service/handler"
	admin "github.com/MahmoudMekki/Rescounts-Task/cmd/product-service/handler"
	user "github.com/MahmoudMekki/Rescounts-Task/cmd/user-service/handler"
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
var l *log.Logger

func init() {
	cfg = config.LoadConfig()
	l = log.New(os.Stdout, "[Rescounts-Task] ", log.LstdFlags)
}
func main() {

	db := cfg.DataBase.OpenDB()
	tknService := jwt.New(cfg.JWT.Secret)
	userRepo := repo.NewUserAccountRepo(db)
	prodRepo := repo.NewProductsRepo(db)
	stripeRepo := repo.NewStripeRepo(db)
	stripeClient := stripe.NewStripe(cfg.JWT.Secret)
	mw := auth.NewMiddleWare(l, tknService)

	userSignupHandler := auth.NewCreateUserAccountHandler(l, userRepo, tknService)
	adminSignupHandler := auth.NewCreateAdminAccountHandler(l, userRepo, tknService)
	loginHandler := auth.NewLoginHandler(l, userRepo, tknService)
	GetAddProductHandler := admin.NewGetAddProductHandler(l, prodRepo, stripeClient, stripeRepo, tknService)
	deleteUpdateProductHandler := admin.NewDeleteUpdateProductHandler(l, prodRepo, stripeClient, stripeRepo, tknService)
	addCardHandler := user.NewAddCardHandler(l, stripeClient, stripeRepo, tknService, userRepo)
	puchaseHandler := user.NewAPurchaseHandler(l, stripeClient, stripeRepo, tknService, userRepo, prodRepo)

	mux := http.NewServeMux()
	mux.Handle("/auth/user/signup", userSignupHandler)
	mux.Handle("/auth/admin/signup", adminSignupHandler)
	mux.Handle("/auth/login", loginHandler)
	mux.Handle("/products", mw.MW(GetAddProductHandler))
	mux.Handle("/products/", mw.MW(deleteUpdateProductHandler))
	mux.Handle("/user/purchase", mw.MW(addCardHandler))
	mux.Handle("/user/purchase/", mw.MW(puchaseHandler))

	httpServer := &http.Server{
		Addr:         cfg.Http.Address,
		ReadTimeout:  time.Duration(cfg.Http.ReadTimeOutInSec) * time.Second,
		WriteTimeout: time.Duration(cfg.Http.WriteTimeOutInSec) * time.Second,
		IdleTimeout:  time.Duration(cfg.Http.IdleTimeOutInSec) * time.Second,
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
