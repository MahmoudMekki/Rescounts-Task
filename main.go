package main

import (
	"context"
	auth "github.com/MahmoudMekki/Rescounts-Task/cmd/auth-service/handler"
	prod "github.com/MahmoudMekki/Rescounts-Task/cmd/product-service/handler"
	user "github.com/MahmoudMekki/Rescounts-Task/cmd/user-service/handler"
	"github.com/MahmoudMekki/Rescounts-Task/config"
	"github.com/MahmoudMekki/Rescounts-Task/pkg/repo"
	"github.com/MahmoudMekki/Rescounts-Task/pkg/stripe"
	"github.com/MahmoudMekki/Rescounts-Task/pkg/token/jwt"
	"github.com/gorilla/mux"
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

	// initiating the internal services
	tknService := jwt.New(cfg.JWT.Secret)
	userRepo := repo.NewUserAccountRepo(db)
	prodRepo := repo.NewProductsRepo(db)
	stripeRepo := repo.NewStripeRepo(db)
	stripeClient := stripe.NewStripe(cfg.JWT.Secret)
	mw := auth.NewMiddleWare(l, tknService)

	//initiating handlers
	userSignupHandler := auth.NewCreateUserAccountHandler(l, userRepo, tknService)
	adminSignupHandler := auth.NewCreateAdminAccountHandler(l, userRepo, tknService)
	loginHandler := auth.NewLoginHandler(l, userRepo, tknService)
	getProdsHandler := prod.NewGetProductHandler(l, prodRepo, stripeClient, stripeRepo, tknService)
	addProdHadnler := prod.NewAddProductHandler(l, prodRepo, stripeClient, stripeRepo, tknService)
	updateProdHandler := prod.NewUpdateProductHandler(l, prodRepo, stripeClient, stripeRepo, tknService)
	deleteProdHandler := prod.NewDeleteProductHandler(l, prodRepo, stripeClient, stripeRepo, tknService)
	addCardHandler := user.NewAddCardHandler(l, stripeClient, stripeRepo, tknService, userRepo)
	purchaseHandler := user.NewAPurchaseHandler(l, stripeClient, stripeRepo, tknService, userRepo, prodRepo)

	s := mux.NewRouter()
	s.Handle("/auth/user/signup", userSignupHandler).Methods(http.MethodPost)
	s.Handle("/auth/admin/signup", adminSignupHandler).Methods(http.MethodPost)
	s.Handle("/auth/login", loginHandler).Methods(http.MethodGet)
	s.Handle("/products", getProdsHandler).Methods(http.MethodGet)
	s.Handle("/products", mw.MW(addProdHadnler)).Methods(http.MethodPost)
	s.Handle("/products/{prod_id}", mw.MW(updateProdHandler)).Methods(http.MethodPut)
	s.Handle("/products/{prod_id}", mw.MW(deleteProdHandler)).Methods(http.MethodDelete)
	s.Handle("/user/purchase", mw.MW(addCardHandler)).Methods(http.MethodPost)
	s.Handle("user/purchase/{prod_id}", mw.MW(purchaseHandler)).Methods(http.MethodGet)

	httpServer := &http.Server{
		Addr:         cfg.Http.Address,
		ReadTimeout:  time.Duration(cfg.Http.ReadTimeOutInSec) * time.Second,
		WriteTimeout: time.Duration(cfg.Http.WriteTimeOutInSec) * time.Second,
		IdleTimeout:  time.Duration(cfg.Http.IdleTimeOutInSec) * time.Second,
		Handler:      s,
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
