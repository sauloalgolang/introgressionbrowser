package web

// https://medium.com/@adigunhammedolalekan/build-and-deploy-a-secure-rest-api-with-go-postgresql-jwt-and-gorm-6fadf3da505b

import (
	// "os"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

import (
	"github.com/sauloalgolang/introgressionbrowser/web/endpoints"
)

const HTTP_ROOT_DIR = "http"

func NewWeb(staticDir string, host string, port int) {
	router := mux.NewRouter()

	fmt.Println("open your browser at http://" + host + ":" + strconv.Itoa(port))

	api := router.PathPrefix("/api").Subrouter()

	newStatic(staticDir, router)
	newApi(api)
	newRoot(HTTP_ROOT_DIR, router)

	srv := &http.Server{
		Handler: router,
		Addr:    host + ":" + strconv.Itoa(port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	err := srv.ListenAndServe() //Launch the app, visit localhost:8000/api

	if err != nil {
		fmt.Print(err)
	}
}

func newRoot(dir string, router *mux.Router) {
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(dir))))
}

func newStatic(dir string, router *mux.Router) {
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))
}

func newApi(router *mux.Router) {
	router.HandleFunc("/databases", endpoints.Databases).Methods("GET").HeadersRegexp("Content-Type", "application/json")
	router.HandleFunc("/databases/{database}/", endpoints.Database).Methods("GET").HeadersRegexp("Content-Type", "application/json")
	router.HandleFunc("/databases/{database}/chromosomes", endpoints.Chromosomes).Methods("GET").HeadersRegexp("Content-Type", "application/json")
	router.HandleFunc("/databases/{database}/chromosomes/{chromosome}", endpoints.Chromosome).Methods("GET").HeadersRegexp("Content-Type", "application/json")
	router.HandleFunc("/databases/{database}/chromosomes/{chromosome}/matrices", endpoints.Matrices).Methods("GET").HeadersRegexp("Content-Type", "application/json")
	router.HandleFunc("/databases/{database}/chromosomes/{chromosome}/matrices/{matrix}", endpoints.Matrix).Methods("GET").HeadersRegexp("Content-Type", "application/json")
}
