package configs

import (
	"encoding/json"
	"log"
	"os"
)

type Configs struct {
	Owner    []string `json:"owner"`
	Public   bool     `json:"public"`
	Prefix   string   `json:"prefix"`
	Premium  []string `json:"premium"`
	AntiCall bool     `json:"anticall"`
}

var CONFIG Configs

func init() {
	data, err := os.ReadFile("./configs.json")
	if err != nil {
		log.Fatalln("Error reading file configs.json")
	}
	if err := json.Unmarshal(data, &CONFIG); err != nil {
		log.Fatalln(err)
	}
}
