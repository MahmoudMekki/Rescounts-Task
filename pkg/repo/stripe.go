
package repo

import (
	"database/sql"
	"github.com/MahmoudMekki/Rescounts-Task/pkg/model"
)

type StripeRepo interface {
	CreateCustomer(customer model.StripeCustomer) error
	IsCustomer(userID int64)(bool,string)
}

func NewStripeRepo(db *sql.DB) StripeRepo {
	return &stripImp{DBEngine: db}
}

type stripImp struct {
	DBEngine *sql.DB
}

func (s *stripImp)CreateCustomer(customer model.StripeCustomer) error{
	_,err := s.DBEngine.Exec("insert into stripe_customers (user_id,customer_id,created_at) values ($1,$2,$3)",customer.UserID,customer.CustomerID,customer.CreatedAt)
	if err !=nil{
		return err
	}
	return nil
}

func (s *stripImp)IsCustomer(userID int64)(bool,string){
	var customer model.StripeCustomer
	row := s.DBEngine.QueryRow("select * from stripe_customer where user_id=$1",userID)
	err := row.Scan(&customer.UserID,&customer.CustomerID,&customer.CreatedAt)
	if err !=nil || customer.UserID <= 0{
		return false,""
	}
	return true,customer.CustomerID
}
