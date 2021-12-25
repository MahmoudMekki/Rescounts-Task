package stripe

import (
	"fmt"
	stripe "github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/charge"
	"github.com/stripe/stripe-go/v72/customer"
	price "github.com/stripe/stripe-go/v72/price"
	"github.com/stripe/stripe-go/v72/product"
	"github.com/stripe/stripe-go/v72/token"
)

type Stripe interface {
	AddProduct(sku, currency string, amount float64) (string, error)
	ChargeCustomer(priceID, customerID string) (string, error)
	CreateCardToken(cardNmuber, expMonth, expYear, cvc string) (string, error)
	CreateCustomer(email, cardToken, userName string) (string, error)
	UpdateCustomer(customerID, cardToken string) error
}

func NewStripe(key string) Stripe {
	return &Client{APIKey: key}
}

type Client struct {
	APIKey string
}

func (c *Client) CreateCardToken(cardNmuber, expMonth, expYear, cvc string) (string, error) {
	stripe.Key = c.APIKey
	cardParm := &stripe.CardParams{
		CVC:      stripe.String(cvc),
		ExpMonth: stripe.String(expMonth),
		ExpYear:  stripe.String(expYear),
		Number:   stripe.String(cardNmuber),
	}
	tok, err := token.New(&stripe.TokenParams{Card: cardParm})
	if err != nil {
		return "", err
	}
	return tok.ID, nil
}

func (c *Client) CreateCustomer(email, cardToken, userName string) (string, error) {
	stripe.Key = c.APIKey
	parameters := stripe.CustomerParams{
		Email:  stripe.String(email),
		Name:   stripe.String(userName),
		Source: &stripe.SourceParams{Token: stripe.String(cardToken)},
	}
	customer, err := customer.New(&parameters)
	if err != nil {
		return "", err
	}
	return customer.ID, nil
}
func (c *Client) UpdateCustomer(customerID, cardToken string) error {
	stripe.Key = c.APIKey
	params := &stripe.CustomerParams{Source: &stripe.SourceParams{Token: stripe.String(cardToken)}}
	_, err := customer.Update(customerID, params)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) AddProduct(sku, currency string, amount float64) (string, error) {
	stripe.Key = c.APIKey
	prod, err := c.addProduct(sku)
	if err != nil {
		return "", err
	}
	price, err := c.addPrice(currency, prod, amount)
	if err != nil {
		return "", err
	}
	return price, nil
}

func (c *Client) addProduct(sku string) (string, error) {
	stripe.Key = c.APIKey
	parameters := stripe.ProductParams{
		Name: stripe.String(sku),
	}
	prod, err := product.New(&parameters)
	if err != nil {
		return "", err
	}
	return prod.ID, nil
}
func (c *Client) addPrice(currency, product string, amount float64) (string, error) {
	stripe.Key = c.APIKey
	parameters := stripe.PriceParams{
		Currency:          stripe.String(currency),
		Product:           stripe.String(product),
		UnitAmountDecimal: stripe.Float64(amount),
	}
	price, err := price.New(&parameters)
	if err != nil {
		return "", err
	}
	return price.ID, nil
}
func (c *Client) ChargeCustomer(priceID, customerID string) (string, error) {
	stripe.Key = c.APIKey
	item, err := price.Get(priceID, nil)
	if err != nil {
		return "", err
	}
	params := stripe.ChargeParams{
		Amount:   stripe.Int64(item.UnitAmount),
		Currency: stripe.String(string(item.Currency)),
		Customer: stripe.String(customerID),
	}
	ch, err := charge.New(&params)
	if err != nil {
		return "", err
	}
	if ch.Status == "succeeded" {
		return ch.ID, nil
	}
	return "", fmt.Errorf("can't perform the transaction %v", err)
}
