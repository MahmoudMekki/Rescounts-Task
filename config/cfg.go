package config

import (
	"encoding/json"
	"errors"
	"github.com/MahmoudMekki/Rescounts-Task/database"
	"github.com/MahmoudMekki/Rescounts-Task/server"
	"io/ioutil"
	"log"
)

type Config struct {
	DataBase  database.DataBase `json:"database"`
	JWT       JWT               `json:"jwt"`
	StripeKey string            `json:"stripe"`
	Http      server.Http       `json:"http"`
}

type JWT struct {
	Secret string `json:"secret"`
}

func LoadConfig() (cfg Config) {
	file, err := ioutil.ReadFile("config/devEnv.json")
	if err != nil {
		log.Panicln(err)
	}
	err = json.Unmarshal(file, &cfg)
	if err != nil {
		log.Panicln(errors.New("can't marshal the config file"))
	}
	return cfg
}
