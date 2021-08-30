package main

import (
	"context"
	"flag"
	"log"

	"github.com/Lind-32/urlshortenergrpc/internal/config"
	api "github.com/Lind-32/urlshortenergrpc/pkg"
	"google.golang.org/grpc"
)

func main() {

	// описание подключения к серверу
	conn, err := grpc.Dial(config.GetConfig().Grpcadr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// определение аргументов
	link := flag.String("link", "link", "link=http://google.com")
	generateflag := flag.Bool("g", false, "generate short link")
	retriveflag := flag.Bool("r", false, "retrive long link")

	flag.Parse()

	if flag.NFlag() < 2 {
		log.Fatal("not enough arguments")
	}
	if *link == "link" {
		log.Fatal("link is missing")
	}

	c := api.NewShortLinkClient(conn)

	//Generate получает длинную ссылку, возвращает короткую
	if *generateflag {
		res, err := c.Generate(context.Background(), &api.LongLinkRequest{Longlink: *link})
		if err != nil {
			log.Fatal(err)
		}
		log.Println(res.GetShortlink())
		return
	}
	// Retrive получает короткую ссылку, возвращает длинную
	if *retriveflag {
		res, err := c.Retrive(context.Background(), &api.ShortLinkRequest{Shortlink: *link})
		if err != nil {
			log.Fatal(err)
		}
		log.Println(res.GetLonglink())
		return
	}
}
