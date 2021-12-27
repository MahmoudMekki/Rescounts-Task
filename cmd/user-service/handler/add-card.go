package handler

import (
	"encoding/json"
	"github.com/MahmoudMekki/Rescounts-Task/cmd/user-service/data"
	"github.com/MahmoudMekki/Rescounts-Task/kit/rhttp"
	"github.com/MahmoudMekki/Rescounts-Task/pkg/model"
	"github.com/MahmoudMekki/Rescounts-Task/pkg/repo"
	"github.com/MahmoudMekki/Rescounts-Task/pkg/stripe"
	"github.com/MahmoudMekki/Rescounts-Task/pkg/token"
	"log"
	"net/http"
	"time"
)

type AddCardHandler struct {
	l            *log.Logger
	StripeClient stripe.Stripe
	StripeRepo   repo.StripeRepo
	tokenService token.Service
	userAccount  repo.UserAccountRepo
}

func NewAddCardHandler(l *log.Logger, sc stripe.Stripe, sp repo.StripeRepo, tkn token.Service, user repo.UserAccountRepo) *AddCardHandler {
	return &AddCardHandler{
		l:            l,
		StripeClient: sc,
		StripeRepo:   sp,
		tokenService: tkn,
		userAccount:  user,
	}
}

func (user *AddCardHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	tok, _ := user.tokenService.ExtractToken(req)
	currentUser, err := user.userAccount.GetUserByID(tok.UserID())
	if err != nil {
		resp := data.AddCardResponse{Status: 0, Message: err.Error()}
		rhttp.RespondJSON(rw, http.StatusInternalServerError, resp)
		return
	}
	var cardData data.AddCardRequest
	err = json.NewDecoder(req.Body).Decode(&cardData)
	if err != nil {
		user.l.Println(err.Error())
		resp := data.AddCardResponse{Status: 0, Message: err.Error()}
		rhttp.RespondJSON(rw, http.StatusInternalServerError, resp)
		return
	}
	cardTok, err := user.StripeClient.CreateCardToken(
		cardData.CardNumber,
		cardData.ExpireMonth,
		cardData.ExpireYear,
		cardData.CVC,
	)
	if err != nil {
		user.l.Println(err.Error())
		resp := data.AddCardResponse{Status: 0, Message: err.Error()}
		rhttp.RespondJSON(rw, http.StatusInternalServerError, resp)
		return
	}
	cus, existed := user.StripeRepo.IsCustomer(currentUser.ID)
	if !existed {
		cusID, err := user.StripeClient.CreateCustomer(currentUser.Email, cardTok, currentUser.FirstName)
		if err != nil {
			resp := data.AddCardResponse{Status: 0, Message: err.Error()}
			rhttp.RespondJSON(rw, http.StatusInternalServerError, resp)
			return
		}
		stripeCustomer := model.StripeCustomer{UserID: tok.UserID(), CustomerID: cusID, CreatedAt: time.Now().UTC().String()}
		err = user.StripeRepo.CreateCustomer(stripeCustomer)
		if err != nil {
			resp := data.AddCardResponse{Status: 0, Message: err.Error()}
			rhttp.RespondJSON(rw, http.StatusInternalServerError, resp)
			return
		}
		resp := data.AddCardResponse{Status: 1}
		rhttp.RespondJSON(rw, http.StatusAccepted, resp)
		return
	}
	err = user.StripeClient.UpdateCustomer(cus, cardTok)
	if err != nil {
		resp := data.AddCardResponse{Status: 0, Message: err.Error()}
		rhttp.RespondJSON(rw, http.StatusInternalServerError, resp)
		return
	}
	resp := data.AddCardResponse{Status: 1}
	rhttp.RespondJSON(rw, http.StatusAccepted, resp)
	return
}
