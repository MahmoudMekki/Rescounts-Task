package model

const (
	UserAcountsInfoTableName = "users"
)

type UserAccount struct {
	ID        int64
	Email     string
	Password  string
	FirstName string
	LastName  string
	IsAdmin   bool
	CreatedAt string
	UpdatedAt string
}
