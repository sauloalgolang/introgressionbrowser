package endpoints

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// router.HandleFunc("/databases/{database}/chromosomes", endpoints.Chromosomes).Methods("GET")                                     //.HeadersRegexp("Content-Type", "application/json")
// router.HandleFunc("/databases/{database}/chromosomes/{chromosome}/block", endpoints.ChromosomeBlock).Methods("GET")              //.HeadersRegexp("Content-Type", "application/json")

func Chromosomes(w http.ResponseWriter, r *http.Request) {
	log.Tracef("Chromosomes %#v", r)

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

func ChromosomeBlock(w http.ResponseWriter, r *http.Request) {
	log.Tracef("ChromosomeBlock %#v", r)

	params := mux.Vars(r)
	database := params["database"]
	chromosome := params["chromosome"]

	db, ok := databases.GetChromosomeBlock(database, chromosome)

	if !ok {
		resp := Message(false, "fail")
		resp["data"] = "No such chromosome: " + chromosome + " in database " + database
		Respond(w, resp)
		return
	}

	resp := Message(true, "success")
	resp["data"] = db

	Respond(w, resp)
}
