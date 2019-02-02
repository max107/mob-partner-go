package main

import (
	"os"
	"io/ioutil"
	"encoding/json"
	"github.com/pkg/errors"
)

func main() {
	if len(os.Args) < 2 {
		errors.New("Missing config file path")
	}

	jsonConfig, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer jsonConfig.Close()
	byteValue, _ := ioutil.ReadAll(jsonConfig)

	var config AppConfig
	json.Unmarshal([]byte(byteValue), &config)

	NewService(config).Listen()
}
