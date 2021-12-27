package handler

import (
	"encoding/json"
	"github.com/MahmoudMekki/Rescounts-Task/cmd/product-service/data"
	"github.com/MahmoudMekki/Rescounts-Task/kit/rhttp"
	"github.com/MahmoudMekki/Rescounts-Task/pkg/repo"
	"github.com/MahmoudMekki/Rescounts-Task/pkg/stripe"
	"github.com/MahmoudMekki/Rescounts-Task/pkg/token"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type UpdateProductHandler struct {
	l            *log.Logger
	productRepo  repo.ProductsRepo
	StripeClient stripe.Stripe
	StripeRepo   repo.StripeRepo
	tokenService token.Service
}

func NewUpdateProductHandler(l *log.Logger, p repo.ProductsRepo, sc stripe.Stripe, sp repo.StripeRepo, tkn token.Service) *UpdateProductHandler {
	return &UpdateProductHandler{
		l:            l,
		productRepo:  p,
		StripeClient: sc,
		StripeRepo:   sp,
		tokenService: tkn,
	}
}

func (product *UpdateProductHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	token, _ := product.tokenService.ExtractToken(req)
	if !token.IsAdmin() {
		rhttp.RespondJSON(rw, http.StatusUnauthorized, "Not allowed")
		return
	}
	var updateProduct data.UpdateProductRequest
	err := json.NewDecoder(req.Body).Decode(&updateProduct)
	if err != nil {
		product.l.Println(err.Error())
		resp := data.UpdateProductResponse{Status: 0, Message: "Unable to marshal request body"}
		rhttp.RespondJSON(rw, http.StatusInternalServerError, resp)
		return
	}
	vars := mux.Vars(req)
	sProdID := vars["prod_id"]
	prodId, err := strconv.Atoi(sProdID)
	if err != nil {
		product.l.Println(err.Error())
		resp := data.UpdateProductResponse{Status: 0, Message: "bad product id"}
		rhttp.RespondJSON(rw, http.StatusNotAcceptable, resp)
		return
	}

	prod, err := product.productRepo.GetProductByID(int64(prodId))
	if err != nil {
		resp := data.UpdateProductResponse{Status: 0, Message: err.Error()}
		rhttp.RespondJSON(rw, http.StatusNotFound, resp)
		return
	}

	if updateProduct.Name != prod.Name && updateProduct.Name != "" {
		prod.Name = updateProduct.Name
	}
	if updateProduct.Currency != prod.Currency && updateProduct.Currency != "" {
		prod.Currency = updateProduct.Currency
	}
	if updateProduct.Price != prod.Price && updateProduct.Price > 0 {
		prod.Price = updateProduct.Price
	}
	priceID, err := product.StripeClient.AddProduct(prod.Name, prod.Currency, prod.Price)
	if err != nil {
		resp := data.UpdateProductResponse{Status: 0, Message: err.Error()}
		rhttp.RespondJSON(rw, http.StatusInternalServerError, resp)
		product.l.Println(err.Error())
		return
	}
	prod.PriceID = priceID
	err = product.productRepo.UpdateProduct(prod)
	if err != nil {
		resp := data.UpdateProductResponse{Status: 0, Message: err.Error()}
		rhttp.RespondJSON(rw, http.StatusInternalServerError, resp)
		product.l.Println(err.Error())
		return
	}
	resp := data.UpdateProductResponse{
		Status: 1,
		Data:   &data.UpdateProductData{ProductID: prod.ID},
	}
	rhttp.RespondJSON(rw, http.StatusOK, resp)
}
