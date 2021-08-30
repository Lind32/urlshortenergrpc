package main

import (
	_ "github.com/Lind-32/urlshortenergrpc/db"
	shortener "github.com/Lind-32/urlshortenergrpc/internal/app"
	"github.com/Lind-32/urlshortenergrpc/internal/config"
	"github.com/Lind-32/urlshortenergrpc/internal/pkg/store"
)

func main() {

	cfg := config.GetConfig()

	// подключение к базе данных
	err := store.Connect(*cfg)
	if err != nil {
		panic(err)
	}
	defer store.Close()

	//запуск gRPC и HTTP серверов
	go shortener.ServerGRPC(*cfg)
	shortener.ServerHTTP(*cfg)

}
