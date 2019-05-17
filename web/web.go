package web

// https://medium.com/@adigunhammedolalekan/build-and-deploy-a-secure-rest-api-with-go-postgresql-jwt-and-gorm-6fadf3da505b

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
	"time"
)

import (
	"github.com/sauloalgolang/introgressionbrowser/web/endpoints"
)

// APIEndpoints is the web api endpoint
const APIEndpoints = "/api"

// DataEndpoint is the web data endpoint
const DataEndpoint = "/data"

// DatabaseEndpoint is the web database endpoint
const DatabaseEndpoint = "/databases"

// PlotsEndpoint is the web plots endpoint
const PlotsEndpoint = "/plots"

// NewWeb starts a new webserver
func NewWeb(databaseDir string, httpDir string, host string, port int) {
	router := mux.NewRouter()

	router.StrictSlash(false)
	// endpoints.SetRouter(router)

	log.Warn("open your browser at http://" + host + ":" + strconv.Itoa(port))

	api := router.PathPrefix(APIEndpoints).Subrouter()
	api.StrictSlash(false)
	router.HandleFunc(APIEndpoints+"/", template).Methods("GET").Name("apiSlash")

	newData(databaseDir, router)
	newAPI(databaseDir, api)
	newRoot(httpDir, router)

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
	router.PathPrefix("").Handler(http.StripPrefix("/", http.FileServer(http.Dir(dir)))).Name("root")
}

func newData(dir string, router *mux.Router) {
	router.PathPrefix(DataEndpoint + "/").Handler(http.StripPrefix(DataEndpoint+"/", http.FileServer(http.Dir(dir)))).Name("dataSlash")
	router.PathPrefix(DataEndpoint).Handler(http.StripPrefix(DataEndpoint, http.FileServer(http.Dir(dir)))).Name("data")
}

func newAPI(dir string, router *mux.Router) {
	//.HeadersRegexp("Content-Type", "application/json")
	router.HandleFunc("", template).Methods("GET").Name("api")
	router.HandleFunc("/update", endpoints.Update).Methods("POST").Name("update")
	router.HandleFunc(DatabaseEndpoint, endpoints.Databases).Methods("GET").Name("databases")
	router.HandleFunc(DatabaseEndpoint+"/{database}", endpoints.Database).Methods("GET").Name("database")
	router.HandleFunc(DatabaseEndpoint+"/{database}/summary", endpoints.DatabaseSummary).Methods("GET").Name("databaseSummary")
	router.HandleFunc(DatabaseEndpoint+"/{database}/summary/matrix", endpoints.DatabaseSummaryMatrix).Methods("GET").Name("databaseSummaryMatrix")
	router.HandleFunc(DatabaseEndpoint+"/{database}/summary/matrix/table", endpoints.DatabaseSummaryMatrixTable).Methods("GET").Name("databaseSummaryMatrixTable")
	router.HandleFunc(DatabaseEndpoint+"/{database}/chromosome", endpoints.Chromosomes).Methods("GET").Name("databaseChromosomes")
	router.HandleFunc(DatabaseEndpoint+"/{database}/chromosomes/{chromosome}", endpoints.Chromosome).Methods("GET").Name("databaseChromosome")
	router.HandleFunc(DatabaseEndpoint+"/{database}/chromosomes/{chromosome}/summary", endpoints.ChromosomeSummary).Methods("GET").Name("databaseChromosomeSummary")
	router.HandleFunc(DatabaseEndpoint+"/{database}/chromosomes/{chromosome}/summary/matrix", endpoints.ChromosomeSummaryMatrix).Methods("GET").Name("databaseChromosomeSummaryMatrix")
	router.HandleFunc(DatabaseEndpoint+"/{database}/chromosomes/{chromosome}/summary/matrix/table", endpoints.ChromosomeSummaryMatrixTable).Methods("GET").Name("databaseChromosomeSummaryTable")
	router.HandleFunc(DatabaseEndpoint+"/{database}/chromosomes/{chromosome}/block", endpoints.Blocks).Methods("GET").Name("databaseChromosomeBlocks")
	router.HandleFunc(DatabaseEndpoint+"/{database}/chromosomes/{chromosome}/blocks/{blockNum:[0-9]+}", endpoints.Block).Methods("GET").Name("databaseChromosomeBlock")
	router.HandleFunc(DatabaseEndpoint+"/{database}/chromosomes/{chromosome}/blocks/{blockNum:[0-9]+}/matrix", endpoints.BlockMatrix).Methods("GET").Name("databaseChromosomeBlockMatrix")
	router.HandleFunc(DatabaseEndpoint+"/{database}/chromosomes/{chromosome}/blocks/{blockNum:[0-9]+}/matrix/table", endpoints.BlocksMatrixTable).Methods("GET").Name("databaseChromosomeBlockMatrixTable")

	router.HandleFunc(PlotsEndpoint+"/{database}/{chromosome}/{referenceName}", endpoints.Plots).Methods("GET").Name("plots")

	route := router.Get("data")

	if route == nil {
		log.Panic("No data router")
	}

	tmpl, tErr := route.GetPathTemplate()
	if tErr != nil {
		log.Panic("No data router template")
	}

	tmpl = strings.TrimSuffix(tmpl, "/")

	endpoints.DataEndpoint = tmpl
	endpoints.DatabaseDir = strings.TrimSuffix(dir, "/")

	endpoints.ListDatabases()
}

