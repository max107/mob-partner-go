package main

import (
	"os"
	"io/ioutil"
	"encoding/json"
)

func main() {
	jsonConfig, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}
	defer jsonConfig.Close()
	byteValue, _ := ioutil.ReadAll(jsonConfig)

	var config AppConfig
	json.Unmarshal([]byte(byteValue), &config)

	NewService(config).Listen()
}
