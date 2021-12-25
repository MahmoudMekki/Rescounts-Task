package handler

import (
	"encoding/json"
	"github.com/MahmoudMekki/Rescounts-Task/cmd/admin-service/data"
	"github.com/MahmoudMekki/Rescounts-Task/kit/rhttp"
	"github.com/MahmoudMekki/Rescounts-Task/pkg/model"
	"github.com/MahmoudMekki/Rescounts-Task/pkg/repo"
	"github.com/MahmoudMekki/Rescounts-Task/pkg/stripe"
	"log"
	"net/http"
	"time"
)

type AddProductHandler struct {
	l            *log.Logger
	productRepo  repo.ProductsRepo
	StripeClient stripe.Stripe
	StripeRepo   repo.StripeRepo
}

func NewAddProductHandler(l *log.Logger, p repo.ProductsRepo, sc stripe.Stripe, sp repo.StripeRepo) *AddProductHandler {
	return &AddProductHandler{
		l:            l,
		productRepo:  p,
		StripeClient: sc,
		StripeRepo:   sp,
	}
}

func (product *AddProductHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		rhttp.RespondJSON(rw, http.StatusMethodNotAllowed, "Not allowed method")
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
