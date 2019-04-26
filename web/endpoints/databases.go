package endpoints

import (
	// "encoding/json"
	// "go-contacts/models"
	// u "go-contacts/utils"
	// "strconv"
	// "fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// router.HandleFunc("/databases", endpoints.Databases).Methods("GET")                                                              //.HeadersRegexp("Content-Type", "application/json")
// router.HandleFunc("/databases/{database}/block", endpoints.DatabaseBlock).Methods("GET")                                         //.HeadersRegexp("Content-Type", "application/json")

func Databases(w http.ResponseWriter, r *http.Request) {
	log.Tracef("Databases %#v", r)

	dbs := databases.GetDatabases()

	resp := Message(true, "success")
	resp["data"] = dbs
	Respond(w, resp)
}

func DatabaseBlock(w http.ResponseWriter, r *http.Request) {
	log.Tracef("DatabaseBlock %#v", r)

	params := mux.Vars(r)
	database := params["database"]

	db, ok := databases.GetDatabaseBlock(database)

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
