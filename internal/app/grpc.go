package app

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/url"
	"strings"

	"github.com/Lind-32/urlshortenergrpc/internal/pkg/store"
	api "github.com/Lind-32/urlshortenergrpc/pkg"
	"google.golang.org/grpc"
)

const (
	grpcadr = "localhost:8020"
	httpadr = "localhost:8080"
)

//Server gRPC
type Server struct {
	api.UnimplementedShortLinkServer
}

//ServerGRPC описание GRPC сервера
func ServerGRPC() {

	lis, err := net.Listen("tcp", grpcadr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	} else {
		fmt.Println("GRPC server is running on", grpcadr)
	}

	s := grpc.NewServer()
	api.RegisterShortLinkServer(s, &Server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve grpc: %v", err)
	}
}

//GRPC методы:

//Generate получает длинную ссылку, возвращает короткую
func (s *Server) Generate(ctx context.Context, req *api.LongLinkRequest) (*api.ShortLinkResponse, error) {

	err := store.Connect()
	if err != nil {
		panic(err)
	}
	defer store.Close()

	link := req.GetLonglink()
	if !ValidURL(link) { //проверка введенной ссылки на наличие хоста и схемы
		link = "invalid link format"
	} else {
		key, unic, err := store.UnicURL(link) //проверка длинной ссылки на уникальность, если есть в базе, возвращает короткую
		if err != nil {
			panic(err)
		}
		if !unic {
			link = "http://" + httpadr + "/to/" + key
		} else {

			sh := short()
			shortlink := "http://" + httpadr + "/to/" + sh
			err = store.Insert(sh, link) //запись в БД
			if err != nil {
				panic(err)
			}
			link = shortlink
		}
	}

	return &api.ShortLinkResponse{Shortlink: "Generated short link: " + link}, nil
}

// Retrive получает короткую ссылку, возвращает длинную
func (s *Server) Retrive(ctx context.Context, req *api.ShortLinkRequest) (*api.LongLinkResponse, error) {

	err := store.Connect()
	if err != nil {
		panic(err)
	}
	defer store.Close()

	url, err := url.Parse(req.GetShortlink())
	if err != nil {
		panic(err)
	}
	key := strings.Trim(url.Path, "to/")
	res, err := store.GetLongURL(key)
	if err != nil {
		panic(err)
	}
	if res == "" {
		res = "not saved"
	}

	return &api.LongLinkResponse{Longlink: "Retrieved lond link: " + res}, nil
}
