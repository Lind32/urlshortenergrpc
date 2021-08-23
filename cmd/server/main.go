package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"
	"text/template"

	api "github.com/Lind-32/urlshortenergrpc/internal/pkg"
	"github.com/Lind-32/urlshortenergrpc/internal/pkg/store"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

//Server gRPC
type Server struct {
	api.UnimplementedShortLinkServer
}

const (
	grpcaddress = "localhost:8020"
	httpaddress = "localhost:8080"
)

var link string

func main() {

	go ServerHTTP() //запуск gRPC и HTTP серверов
	ServerGRPC()
}

//ServerHTTP описание сервера webpage
func ServerHTTP() {

	r := mux.NewRouter()

	r.HandleFunc("/", homepage)
	r.HandleFunc("/to/{key}", redirect)
	http.Handle("/", r)

	lis, err := net.Listen("tcp", httpaddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	} else {
		fmt.Println("HTTP server is running on", httpaddress)
	}
	s := &http.Server{
		Handler: r,
	}
	if err := s.Serve((lis)); err != nil {
		log.Fatalf("failed to serve http: %v", err)
	}
}

//ServerGRPC описание GRPC сервера
func ServerGRPC() {

	lis, err := net.Listen("tcp", grpcaddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	} else {
		fmt.Println("GRPC server is running on", grpcaddress)
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

	link = req.GetLonglink()
	if !ValidURL(link) { //проверка введенной ссылки на наличие хоста и схемы
		link = "invalid link format"
	} else {
		key, unic, err := store.UnicURL(link) //проверка длинной ссылки на уникальность, если есть в базе, возвращает короткую
		if err != nil {
			panic(err)
		}
		if !unic {
			link = "http://" + httpaddress + "/to/" + key
		} else {

			sh := short()
			shortlink := "http://" + httpaddress + "/to/" + sh
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
	key := strings.Trim(url.Path, "/to/")
	res, err := store.GetLongURL(key)
	if err != nil {
		panic(err)
	}
	if res == "" {
		res = "not saved"
	}

	return &api.LongLinkResponse{Longlink: "Retrieved lond link: " + res}, nil
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
				link = "http://" + httpaddress + "/to/" + key //вывод сохраненной ссылки, если таковая найдена в базе
			} else {
				sh := short()

				shortlink := "http://" + httpaddress + "/to/" + sh //вывод сгенерированной ссылки

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
