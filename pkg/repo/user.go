package repo

import (
	"database/sql"
	"github.com/MahmoudMekki/Rescounts-Task/pkg/model"
)

type UserAccountRepo interface {
	GetUserByEmail(email string) (*model.UserAccount, error)
	GetUserByID(id int64) (*model.UserAccount, error)
	CreateUser(user *model.UserAccount) (userId int64, err error)
}

func NewUserAccountRepo(db *sql.DB) UserAccountRepo {
	return &userAcccountImp{DBEngine: db}
}

type userAcccountImp struct {
	DBEngine *sql.DB
}

func (u *userAcccountImp) GetUserByEmail(email string) (*model.UserAccount, error) {

}

func (u *userAcccountImp) GetUserByID(id int64) (*model.UserAccount, error) {

}

func (u *userAcccountImp) CreateUser(user *model.UserAccount) (userId int64, err error) {

}
