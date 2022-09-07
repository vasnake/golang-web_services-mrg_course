package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

// local config
type Config struct {
	Comments bool `json:"comments"`
	Limit    int  // 0 by default
	Servers  []string
}

// not easy to share with other packages
var (
	config = &Config{}
)

func main() {
	// load config

	data, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatalln("cant read config file:", err)
	}

	err = json.Unmarshal(data, config)
	if err != nil {
		log.Fatalln("cant parse config:", err)
	}

	// use config

	if config.Comments {
		fmt.Println("Comments per page", config.Limit)
		fmt.Println("Comments services", config.Servers)
	} else {
		fmt.Println("Comments disabled")
	}
}
