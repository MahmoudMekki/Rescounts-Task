package config

import (
	"encoding/json"
	"errors"
	"github.com/MahmoudMekki/Rescounts-Task/database"
	"io/ioutil"
	"log"
)

type Config struct {
	DataBase  database.DataBase `json:"database"`
	JWT       JWT               `json:"jwt"`
	StripeKey string            `json:"stripe"`
}

type JWT struct {
	Secret string `json:"secret"`
}

func (cfg *Config) LoadConfig() {
	var config Config
	file, err := ioutil.ReadFile("config/devEnv.json")
	if err != nil {
		log.Panicln(err)
	}
	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Panicln(errors.New("can't marshal the config file"))
	}
}
