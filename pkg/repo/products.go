package repo

import (
	"database/sql"
	"github.com/MahmoudMekki/Rescounts-Task/pkg/model"
)

type ProductsRepo interface {
	CreateProduct(product model.Product) (productID int64,err error)
	UpdateProduct(product model.Product) error
	GetProducts() (products []model.Product, err error)
}

func NewProductsRepo(db *sql.DB) ProductsRepo {
	return &productsImp{DBEngine: db}
}

type productsImp struct {
	DBEngine *sql.DB
}

func (p *productsImp) CreateProduct(product model.Product) (productID int64,err error) {
	stmnt := `INSERT INTO rescounts.products (name,price,currency,created_at) VALUES ($1,$2,$3,$4) RETURNING id`
	err = p.DBEngine.QueryRow(stmnt, product.Name, product.Price, product.Currency, product.CreatedAt).Scan(&productID)
	if err != nil {
		return 0, err
	}
	return productID, nil
}

func (p *productsImp) UpdateProduct(product model.Product) error {
	_,err := p.DBEngine.
		Exec("update rescounts.products set name=$1,price=$2,currency=$3 where id =$4",product.Name,product.Price,product.Currency,product.ID)
	if err !=nil{
		return err
	}
	return nil
}

// GetProducts it will be done in pagination but later coz there is no time

func (p *productsImp) GetProducts() (products []model.Product, err error) {
	rows,err := p.DBEngine.Query("select id,name,price,currency from rescounts.products")
	if err !=nil{
		return nil,err
	}
	defer rows.Close()
	for rows.Next() {
		var product model.Product
		err = rows.Scan(&product.ID, &product.Name, &product.Price, &product.Currency)

		if err != nil {
			return nil,err
		}
		products = append(products, product)
	}
	return products,nil
}

