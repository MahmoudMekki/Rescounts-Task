package handler

import (
	"github.com/MahmoudMekki/Rescounts-Task/cmd/user-service/data"
	"github.com/MahmoudMekki/Rescounts-Task/kit/rhttp"
	"github.com/MahmoudMekki/Rescounts-Task/pkg/repo"
	"github.com/MahmoudMekki/Rescounts-Task/pkg/stripe"
	"github.com/MahmoudMekki/Rescounts-Task/pkg/token"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type PurchaseHandler struct {
	l            *log.Logger
	StripeClient stripe.Stripe
	StripeRepo   repo.StripeRepo
	tokenService token.Service
	userAccount  repo.UserAccountRepo
	productRepo  repo.ProductsRepo
}

func NewAPurchaseHandler(l *log.Logger, sc stripe.Stripe, sp repo.StripeRepo, tkn token.Service, user repo.UserAccountRepo, product repo.ProductsRepo) *PurchaseHandler {
	return &PurchaseHandler{
		l:            l,
		StripeClient: sc,
		StripeRepo:   sp,
		tokenService: tkn,
		userAccount:  user,
	}
}

func (user *PurchaseHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	tok, _ := user.tokenService.ExtractToken(req)
	currentUser, err := user.userAccount.GetUserByID(tok.UserID())
	if err != nil {
		resp := data.PurchaseResponse{Status: 0, Message: err.Error()}
		rhttp.RespondJSON(rw, http.StatusInternalServerError, resp)
		return
	}
	vars := mux.Vars(req)
	sProdID := vars["prod_id"]
	prodId, err := strconv.Atoi(sProdID)
	if err != nil {
		user.l.Println(err.Error())
		resp := data.PurchaseResponse{Status: 0, Message: "bad product id"}
		rhttp.RespondJSON(rw, http.StatusBadRequest, resp)
		return
	}
	cus, existed := user.StripeRepo.IsCustomer(currentUser.ID)
	if !existed {
		resp := data.PurchaseResponse{Status: 0, Message: "not a customer"}
		rhttp.RespondJSON(rw, http.StatusBadRequest, resp)
		return
	}
	prod, err := user.productRepo.GetProductByID(int64(prodId))
	if err != nil || prod.ID <= 0 {
		resp := data.PurchaseResponse{Status: 0, Message: "No product with this ID, try again!"}
		rhttp.RespondJSON(rw, http.StatusNoContent, resp)
		return
	}
	transactionID, err := user.StripeClient.ChargeCustomer(prod.PriceID, cus)
	if err != nil || transactionID == "" {
		user.l.Println(err.Error())
		resp := data.PurchaseResponse{Status: 0, Message: "unable to complete the payment"}
		rhttp.RespondJSON(rw, http.StatusInternalServerError, resp)
		return
	}
	resp := data.PurchaseResponse{Status: 1, Data: &data.PurchaseData{TransactionID: transactionID}}
	rhttp.RespondJSON(rw, http.StatusOK, resp)
}
