package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

var DB *sql.DB

type DataBase struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
	Debug    bool   `json:"debug"`
	Dialect  string `json:"dialect"`
}

func (d *DataBase) OpenDB() *sql.DB {
	if DB == nil {
		dsn := d.dsn()
		DB, err := sql.Open("postgres", dsn)
		if err != nil {
			log.Panicln(err.Error())
		}
		return DB
	}
	return DB
}

func (d *DataBase) dsn() (dsn string) {
	dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", d.Host, d.User, d.Password, d.Name, d.Port)
	return dsn
}
