package main

import (
	"log"
	"net/http"
)

func main() {
	// 写一个启动 http 服务代码
	http.HandleFunc("/traefik", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("start traefik v2"))
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
