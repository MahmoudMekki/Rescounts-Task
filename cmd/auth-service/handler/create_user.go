package handler

import (
	"encoding/json"
	"errors"
	"github.com/MahmoudMekki/Rescounts-Task/cmd/auth-service/data"
	codes "github.com/MahmoudMekki/Rescounts-Task/kit/error-codes"
	"github.com/MahmoudMekki/Rescounts-Task/kit/rhttp"
	"github.com/MahmoudMekki/Rescounts-Task/pkg/model"
	"github.com/MahmoudMekki/Rescounts-Task/pkg/password"
	"github.com/MahmoudMekki/Rescounts-Task/pkg/regex"
	"github.com/MahmoudMekki/Rescounts-Task/pkg/repo"
	"github.com/MahmoudMekki/Rescounts-Task/pkg/token"
	"log"
	"net/http"
	"time"
)

type CreateUserAccountHandler struct {
	l            *log.Logger
	userAccount  repo.UserAccountRepo
	tokenService token.Service
}

func NewCreateUserAccountHandler(l *log.Logger, u repo.UserAccountRepo, tkn token.Service) *CreateUserAccountHandler {
	return &CreateUserAccountHandler{
		l:            l,
		userAccount:  u,
		tokenService: tkn,
	}
}

func (u *CreateUserAccountHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	var signUpRequest data.CreateAccountRequest
	err := json.NewDecoder(req.Body).Decode(&signUpRequest)
	if err != nil {
		u.l.Println(err.Error())
		rhttp.RespondJSON(rw, http.StatusInternalServerError, "Unable to marshal request body")
		return
	}
	code, err := u.validate(signUpRequest)
	if code != codes.SuccessfulState {
		if code == codes.InternalServerError {
			u.l.Println(err.Error())
			rhttp.RespondJSON(rw, http.StatusInternalServerError, err.Error())
			return
		}
		rhttp.RespondJSON(rw, http.StatusNotAcceptable, err.Error())
		return
	}
	userAccount := &model.UserAccount{
		Email:     signUpRequest.Email,
		Password:  password.Hash([]byte(signUpRequest.Password)),
		FirstName: signUpRequest.FirstName,
		LastName:  signUpRequest.LastName,
		IsAdmin:   false,
		CreatedAt: time.Now().UTC().String(),
		UpdatedAt: time.Now().UTC().String(),
	}
	userID, err := u.userAccount.CreateUser(userAccount)
	if err != nil {
		u.l.Println(err.Error())
		rhttp.RespondJSON(rw, http.StatusInternalServerError, err.Error())
		return
	}
	token, err := u.tokenService.New(userID, 24*time.Hour)
	if err != nil {
		rhttp.RespondJSON(rw, http.StatusInternalServerError, err.Error())
		return
	}
	resp := data.CreateAccountResponse{UserID: userID, AccessToken: token.Signed()}
	rhttp.RespondJSON(rw, http.StatusCreated, resp)
}

func (u *CreateUserAccountHandler) validate(account data.CreateAccountRequest) (code int, err error) {
	if account.Email == "" || !regex.MatchEmail(account.Email) {
		return codes.BadRequestArguments, errors.New("provided email is invalid")
	}
	if len(account.Password) < 8 {
		return codes.BadRequestArguments, errors.New("provided password isn't satisfying")
	}
	user, err := u.userAccount.GetUserByEmail(account.Email)
	if err != nil {
		return codes.InternalServerError, err
	}
	if user.ID > 0 {
		return codes.BadRequestArguments, errors.New("email is already taken")
	}
	return codes.SuccessfulState, nil
}
