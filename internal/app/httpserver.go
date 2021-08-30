package shortener

import (
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"

	"github.com/Lind-32/urlshortenergrpc/internal/config"
	"github.com/Lind-32/urlshortenergrpc/internal/pkg/store"

	"github.com/gorilla/mux"
)

//ServerHTTP описание сервера webpage
func ServerHTTP(cfg config.Config) {

	r := mux.NewRouter()

	r.HandleFunc("/", homepage)
	r.HandleFunc("/to/{key}", redirect)
	http.Handle("/", r)

	lis, err := net.Listen("tcp", cfg.Httpadr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	} else {
		log.Println("HTTP server is running on", cfg.Httpadr)
	}
	s := &http.Server{
		Handler: r,
	}
	if err := s.Serve((lis)); err != nil {
		log.Fatalf("failed to serve http: %v", err)
	}
}

// редирект на сохраненную страницу по ключу
func redirect(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	key := vars["key"]
	l, err := store.GetLongURL(key)
	if err != nil {
		panic(err)
	}
	http.Redirect(w, r, l, http.StatusSeeOther)
}

// интерфейс базовой страницы
func homepage(w http.ResponseWriter, r *http.Request) {

	var link string

	temp, err := template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Fprintf(w, "Template error: %s/n", err.Error()) //вывод шаблона на homepage
	}

	if r.Method == "POST" {

		link = r.FormValue("link")

		if !ValidURL(link) { //проверка введенной ссылки на наличие хоста и схемы
			link = "invalid link format"
		} else {
			key, unic, err := store.UnicURL(link) //проверка длинной ссылки на уникальность, если есть в базе, возвращает короткую
			if err != nil {
				panic(err)
			}
			if !unic {

				link = "http://" + config.GetConfig().Httpadr + "/to/" + key //вывод сохраненной ссылки, если таковая найдена в базе
			} else {
				sh := short()

				shortlink := "http://" + config.GetConfig().Httpadr + "/to/" + sh //вывод сгенерированной ссылки

				err = store.Insert(sh, link) //запись в БД сгенерированного ключа и длинной ссылки
				if err != nil {
					panic(err)
				}
				link = shortlink
			}

		}
	}
	temp.Execute(w, link)
}
