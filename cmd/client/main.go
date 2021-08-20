package main

import (
	"context"
	"flag"
	"log"

	"github.com/Lind-32/urlshortenergrpc/api"
	"google.golang.org/grpc"
)

const (
	address = "localhost:8000"
)

func main() {

	// описание подключения к серверу
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// определение аргументов
	longlink := flag.String("link", "link", "link=http://google.com")
	generateflag := flag.Bool("g", false, "generate short link")
	retriveflag := flag.Bool("r", false, "retrive long link")

	flag.Parse()

	if flag.NFlag() < 2 {
		log.Fatal("not enough arguments")
	}
	if *longlink == "link" {
		log.Fatal("link is missing")
	}
	//Generate получает длинную ссылку, возвращает короткую
	c := api.NewShortLinkClient(conn)
	if *generateflag {
		res, err := c.Generate(context.Background(), &api.LongLinkRequest{Longlink: *longlink})
		if err != nil {
			log.Fatal(err)
		}
		log.Println(res.GetShortlink())
		return
	}
	// Retrive получает короткую ссылку, возвращает длинную
	if *retriveflag {
		res, err := c.Retrive(context.Background(), &api.ShortLinkRequest{Shortlink: *longlink})
		if err != nil {
			log.Fatal(err)
		}
		log.Println(res.GetLonglink())
		return
	}
}
