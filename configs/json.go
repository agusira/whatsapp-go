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
		log.Println("Error reading file configs.json")
		file, err := os.Create("./configs.json")
		if err != nil {
			log.Fatalln("Error creating file configs.json")
		}
		defer file.Close()
		newcfg := Configs{
			Owner:    []string{"62895359263399"},
			Public:   false,
			Prefix:   "!",
			Premium:  []string{"62895359263399"},
			AntiCall: false,
		}
		cfgbyte, err := json.Marshal(newcfg)
		if err != nil {
			log.Fatalln("Error Marshal config")
		}
		file.WriteString(string(cfgbyte))
	}
	if err := json.Unmarshal(data, &CONFIG); err != nil {
		log.Fatalln(err)
	}
}
