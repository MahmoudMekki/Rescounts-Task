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
	var user model.UserAccount
	row := u.DBEngine.QueryRow("select "+
		"id,first_name,last_name,email,password,is_admin"+
		" from rescounts.users where email=$1", email)
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.IsAdmin)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}
	return &user, nil
}

func (u *userAcccountImp) GetUserByID(id int64) (*model.UserAccount, error) {
	var user model.UserAccount
	row := u.DBEngine.QueryRow("select "+
		"id,first_name,last_name,email,password,is_admin"+
		" from rescounts.users where email=$1", id)
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.IsAdmin)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}
	return &user, nil
}

func (u *userAcccountImp) CreateUser(user *model.UserAccount) (userId int64, err error) {
	stmnt := `INSERT INTO rescounts.users (first_name,last_name,email,password,is_admin,created_at) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id`
	err = u.DBEngine.QueryRow(stmnt, user.FirstName, user.LastName, user.Email, user.Password, user.IsAdmin, user.CreatedAt).Scan(&userId)
	if err != nil {
		return 0, err
	}
	return userId, nil
}
