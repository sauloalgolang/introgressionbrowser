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
const DATA_ENDPOINT = "/data"
const DATABASE_ENDPOINT = "/database"

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
	router.PathPrefix(DATA_ENDPOINT).Handler(http.StripPrefix(DATA_ENDPOINT+"/", http.FileServer(http.Dir(dir))))
}

func newApi(dir string, router *mux.Router, verbosityLevel log.Level) {
	router.HandleFunc("update", endpoints.Update).Methods("POST")                                                                                               //.HeadersRegexp("Content-Type", "application/json")
	router.HandleFunc(DATABASE_ENDPOINT, endpoints.Databases).Methods("GET")                                                                                    //.HeadersRegexp("Content-Type", "application/json")
	router.HandleFunc(DATABASE_ENDPOINT+"/{database}", endpoints.Database).Methods("GET")                                                                       //.HeadersRegexp("Content-Type", "application/json")
	router.HandleFunc(DATABASE_ENDPOINT+"/{database}/summary", endpoints.DatabaseSummary).Methods("GET")                                                        //.HeadersRegexp("Content-Type", "application/json")
	router.HandleFunc(DATABASE_ENDPOINT+"/{database}/summary/matrix", endpoints.DatabaseSummaryMatrix).Methods("GET")                                           //.HeadersRegexp("Content-Type", "application/json")
	router.HandleFunc(DATABASE_ENDPOINT+"/{database}/summary/matrix/table", endpoints.DatabaseSummaryMatrixTable).Methods("GET")                                //.HeadersRegexp("Content-Type", "application/json")
	router.HandleFunc(DATABASE_ENDPOINT+"/{database}/chromosome", endpoints.Chromosomes).Methods("GET")                                                         //.HeadersRegexp("Content-Type", "application/json")
	router.HandleFunc(DATABASE_ENDPOINT+"/{database}/chromosome/{chromosome}", endpoints.Chromosome).Methods("GET")                                             //.HeadersRegexp("Content-Type", "application/json")
	router.HandleFunc(DATABASE_ENDPOINT+"/{database}/chromosome/{chromosome}/summary", endpoints.ChromosomeSummary).Methods("GET")                              //.HeadersRegexp("Content-Type", "application/json")
	router.HandleFunc(DATABASE_ENDPOINT+"/{database}/chromosome/{chromosome}/summary/matrix", endpoints.ChromosomeSummaryMatrix).Methods("GET")                 //.HeadersRegexp("Content-Type", "application/json")
	router.HandleFunc(DATABASE_ENDPOINT+"/{database}/chromosome/{chromosome}/summary/table", endpoints.ChromosomeSummaryMatrixTable).Methods("GET")             //.HeadersRegexp("Content-Type", "application/json")
	router.HandleFunc(DATABASE_ENDPOINT+"/{database}/chromosome/{chromosome}/block", endpoints.Blocks).Methods("GET")                                           //.HeadersRegexp("Content-Type", "application/json")
	router.HandleFunc(DATABASE_ENDPOINT+"/{database}/chromosome/{chromosome}/block/{blockNum:[0-9]+}", endpoints.Block).Methods("GET")                          //.HeadersRegexp("Content-Type", "application/json")
	router.HandleFunc(DATABASE_ENDPOINT+"/{database}/chromosome/{chromosome}/block/{blockNum:[0-9]+}/matrix", endpoints.BlockMatrix).Methods("GET")             //.HeadersRegexp("Content-Type", "application/json")
	router.HandleFunc(DATABASE_ENDPOINT+"/{database}/chromosome/{chromosome}/block/{blockNum:[0-9]+}/matrix/table", endpoints.BlocksMatrixTable).Methods("GET") //.HeadersRegexp("Content-Type", "application/json")

	endpoints.DATABASE_DIR = dir
	endpoints.VERBOSITY = verbosityLevel

	endpoints.ListDatabases()
}