func template(w http.ResponseWriter, r *http.Request) {
	log.Tracef("Template %#v", r)

	resp := endpoints.Message(true, "success")

	tmp := map[string]interface{}{
		APIEndpoints + DatabaseEndpoint + "":                                                                           []string{""},
		APIEndpoints + DatabaseEndpoint + "/{database}":                                                                endpoints.DatabaseInfo{},
		APIEndpoints + DatabaseEndpoint + "/{database}/summary":                                                        endpoints.BlockInfo{},
		APIEndpoints + DatabaseEndpoint + "/{database}/summary/matrix":                                                 endpoints.MatrixInfo{},
		APIEndpoints + DatabaseEndpoint + "/{database}/summary/matrix/table":                                           endpoints.TableInfo{},
		APIEndpoints + DatabaseEndpoint + "/{database}/chromosome":                                                     []endpoints.ChromosomeInfo{endpoints.ChromosomeInfo{}},
		APIEndpoints + DatabaseEndpoint + "/{database}/chromosomes/{chromosome}":                                       endpoints.ChromosomeInfo{},
		APIEndpoints + DatabaseEndpoint + "/{database}/chromosomes/{chromosome}/summary":                               endpoints.BlockInfo{},
		APIEndpoints + DatabaseEndpoint + "/{database}/chromosomes/{chromosome}/summary/matrix":                        endpoints.MatrixInfo{},
		APIEndpoints + DatabaseEndpoint + "/{database}/chromosomes/{chromosome}/summary/matrix/table":                  endpoints.TableInfo{},
		APIEndpoints + DatabaseEndpoint + "/{database}/chromosomes/{chromosome}/block":                                 []endpoints.BlockInfo{endpoints.BlockInfo{}},
		APIEndpoints + DatabaseEndpoint + "/{database}/chromosomes/{chromosome}/blocks/{blockNum:[0-9]+}":              endpoints.BlockInfo{},
		APIEndpoints + DatabaseEndpoint + "/{database}/chromosomes/{chromosome}/blocks/{blockNum:[0-9]+}/matrix":       endpoints.MatrixInfo{},
		APIEndpoints + DatabaseEndpoint + "/{database}/chromosomes/{chromosome}/blocks/{blockNum:[0-9]+}/matrix/table": endpoints.TableInfo{},
	}

	resp["data"] = tmp

	endpoints.Respond(w, resp)
}

//
//TODO
//

//
// Query parameters
//
// https://stackoverflow.com/questions/45378566/gorilla-mux-optional-query-values
//
// router.Path("/articles/{id:[0-9]+}").
//     Queries("key", "{[0-9]*?}").
//     HandlerFunc(YourHandler).
//     Name("YourHandler")
// router.Path("/articles/{id:[0-9]+}").HandlerFunc(YourHandler)
//
// router.Path("/articles/{id:[0-9]+}").Queries("key", "{key}").HandlerFunc(YourHandler).Name("YourHandler")
//     router.Path("/articles/{id:[0-9]+}").HandlerFunc(YourHandler)
//
// id := mux.Vars(r)["id"]
//     key := r.FormValue("key")

//     u, err := router.Get("YourHandler").URL("id", id, "key", key)
//     if err != nil {
//         http.Error(w, err.Error(), 500)
//         return
//     }
