package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"gopkg.in/sirupsen/logrus.v1"
	"gopkg.in/spf13/viper.v1"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
)

type DataResponse struct {
	Hostname         string      `json:"hostname,omitempty"`
	IP               []string    `json:"ip,omitempty"`
	Headers          http.Header `json:"header,omitempty"`
	Environment      []string    `json:"env,omitempty"`
}
var port string

func init() {
	flag.StringVar(&port, "port", "8081", "give me a port number")

	lvl, err := logrus.ParseLevel(viper.GetString("loglevel"))
	if err != nil {
		lvl = logrus.WarnLevel
	}
	logrus.SetLevel(lvl)
}

func main() {
	flag.Parse()

	http.HandleFunc("/", index)
	http.HandleFunc("/api", api)

	log.Println("Starting up on port " + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		http.ListenAndServe(":"+port, nil)
	}
}

func index(w http.ResponseWriter, req *http.Request) {
	u, _ := url.Parse(req.URL.String())
	queryParams := u.Query()

	wait := queryParams.Get("wait")
	if len(wait) > 0 {
		duration, err := time.ParseDuration(wait)
		if err == nil {
			time.Sleep(duration)
		}
	}

	data := fetchData(req)
	fmt.Fprintln(w, "This is some xkpfheyc test string")
	fmt.Fprintln(w, "Hostname:", data.Hostname)

	for _, ip := range data.IP {
		fmt.Fprintln(w, "IP:", ip)
	}

	for _, env := range data.Environment {
		fmt.Fprintln(w, "ENV:", env)
	}
	req.Write(w)
}

func api(w http.ResponseWriter, req *http.Request) {
	data := fetchData(req)
	json.NewEncoder(w).Encode(data)
}

func fetchData(req *http.Request) DataResponse {
	hostname, _ := os.Hostname()
	data := DataResponse{
		hostname,
		[]string{},
		req.Header,
		os.Environ(),
	}

	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			data.IP = append(data.IP, ip.String())
		}
	}

	return data
}
