package common

import (
	"encoding/json"
	"log"
	"os"
)

type configuration struct {
	Protocol string `json:"protocol"`
	Server   string `json:"server"`
	DBUser   string `json:"databaseUser"`
	DBPwd    string `json:"databasePassword"`
	Database string `json:"database"`
}

// AppConfig holds the configuration values from config.json
var AppConfig configuration

func InitConfig() {
	file, err := os.Open("common/config.json")
	defer file.Close()
	if err != nil {
		log.Fatalf("[loadConfig]: %s\n", err)
	}
	decoder := json.NewDecoder(file)
	AppConfig = configuration{}
	err = decoder.Decode(&AppConfig)
	if err != nil {
		log.Fatalf("[loadAppConfig]: %s\n", err)
	}
}
