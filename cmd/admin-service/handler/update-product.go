package handler

import (
	"encoding/json"
	"github.com/MahmoudMekki/Rescounts-Task/cmd/admin-service/data"
	"github.com/MahmoudMekki/Rescounts-Task/kit/rhttp"
	"github.com/MahmoudMekki/Rescounts-Task/pkg/repo"
	"github.com/MahmoudMekki/Rescounts-Task/pkg/stripe"
	"log"
	"net/http"
	"strconv"
)

type UpdateProductHandler struct {
	l            *log.Logger
	productRepo  repo.ProductsRepo
	StripeClient stripe.Stripe
	StripeRepo   repo.StripeRepo
}

func NewUpdateProductHandler(l *log.Logger, p repo.ProductsRepo, sc stripe.Stripe, sp repo.StripeRepo) *UpdateProductHandler {
	return &UpdateProductHandler{
		l:            l,
		productRepo:  p,
		StripeClient: sc,
		StripeRepo:   sp,
	}
}

func (product *UpdateProductHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPut {
		rhttp.RespondJSON(rw, http.StatusMethodNotAllowed, "Not allowed method")
		return
	}
	var updateProduct data.UpdateProductRequest
	err := json.NewDecoder(req.Body).Decode(&updateProduct)
	if err != nil {
		product.l.Println(err.Error())
		rhttp.RespondJSON(rw, http.StatusInternalServerError, "Unable to marshal request body")
		return
	}
	sprodID := req.FormValue("prod_id")
	prodId, err := strconv.Atoi(sprodID)
	if err != nil {
		product.l.Println(err.Error())
		rhttp.RespondJSON(rw, http.StatusNotAcceptable, "bad product id")
		return
	}

	prod, err := product.productRepo.GetProductByID(int64(prodId))
	if err != nil {
		rhttp.RespondJSON(rw, http.StatusNotFound, err.Error())
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
		rhttp.RespondJSON(rw, http.StatusInternalServerError, "Unable add the new prod on stripe")
		product.l.Println(err.Error())
		return
	}
	prod.PriceID = priceID
	err = product.productRepo.UpdateProduct(prod)
	if err != nil {
		rhttp.RespondJSON(rw, http.StatusInternalServerError, "Unable to update the product")
		product.l.Println(err.Error())
		return
	}
	resp := data.UpdateProductResponse{
		Message:   "updated successfully",
		ProductID: prod.ID,
	}
	rhttp.RespondJSON(rw, http.StatusOK, resp)
}
