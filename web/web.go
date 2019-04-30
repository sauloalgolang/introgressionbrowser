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

const API_ENDPOINT = "/api"
const DATA_ENDPOINT = "/data"
const DATABASE_ENDPOINT = "/databases"
const PLOTS_ENDPOINT = "/plots"

func NewWeb(databaseDir string, httpDir string, host string, port int, verbosityLevel log.Level) {
	router := mux.NewRouter()

	router.StrictSlash(false)
	// endpoints.SetRouter(router)

	log.Warn("open your browser at http://" + host + ":" + strconv.Itoa(port))

	api := router.PathPrefix(API_ENDPOINT).Subrouter()
	api.StrictSlash(false)
	router.HandleFunc(API_ENDPOINT+"/", Template).Methods("GET").Name("apiSlash")

	newData(databaseDir, router)
	newApi(databaseDir, api, verbosityLevel)
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
	router.PathPrefix(DATA_ENDPOINT + "/").Handler(http.StripPrefix(DATA_ENDPOINT+"/", http.FileServer(http.Dir(dir)))).Name("dataSlash")
	router.PathPrefix(DATA_ENDPOINT).Handler(http.StripPrefix(DATA_ENDPOINT, http.FileServer(http.Dir(dir)))).Name("data")
}

func newApi(dir string, router *mux.Router, verbosityLevel log.Level) {
	//.HeadersRegexp("Content-Type", "application/json")
	router.HandleFunc("", Template).Methods("GET").Name("api")
	router.HandleFunc("/update", endpoints.Update).Methods("POST").Name("update")
	router.HandleFunc(DATABASE_ENDPOINT, endpoints.Databases).Methods("GET").Name("databases")
	router.HandleFunc(DATABASE_ENDPOINT+"/{database}", endpoints.Database).Methods("GET").Name("database")
	router.HandleFunc(DATABASE_ENDPOINT+"/{database}/summary", endpoints.DatabaseSummary).Methods("GET").Name("databaseSummary")
	router.HandleFunc(DATABASE_ENDPOINT+"/{database}/summary/matrix", endpoints.DatabaseSummaryMatrix).Methods("GET").Name("databaseSummaryMatrix")
	router.HandleFunc(DATABASE_ENDPOINT+"/{database}/summary/matrix/table", endpoints.DatabaseSummaryMatrixTable).Methods("GET").Name("databaseSummaryMatrixTable")
	router.HandleFunc(DATABASE_ENDPOINT+"/{database}/chromosome", endpoints.Chromosomes).Methods("GET").Name("databaseChromosomes")
	router.HandleFunc(DATABASE_ENDPOINT+"/{database}/chromosomes/{chromosome}", endpoints.Chromosome).Methods("GET").Name("databaseChromosome")
	router.HandleFunc(DATABASE_ENDPOINT+"/{database}/chromosomes/{chromosome}/summary", endpoints.ChromosomeSummary).Methods("GET").Name("databaseChromosomeSummary")
	router.HandleFunc(DATABASE_ENDPOINT+"/{database}/chromosomes/{chromosome}/summary/matrix", endpoints.ChromosomeSummaryMatrix).Methods("GET").Name("databaseChromosomeSummaryMatrix")
	router.HandleFunc(DATABASE_ENDPOINT+"/{database}/chromosomes/{chromosome}/summary/matrix/table", endpoints.ChromosomeSummaryMatrixTable).Methods("GET").Name("databaseChromosomeSummaryTable")
	router.HandleFunc(DATABASE_ENDPOINT+"/{database}/chromosomes/{chromosome}/block", endpoints.Blocks).Methods("GET").Name("databaseChromosomeBlocks")
	router.HandleFunc(DATABASE_ENDPOINT+"/{database}/chromosomes/{chromosome}/blocks/{blockNum:[0-9]+}", endpoints.Block).Methods("GET").Name("databaseChromosomeBlock")
	router.HandleFunc(DATABASE_ENDPOINT+"/{database}/chromosomes/{chromosome}/blocks/{blockNum:[0-9]+}/matrix", endpoints.BlockMatrix).Methods("GET").Name("databaseChromosomeBlockMatrix")
	router.HandleFunc(DATABASE_ENDPOINT+"/{database}/chromosomes/{chromosome}/blocks/{blockNum:[0-9]+}/matrix/table", endpoints.BlocksMatrixTable).Methods("GET").Name("databaseChromosomeBlockMatrixTable")

	router.HandleFunc(PLOTS_ENDPOINT+"/{database}/{chromosome}/{referenceName}", endpoints.Plots).Methods("GET").Name("plots")

	route := router.Get("data")

	if route == nil {
		log.Panic("No data router")
	}

	tmpl, t_err := route.GetPathTemplate()
	if t_err != nil {
		log.Panic("No data router template")
	}

	tmpl = strings.TrimSuffix(tmpl, "/")

	endpoints.DATA_ENDPOINT = tmpl
	endpoints.DATABASE_DIR = strings.TrimSuffix(dir, "/")
	endpoints.VERBOSITY = verbosityLevel

	endpoints.ListDatabases()
}

func Template(w http.ResponseWriter, r *http.Request) {
	log.Tracef("Template %#v", r)

	resp := endpoints.Message(true, "success")

	tmp := map[string]interface{}{
		API_ENDPOINT + DATABASE_ENDPOINT + "":                                                                           []string{""},
		API_ENDPOINT + DATABASE_ENDPOINT + "/{database}":                                                                endpoints.DatabaseInfo{},
		API_ENDPOINT + DATABASE_ENDPOINT + "/{database}/summary":                                                        endpoints.BlockInfo{},
		API_ENDPOINT + DATABASE_ENDPOINT + "/{database}/summary/matrix":                                                 endpoints.MatrixInfo{},
		API_ENDPOINT + DATABASE_ENDPOINT + "/{database}/summary/matrix/table":                                           endpoints.TableInfo{},
		API_ENDPOINT + DATABASE_ENDPOINT + "/{database}/chromosome":                                                     []endpoints.ChromosomeInfo{endpoints.ChromosomeInfo{}},
		API_ENDPOINT + DATABASE_ENDPOINT + "/{database}/chromosomes/{chromosome}":                                       endpoints.ChromosomeInfo{},
		API_ENDPOINT + DATABASE_ENDPOINT + "/{database}/chromosomes/{chromosome}/summary":                               endpoints.BlockInfo{},
		API_ENDPOINT + DATABASE_ENDPOINT + "/{database}/chromosomes/{chromosome}/summary/matrix":                        endpoints.MatrixInfo{},
		API_ENDPOINT + DATABASE_ENDPOINT + "/{database}/chromosomes/{chromosome}/summary/matrix/table":                  endpoints.TableInfo{},
		API_ENDPOINT + DATABASE_ENDPOINT + "/{database}/chromosomes/{chromosome}/block":                                 []endpoints.BlockInfo{endpoints.BlockInfo{}},
		API_ENDPOINT + DATABASE_ENDPOINT + "/{database}/chromosomes/{chromosome}/blocks/{blockNum:[0-9]+}":              endpoints.BlockInfo{},
		API_ENDPOINT + DATABASE_ENDPOINT + "/{database}/chromosomes/{chromosome}/blocks/{blockNum:[0-9]+}/matrix":       endpoints.MatrixInfo{},
		API_ENDPOINT + DATABASE_ENDPOINT + "/{database}/chromosomes/{chromosome}/blocks/{blockNum:[0-9]+}/matrix/table": endpoints.TableInfo{},
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
