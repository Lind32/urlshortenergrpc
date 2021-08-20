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

	"github.com/Lind-32/urlshortenergrpc/api"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

//Data ... in-memory DB
type Data struct {
	db map[string]string
}

var data = &Data{db: make(map[string]string)}

//Result вывод на страницу
type Result struct {
	Link string
}

var result Result

//Server gRPC
type Server struct {
	api.UnimplementedShortLinkServer
}

// адрес GRPC сервера
const grpcaddress = "localhost:8000"
const httpaddress = "localhost:8080"

func main() {
	//ServerHTTP()
	go ServerHTTP()
	ServerGRPC()
}

//ServerHTTP описание сервера webpage
func ServerHTTP() {

	r := mux.NewRouter()

	r.HandleFunc("/", data.homepage)
	r.HandleFunc("/to/{key}", data.redirect)
	http.Handle("/", r)

	lis, err := net.Listen("tcp", httpaddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	} else {
		fmt.Println("HTTP server is running on", httpaddress)
	}
	server := &http.Server{
		Handler: r,
	}
	if err := server.Serve((lis)); err != nil {
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

	result.Link = req.GetLonglink()
	if !ValidURL(result.Link) {
		result.Link = "invalid link format"
	} else {
		sh := short()
		shortlink := "http://" + httpaddress + "/to/" + sh
		data.db[sh] = result.Link
		result.Link = shortlink
	}

	return &api.ShortLinkResponse{Shortlink: "Generated short link: " + result.Link}, nil
}

// Retrive получает короткую ссылку, возвращает длинную
func (s *Server) Retrive(ctx context.Context, req *api.ShortLinkRequest) (*api.LongLinkResponse, error) {
	url, err := url.Parse(req.GetShortlink())
	if err != nil {
		panic(err)
	}

	key := strings.Trim(url.Path, "/to/")
	res := data.db[key]
	if res == "" {
		res = "not saved"
	}

	return &api.LongLinkResponse{Longlink: "Retrieved lond link: " + res}, nil
}

//генерация ключа
func short() string {

	var letters string = "_QWERTYUIOPASDFGHJKLZXCVBNMqwertyuiopasdfghjklzxcvbnm1234567890"
	var key string
	randLetters := make([]byte, 10)

	for !ValidKey(key) {
		for i := range randLetters {
			randLetters[i] = letters[rand.Intn(len(letters))]
		}
		key = string(randLetters)
	}
	return key

}

//ValidKey проверка ключа на уникальность
func ValidKey(key string) bool {
	if key == "" {
		return false
	}
	for keydb := range data.db {
		if keydb == key {
			return false
		}
	}
	return true

}

//UnicURL проверка длинной ссылки на уникальность
func UnicURL(url string) (string, bool) {

	for k, u := range data.db {
		if url == u {
			return k, false
		}
	}
	return "", true
}

// редирект на сохраненную страницу по ключу
func (data *Data) redirect(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	key := vars["key"]
	http.Redirect(w, r, data.db[key], http.StatusSeeOther)
	//	w.WriteHeader(http.StatusOK)

}

// интерфейс базовой страницы
func (data *Data) homepage(w http.ResponseWriter, r *http.Request) {

	temp, err := template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Fprintf(w, "Template error: %s/n", err.Error())
	}

	if r.Method == "POST" {

		result.Link = r.FormValue("link")

		if !ValidURL(result.Link) {
			result.Link = "invalid link format"
		} else {
			k, unic := UnicURL(result.Link)
			if !unic {
				result.Link = "http://" + httpaddress + "/to/" + k
			} else {
				sh := short()
				shortlink := "http://" + httpaddress + "/to/" + sh
				data.db[sh] = result.Link
				result.Link = shortlink
			}
			for key, value := range data.db {
				fmt.Printf("%s === %s \n", key, value)
			}
			fmt.Println("____________________________________")
		}
	}
	temp.Execute(w, result)
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
