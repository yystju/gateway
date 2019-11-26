package main

import (
	"fmt"
	"io/ioutil"
	"log"

	toml "github.com/BurntSushi/toml"
)

// Server is structure for toml
type Server struct {
	IP   string `toml:"ip"`
	Port int    `toml:"port"`
}

// Client is structure for toml
type Client struct {
	Name string `toml:"name"`
	IP   string `toml:"ip"`
	Port int    `toml:"port"`
}

// Config is structure for toml
type Config struct {
	Server  Server   `toml:"server"`
	Clients []Client `toml:"client"`
}

func main() {
	fmt.Println("MENTOS GATEWAY")

	f, err := ioutil.ReadFile("config.toml")

	if err != nil {
		log.Fatal(err)
	}

	var config Config

	_, err = toml.Decode(string(f), &config)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("[SERVER] -> %s:%d\n", config.Server.IP, config.Server.Port)

	for _, c := range config.Clients {
		fmt.Printf("[CLIENT] %s -> %s:%d\n", c.Name, c.IP, c.Port)
	}
}
