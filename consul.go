package main


import (
	"log"
	"fmt"
	"net/http"
	"net"
	consul "github.com/hashicorp/consul/api"
)

var (
	count int64
)

func RegisterService(serviceName string, instanceId string, port int, checkPort int) {
	config := consul.DefaultConfig()

	config.Address = "consul:8500"

	client, err := consul.NewClient(config)

    if err != nil {
        log.Fatal("consul client error : ", err)
	}

	registration := new(consul.AgentServiceRegistration)
    registration.ID = instanceId
    registration.Name = serviceName
    registration.Port = port
    registration.Tags = []string{"Test"}
	registration.Address = localIP()
	
    registration.Check = &consul.AgentServiceCheck{
        HTTP:                           fmt.Sprintf("http://%s:%d%s", registration.Address, checkPort, "/check"),
        Timeout:                        "3s",
        Interval:                       "5s",
        DeregisterCriticalServiceAfter: "30s",
        // GRPC:     fmt.Sprintf("%v:%v/%v", IP, r.Port, r.Service),
    }

	err = client.Agent().ServiceRegister(registration)
	
    if err != nil {
        log.Fatal("register server error : ", err)
    }

    http.HandleFunc("/check", consulCheck)
    http.ListenAndServe(fmt.Sprintf(":%d", checkPort), nil)
}

func consulCheck(w http.ResponseWriter, r *http.Request) {
    s := "consulCheck" + fmt.Sprint(count) + "remote:" + r.RemoteAddr + " " + r.URL.String()
    fmt.Println(s)
    fmt.Fprintln(w, s)
    count++
}

func localIP() string {
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        return ""
	}
	
    for _, address := range addrs {
        if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                return ipnet.IP.String()
            }
        }
	}
	
    return ""
}
