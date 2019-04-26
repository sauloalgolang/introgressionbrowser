package web

// https://medium.com/@adigunhammedolalekan/build-and-deploy-a-secure-rest-api-with-go-postgresql-jwt-and-gorm-6fadf3da505b

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

import (
	"github.com/sauloalgolang/introgressionbrowser/web/endpoints"
)

const HTTP_ROOT_DIR = "http"

func NewWeb(databaseDir string, host string, port int, verbosityLevel log.Level) {
	router := mux.NewRouter()

	log.Warn("open your browser at http://" + host + ":" + strconv.Itoa(port))

	api := router.PathPrefix("/api/").Subrouter()

	newStatic(databaseDir, router)
	newApi(databaseDir, api, verbosityLevel)
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
		log.Panic(err)
	}
}

func newRoot(dir string, router *mux.Router) {
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(dir))))
}

func newStatic(dir string, router *mux.Router) {
	router.PathPrefix("/database/").Handler(http.StripPrefix("/database/", http.FileServer(http.Dir(dir))))
}

func newApi(dir string, router *mux.Router, verbosityLevel log.Level) {
	router.HandleFunc("/databases", endpoints.Databases).Methods("GET")                                                                        //.HeadersRegexp("Content-Type", "application/json")
	router.HandleFunc("/databases/{database}/block", endpoints.DatabaseBlock).Methods("GET")                                                   //.HeadersRegexp("Content-Type", "application/json")
	router.HandleFunc("/databases/{database}/chromosomes", endpoints.Chromosomes).Methods("GET")                                               //.HeadersRegexp("Content-Type", "application/json")
	router.HandleFunc("/databases/{database}/chromosomes/{chromosome}/block", endpoints.ChromosomeBlock).Methods("GET")                        //.HeadersRegexp("Content-Type", "application/json")
	router.HandleFunc("/databases/{database}/chromosomes/{chromosome}/blocks", endpoints.Blocks).Methods("GET")                                //.HeadersRegexp("Content-Type", "application/json")
	router.HandleFunc("/databases/{database}/chromosomes/{chromosome}/blocks/{blockNum:[0-9]+}/block", endpoints.BlocksBlock).Methods("GET")   //.HeadersRegexp("Content-Type", "application/json")
	router.HandleFunc("/databases/{database}/chromosomes/{chromosome}/blocks/{blockNum:[0-9]+}/matrix", endpoints.BlocksMatrix).Methods("GET") //.HeadersRegexp("Content-Type", "application/json")

	endpoints.DATABASE_DIR = dir
	endpoints.VERBOSITY = verbosityLevel

	endpoints.ListDatabases()
}
