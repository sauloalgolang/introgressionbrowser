package endpoints

import (
	// "encoding/json"
	// "go-contacts/models"
	// u "go-contacts/utils"
	// "strconv"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// router.HandleFunc(DATABASE_ENDPOINT, endpoints.Databases).Methods("GET")                                                                                    //.HeadersRegexp("Content-Type", "application/json")
// router.HandleFunc(DATABASE_ENDPOINT+"/{database}", endpoints.Database).Methods("GET")                                                                       //.HeadersRegexp("Content-Type", "application/json")
// router.HandleFunc(DATABASE_ENDPOINT+"/{database}/summary", endpoints.DatabaseSummary).Methods("GET")                                                        //.HeadersRegexp("Content-Type", "application/json")
// router.HandleFunc(DATABASE_ENDPOINT+"/{database}/summary/matrix", endpoints.DatabaseSummaryMatrix).Methods("GET")                                           //.HeadersRegexp("Content-Type", "application/json")
// router.HandleFunc(DATABASE_ENDPOINT+"/{database}/summary/matrix/table", endpoints.DatabaseSummaryMatrixTable).Methods("GET")                                //.HeadersRegexp("Content-Type", "application/json")

func Databases(w http.ResponseWriter, r *http.Request) {
	log.Tracef("Databases %#v", r)

	dbs := databases.GetDatabases()

	resp := Message(true, "success")
	resp["data"] = dbs
	Respond(w, resp)
}

func Database(w http.ResponseWriter, r *http.Request) {
	log.Tracef("Database %#v", r)

	params := mux.Vars(r)
	database := params["database"]

	db, ok := databases.GetDatabase(database)

	if !ok {
		resp := Message(false, "fail")
		resp["data"] = "No such database: " + database
		Respond(w, resp)
		return
	}

	resp := Message(true, "success")
	resp["data"] = db

	Respond(w, resp)
}

func DatabaseSummary(w http.ResponseWriter, r *http.Request) {
	log.Tracef("DatabaseSummary %#v", r)

	params := mux.Vars(r)
	database := params["database"]

	db, ok := databases.GetDatabaseSummaryBlock(database)

	if !ok {
		resp := Message(false, "fail")
		resp["data"] = "No such database: " + database
		Respond(w, resp)
		return
	}

	resp := Message(true, "success")
	resp["data"] = db

	Respond(w, resp)
}

func DatabaseSummaryMatrix(w http.ResponseWriter, r *http.Request) {
	log.Tracef("DatabaseBlockMatrix %#v", r)

	params := mux.Vars(r)
	database := params["database"]

	db, ok := databases.GetDatabaseSummaryBlockMatrix(database)

	if !ok {
		resp := Message(false, "fail")
		resp["data"] = "No such database: " + database
		Respond(w, resp)
		return
	}

	resp := Message(true, "success")
	resp["data"] = db

	Respond(w, resp)
}

func DatabaseSummaryMatrixTable(w http.ResponseWriter, r *http.Request) {
	log.Tracef("DatabaseBlockMatrix %#v", r)

	params := mux.Vars(r)
	database := params["database"]

	db, ok := databases.GetDatabaseSummaryBlockMatrixTable(database)

	if !ok {
		resp := Message(false, "fail")
		resp["data"] = "No such database: " + database
		Respond(w, resp)
		return
	}

	resp := Message(true, "success")
	resp["data"] = db

	Respond(w, resp)
}
