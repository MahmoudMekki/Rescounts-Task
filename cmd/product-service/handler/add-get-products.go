package handler

import (
	"encoding/json"
	"github.com/MahmoudMekki/Rescounts-Task/cmd/product-service/data"
	"github.com/MahmoudMekki/Rescounts-Task/kit/rhttp"
	"github.com/MahmoudMekki/Rescounts-Task/pkg/model"
	"github.com/MahmoudMekki/Rescounts-Task/pkg/repo"
	"github.com/MahmoudMekki/Rescounts-Task/pkg/stripe"
	"github.com/MahmoudMekki/Rescounts-Task/pkg/token"
	"log"
	"net/http"
	"time"
)

type GetAddProductHandler struct {
	l            *log.Logger
	productRepo  repo.ProductsRepo
	StripeClient stripe.Stripe
	StripeRepo   repo.StripeRepo
	tokenService token.Service
}

func NewGetAddProductHandler(l *log.Logger, p repo.ProductsRepo, sc stripe.Stripe, sp repo.StripeRepo, tkn token.Service) *GetAddProductHandler {
	return &GetAddProductHandler{
		l:            l,
		productRepo:  p,
		StripeClient: sc,
		StripeRepo:   sp,
		tokenService: tkn,
	}
}

func (product *GetAddProductHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		product.addProduct(rw, req)
	case http.MethodGet:
		product.getProducts(rw, req)
	default:
		rhttp.RespondJSON(rw, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (product *GetAddProductHandler) getProducts(rw http.ResponseWriter, req *http.Request) {
	var addProduct data.AddProductRequest
	err := json.NewDecoder(req.Body).Decode(&addProduct)
	if err != nil {
		product.l.Println(err.Error())
		rhttp.RespondJSON(rw, http.StatusInternalServerError, "Unable to marshal request body")
		return
	}

	priceID, err := product.StripeClient.AddProduct(addProduct.Name, addProduct.Currency, addProduct.Price)
	if err != nil {
		product.l.Println(err.Error())
		rhttp.RespondJSON(rw, http.StatusInternalServerError, "Unable to create product and price on stripe")
		return
	}
	item := model.Product{
		Name:      addProduct.Name,
		Price:     addProduct.Price,
		Currency:  addProduct.Currency,
		PriceID:   priceID,
		CreatedAt: time.Now().UTC().String(),
	}
	prodID, err := product.productRepo.CreateProduct(item)
	if err != nil {
		product.l.Println(err.Error())
		rhttp.RespondJSON(rw, http.StatusInternalServerError, "Unable to create product ")
		return
	}
	resp := data.AddProductResponse{
		Message:   "Created successfully",
		ProductID: prodID,
	}
	rhttp.RespondJSON(rw, http.StatusCreated, resp)
}

func (product *GetAddProductHandler) addProduct(rw http.ResponseWriter, req *http.Request) {
	token, _ := product.tokenService.ExtractToken(req)
	if !token.IsAdmin() {
		rhttp.RespondJSON(rw, http.StatusUnauthorized, "Not allowed")
		return
	}
	var addProduct data.AddProductRequest
	err := json.NewDecoder(req.Body).Decode(&addProduct)
	if err != nil {
		product.l.Println(err.Error())
		rhttp.RespondJSON(rw, http.StatusInternalServerError, "Unable to marshal request body")
		return
	}

	priceID, err := product.StripeClient.AddProduct(addProduct.Name, addProduct.Currency, addProduct.Price)
	if err != nil {
		product.l.Println(err.Error())
		rhttp.RespondJSON(rw, http.StatusInternalServerError, "Unable to create product and price on stripe")
		return
	}
	item := model.Product{
		Name:      addProduct.Name,
		Price:     addProduct.Price,
		Currency:  addProduct.Currency,
		PriceID:   priceID,
		CreatedAt: time.Now().UTC().String(),
	}
	prodID, err := product.productRepo.CreateProduct(item)
	if err != nil {
		product.l.Println(err.Error())
		rhttp.RespondJSON(rw, http.StatusInternalServerError, "Unable to create product ")
		return
	}
	resp := data.AddProductResponse{
		Message:   "Created successfully",
		ProductID: prodID,
	}
	rhttp.RespondJSON(rw, http.StatusCreated, resp)
}
