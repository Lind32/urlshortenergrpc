package app

import (
	"math/rand"
	"net/url"

	"github.com/Lind-32/urlshortenergrpc/internal/pkg/store"
)

//ValidURL проверка URL адреса
func ValidURL(tocen string) bool {
	_, err := url.ParseRequestURI(tocen)
	if err != nil {
		return false
	}
	u, err := url.Parse(tocen)
	if err != nil || u.Host == "" {
		return false
	}
	return true
}

//генерация ключа
func short() string {

	var letters string = "_QWERTYUIOPASDFGHJKLZXCVBNMqwertyuiopasdfghjklzxcvbnm1234567890"
	var key string
	var err error
	randLetters := make([]byte, 10)
	u := false

	for !u {
		for i := range randLetters {
			randLetters[i] = letters[rand.Intn(len(letters))]
		}
		key = string(randLetters)
		u, err = store.UnicKey(key) //проверка ключа на уникальность (true если уникальный)
		if err != nil {
			panic(err)
		}
	}
	return key
}
