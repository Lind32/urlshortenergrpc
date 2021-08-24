package main

import (
	"github.com/Lind-32/urlshortenergrpc/internal/app"
)

func main() {

	//запуск gRPC и HTTP серверов
	go app.ServerGRPC()
	app.ServerGRPC()
}
