package app

import (
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"

	"github.com/Lind-32/urlshortenergrpc/internal/pkg/store"
	"github.com/gorilla/mux"
)

//ServerHTTP описание сервера webpage
func ServerHTTP() {

	r := mux.NewRouter()

	r.HandleFunc("/", homepage)
	r.HandleFunc("/to/{key}", redirect)
	http.Handle("/", r)

	lis, err := net.Listen("tcp", httpadr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	} else {
		fmt.Println("HTTP server is running on", httpadr)
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

	err := store.Connect()
	if err != nil {
		panic(err)
	}
	defer store.Close()

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
	err := store.Connect()
	if err != nil {
		panic(err)
	}
	defer store.Close()

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
				link = "http://" + httpadr + "/to/" + key //вывод сохраненной ссылки, если таковая найдена в базе
			} else {
				sh := short()

				shortlink := "http://" + httpadr + "/to/" + sh //вывод сгенерированной ссылки

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
