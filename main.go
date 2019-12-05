package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"

	toml "github.com/BurntSushi/toml"
)

const (
	// PacketStart is the flag for the start of a packet.
	PacketStart byte = 0x02

	// PacketHead is the flag for the end of a packet head.
	PacketHead byte = 0x01

	// PacketEnd is the flag for the end of a packet.
	PacketEnd byte = 0x03
)

// Server is structure for toml
type Server struct {
	IP      string `toml:"ip"`
	Port    int    `toml:"port"`
	Network string `toml:"network"`
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
	argIP      string
	argPort    int
	argNetwork string
	argMode    string
)

func init() {
	flag.StringVar(&argIP, "ip", "", "IP address")
	flag.IntVar(&argPort, "port", -1, "The service port")
	flag.StringVar(&argNetwork, "network", "tcp", "Network type. tcp or udp")
	flag.StringVar(&argMode, "mode", "stream", "Mode: stream or packet")
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

	ln, err := net.Listen(argNetwork, fmt.Sprintf("%s:%d", config.Server.IP, config.Server.Port))

	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handler(config, conn)
	}
}

func handler(config Config, server net.Conn) {
	log.Printf("[SERVER](%s,%s) -> (%s,%s)", server.LocalAddr().Network(), server.LocalAddr().String(), server.RemoteAddr().Network(), server.RemoteAddr().String())

	for _, c := range config.Clients {
		client, err := net.Dial(argNetwork, net.JoinHostPort(c.IP, fmt.Sprintf("%d", c.Port)))

		if err != nil {
			log.Fatal(err)
		} else {
			log.Printf("[CLIENT](%s,%s) -> (%s,%s)", client.LocalAddr().Network(), client.LocalAddr().String(), client.RemoteAddr().Network(), client.RemoteAddr().String())

			if "packet" == argMode {
				var buff bytes.Buffer

				b := make([]byte, 1024)

				for {
					n, err := server.Read(b)

					if err != nil {
						log.Fatal(err)
						break
					}

					if n < 0 {
						break
					}

					n, err = buff.Write(b[0:n])

					if err != nil {
						log.Fatal(err)
					}

					line, err := buff.ReadString('\003')

					bbuf := []byte(line)

					startIndex := bytes.IndexByte(bbuf, PacketStart)
					headIndex := bytes.IndexByte(bbuf, PacketHead)

					log.Printf(">> startIndex : %d, headIndex : %d", startIndex, headIndex)

					if headIndex > startIndex && startIndex >= 0 {
						packet := bbuf[startIndex:]

						log.Printf("packet : %s", string(packet))

						client.Write([]byte(packet))
					}
				}
			} else {
				go func() {
					_, err := io.Copy(server, client)

					if err != nil {
						log.Fatal(err)
					}
				}()

				_, err := io.Copy(client, server)

				if err != nil {
					log.Fatal(err)
				}
			}

			break
		}
	}
}
