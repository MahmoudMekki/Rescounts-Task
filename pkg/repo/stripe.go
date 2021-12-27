package repo

import (
	"database/sql"
	"github.com/MahmoudMekki/Rescounts-Task/pkg/model"
)

type StripeRepo interface {
	CreateCustomer(customer model.StripeCustomer) error
	IsCustomer(userID int64) (string, bool, error)
	UpdateCustomer(customer model.StripeCustomer) error
}

func NewStripeRepo(db *sql.DB) StripeRepo {
	return &stripImp{DBEngine: db}
}

type stripImp struct {
	DBEngine *sql.DB
}

func (s *stripImp) CreateCustomer(customer model.StripeCustomer) error {
	_, err := s.DBEngine.Exec("insert into rescounts.stripe_customers (user_id,customer_id,created_at) values ($1,$2,$3)", customer.UserID, customer.CustomerID, customer.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *stripImp) IsCustomer(userID int64) (string, bool,error) {
	var customer model.StripeCustomer
	row := s.DBEngine.QueryRow("select * from rescounts.stripe_customers where user_id=$1", userID)
	err := row.Scan(&customer.UserID, &customer.CustomerID, &customer.CreatedAt)
	 if err == sql.ErrNoRows{
		 return "",false,nil
	}
	return customer.CustomerID, true,nil
}
func (s *stripImp) UpdateCustomer(customer model.StripeCustomer) error {
	_, err := s.DBEngine.Exec("update rescounts.stripe_customers set customer_id=$1 where user_id=$2", customer.CustomerID, customer.UserID)
	if err != nil {
		return err
	}
	return nil
}
