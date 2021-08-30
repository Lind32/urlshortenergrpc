package shortener

import (
	"context"

	"github.com/Lind-32/urlshortenergrpc/internal/config"
	"github.com/Lind-32/urlshortenergrpc/internal/pkg/store"

	api "github.com/Lind-32/urlshortenergrpc/pkg"
)

//Generate получает длинную ссылку, возвращает короткую
func (s *Server) Generate(ctx context.Context, req *api.LongLinkRequest) (*api.ShortLinkResponse, error) {

	link := req.GetLonglink()
	if !ValidURL(link) { //проверка введенной ссылки на наличие хоста и схемы
		link = "invalid link format"
	} else {
		key, unic, err := store.UnicURL(link) //проверка длинной ссылки на уникальность, если есть в базе, возвращает короткую
		if err != nil {
			panic(err)
		}
		if !unic {
			link = "http://" + config.GetConfig().Httpadr + "/to/" + key
		} else {

			sh := short()
			shortlink := "http://" + config.GetConfig().Httpadr + "/to/" + sh

			err = store.Insert(sh, link) //запись в БД
			if err != nil {
				panic(err)
			}
			link = shortlink
		}
	}

	return &api.ShortLinkResponse{Shortlink: "Generated short link: " + link}, nil
}
