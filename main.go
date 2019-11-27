package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"

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

var (
	argIP   string
	argPort int
)

func init() {
	flag.StringVar(&argIP, "ip", "", "IP address")
	flag.IntVar(&argPort, "port", -1, "The service port")
	flag.Parse()
}

func main() {
	log.Println("[MENTOS GATEWAY]")

	f, err := ioutil.ReadFile("config.toml")

	if err != nil {
		log.Fatal(err)
	}

	var config Config

	_, err = toml.Decode(string(f), &config)

	if err != nil {
		log.Fatal(err)
	}

	if "" != argIP {
		config.Server.IP = argIP
	}

	if argPort >= 0 {
		config.Server.Port = argPort
	}

	log.Printf("[SERVER] -> %s:%d\n", config.Server.IP, config.Server.Port)

	for _, c := range config.Clients {
		log.Printf("[CLIENT] %s -> %s:%d\n", c.Name, c.IP, c.Port)
	}

	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", config.Server.IP, config.Server.Port))

	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handler(conn)
	}
}

func handler(conn net.Conn) {
	log.Printf("(%s,%s) -> (%s,%s)", conn.LocalAddr().Network(), conn.LocalAddr().String(), conn.RemoteAddr().Network(), conn.RemoteAddr().String())

	conn.Write([]byte("NG"))
	conn.Close()
}
