package shortener

import (
	"context"
	"net/url"
	"strings"

	"github.com/Lind-32/urlshortenergrpc/internal/pkg/store"
	api "github.com/Lind-32/urlshortenergrpc/pkg"
)

// Retrive получает короткую ссылку, возвращает длинную
func (s *Server) Retrive(ctx context.Context, req *api.ShortLinkRequest) (*api.LongLinkResponse, error) {

	url, err := url.Parse(req.GetShortlink())
	if err != nil {
		panic(err)
	}

	key := strings.TrimPrefix(url.Path, "/to/")

	res, err := store.GetLongURL(key)
	if err != nil {
		panic(err)
	}
	if res == "" {
		res = "not saved"
	}

	return &api.LongLinkResponse{Longlink: "Retrieved long link: " + res}, nil
}
