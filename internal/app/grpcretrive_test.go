package shortener

import (
	"context"
	"testing"

	_ "github.com/Lind-32/urlshortenergrpc/db"
	"github.com/Lind-32/urlshortenergrpc/internal/config"
	"github.com/Lind-32/urlshortenergrpc/internal/pkg/store"

	api "github.com/Lind-32/urlshortenergrpc/pkg"
)

func TestRetrive(t *testing.T) {
	s := Server{}
	llink := "https://www.5gI3SC6q.test/"
	key := "test__test"

	//подключение к БД
	err := store.Connect(*config.GetConfig())
	if err != nil {
		panic(err)
	}
	defer store.Close()

	//UnicURL проверка длинной ссылки на уникальность (true если уникальна)
	_, unic, err := store.UnicURL(llink)
	if err != nil {
		panic(err)
	}
	if unic {
		err = store.Insert(key, llink) //запись в БД
		if err != nil {
			panic(err)
		}
	}
	//тестирование
	slink := "http://" + config.GetConfig().Httpadr + "/to/" + key

	res, err := s.Retrive(context.Background(), &api.ShortLinkRequest{Shortlink: slink})
	if err != nil {
		t.Errorf("Retrive failed: %v", err)
	}
	if res.GetLonglink() != "Retrieved long link: "+llink {
		t.Errorf("Retrive failed: %v != %v", res.GetLonglink(), llink)
	}

}
