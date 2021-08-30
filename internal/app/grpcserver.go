package shortener

import (
	"log"
	"net"

	"github.com/Lind-32/urlshortenergrpc/internal/config"
	api "github.com/Lind-32/urlshortenergrpc/pkg"
	"google.golang.org/grpc"
)

//Server gRPC
type Server struct {
	api.UnimplementedShortLinkServer
}

//ServerGRPC описание GRPC сервера
func ServerGRPC(cfg config.Config) {

	lis, err := net.Listen("tcp", cfg.Grpcadr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	} else {
		log.Println("GRPC server is running on", cfg.Grpcadr)
	}

	s := grpc.NewServer()
	api.RegisterShortLinkServer(s, &Server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve grpc: %v", err)
	}
}
