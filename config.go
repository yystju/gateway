package main


// Server is structure for toml
type Server struct {
	IP      string `toml:"ip"`
	Port    int    `toml:"port"`
	Network string `toml:"network"`
	CheckPort    int    `toml:"checkPort"`
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

