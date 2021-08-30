package shortener

import (
	"context"
	"testing"

	_ "github.com/Lind-32/urlshortenergrpc/db"
	"github.com/Lind-32/urlshortenergrpc/internal/config"
	"github.com/Lind-32/urlshortenergrpc/internal/pkg/store"

	api "github.com/Lind-32/urlshortenergrpc/pkg"
)

func TestGenerate(t *testing.T) {
	s := Server{}
	llink := "https://www.5gI3SC6q.test/"
	key := "testtest"

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
	slink := "Generated short link: http://" + config.GetConfig().Httpadr + "/to/" + key

	res, err := s.Generate(context.Background(), &api.LongLinkRequest{Longlink: llink})
	if err != nil {
		t.Errorf("Generate failed: %v", err)
	}
	if res.GetShortlink() != slink {
		t.Errorf("Generate failed: %v != %v", res.GetShortlink(), slink)
	}

}
