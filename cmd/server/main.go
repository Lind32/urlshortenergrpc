package main

import (
	"github.com/Lind-32/urlshortenergrpc/internal/app"
	_ "github.com/lib/pq"
)

func main() {

	//запуск gRPC и HTTP серверов
	go app.ServerHTTP()
	app.ServerGRPC()
}
