package main

import (
	"errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net"
	"net/http"
	"strings"
)

// GetIP returns request real ip.
func GetIP(r *http.Request) (string, error) {
	ip := r.Header.Get("X-Real-IP")
	if net.ParseIP(ip) != nil {
		return ip, nil
	}
	ip = r.Header.Get("X-Forward-For")
	for _, i := range strings.Split(ip, ",") {
		if net.ParseIP(i) != nil {
			return i, nil
		}
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	if net.ParseIP(ip) != nil {
		return ip, nil
	}
	return "", errors.New("no valid ip found")
}

func main() {
	// 写一个启动 http 服务代码
	http.HandleFunc("/traefik", func(writer http.ResponseWriter, request *http.Request) {
		log.Println("start traefik v2")
		ip, _ := GetIP(request)
		log.Printf("ip:%s", ip)
		writer.Write([]byte("start traefik v2"))
	})

	http.Handle("/metrics", promhttp.Handler())

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
