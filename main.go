package main

import (
	"errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

// https://cloud.tencent.com/document/product/1416/56033
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

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})
)

func recordMetrics() {
	go func() {
		for {
			opsProcessed.Inc()
			time.Sleep(2 * time.Second)
		}
	}()
}

func main() {
	recordMetrics()
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
