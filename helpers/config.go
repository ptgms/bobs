package helpers

import (
	"encoding/json"
	"os"
)

type Server struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Priority int    `json:"priority"`
}

type Domain struct {
	Zone string `json:"zone"`
	Name string `json:"name"`
}

type Config struct {
	Servers            []Server `json:"servers"`
	IP                 string   `json:"ip"`
	Domains            []Domain `json:"domains"`
	CloudflareAPIToken string   `json:"cloudflare_api_token"`
}

func SaveEmptyConfig() {
	var config Config
	config.Servers = []Server{
		{
			Name:     "Server 1",
			Host:     "https://server1.example.com",
			Priority: 1,
		},
	}

	config.Domains = []Domain{
		{
			Name: "example.com",
		},
	}

	config.CloudflareAPIToken = ""

	// open file
	file, err := os.Create("config.json")
	HandleError(err, false)
	defer file.Close()

	b, err := json.MarshalIndent(config, "", "    ")
	HandleError(err, false)

	_, err = file.Write(b)
	HandleError(err, false)
}

func LoadConfig() Config {
	file, err := os.Open("config.json")
	if err != nil {
		SaveEmptyConfig()
		println("Please fill out config.json and restart the program.")
		os.Exit(1)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := Config{}
	err = decoder.Decode(&config)
	HandleError(err, false)

	return config
}

func HandleError(err error, fatal bool) {
	if err != nil {
		if fatal {
			panic(err)
		} else {
			println(err.Error())
		}
	}
}
