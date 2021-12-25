package handler

import (
	"encoding/json"
	"github.com/MahmoudMekki/Rescounts-Task/cmd/auth-service/data"
	"github.com/MahmoudMekki/Rescounts-Task/kit/rhttp"
	"github.com/MahmoudMekki/Rescounts-Task/pkg/password"
	"github.com/MahmoudMekki/Rescounts-Task/pkg/repo"
	"github.com/MahmoudMekki/Rescounts-Task/pkg/token"
	"log"
	"net/http"
	"time"
)

type LoginHandler struct {
	l            *log.Logger
	userAccount  repo.UserAccountRepo
	tokenService token.Service
}

func NewLoginHandler(l *log.Logger, u repo.UserAccountRepo, tkn token.Service) *LoginHandler {
	return &LoginHandler{
		l:            l,
		userAccount:  u,
		tokenService: tkn,
	}
}

func (login *LoginHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		rhttp.RespondJSON(rw, http.StatusMethodNotAllowed, "Not allowed method")
		return
	}
	var loginRequest data.LoginRequest
	err := json.NewDecoder(req.Body).Decode(&loginRequest)
	if err != nil {
		login.l.Println(err.Error())
		rhttp.RespondJSON(rw, http.StatusInternalServerError, "Unable to marshal request body")
		return
	}

	userAccount,err :=login.userAccount.GetUserByEmail(loginRequest.Email)
	if err !=nil{
		rhttp.RespondJSON(rw,http.StatusInternalServerError,err.Error())
		return
	}
	if userAccount.ID<=0{
		rhttp.RespondJSON(rw,http.StatusNotFound,"this email not registered")
		return
	}
	if !password.ComparePasswords(userAccount.Password,[]byte(loginRequest.Password)){
		rhttp.RespondJSON(rw,http.StatusNotAcceptable,"Wrong Password")
		return
	}

	if userAccount.IsAdmin{
		accessToken, err := login.tokenService.NewAdminToken(userAccount.ID, 24*time.Hour)
		if err != nil {
			rhttp.RespondJSON(rw, http.StatusInternalServerError, err.Error())
			return
		}
		resp := data.CreateAccountResponse{UserID: userAccount.ID, AccessToken: accessToken.Signed()}
		rhttp.RespondJSON(rw, http.StatusOK, resp)
		return
	}

	accessToken, err := login.tokenService.New(userAccount.ID, 24*time.Hour)
	if err != nil {
		rhttp.RespondJSON(rw, http.StatusInternalServerError, err.Error())
		return
	}
	resp := data.CreateAccountResponse{UserID: userAccount.ID, AccessToken: accessToken.Signed()}
	rhttp.RespondJSON(rw, http.StatusOK, resp)
}



