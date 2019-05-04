package endpoints

import (
	// "errors"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	// "strconv"
)

// router.HandleFunc(PLOTS_ENDPOINT+"/{database}/{chromosome}/{referenceName}", endpoints.Plots).Methods("GET").Name("plots")

func plotsParams(w http.ResponseWriter, r *http.Request) (database string, chromosome string, referenceName string, msg string, ok bool) {
	params := mux.Vars(r)
	database = params["database"]
	chromosome = params["chromosome"]
	referenceName = params["referenceName"]
	ok = true
	msg = ""

	if !ok {
		msg = fmt.Sprintf("error getting plot parameters")
		resp := Message(false, "fail")
		resp["data"] = msg
		Respond(w, resp)
		return
	}

	return
}

// Plots handle plots requests
func Plots(w http.ResponseWriter, r *http.Request) {
	log.Tracef("Plots %#v", r)

	database, chromosome, referenceName, msg, ok := plotsParams(w, r)

	if !ok {
		return
	}

	table, pOk := databases.GetPlotTable(database, chromosome, referenceName)

	if !pOk {
		msg = fmt.Sprintf("error getting plot for database %s chromosome %s reference %s", database, chromosome, referenceName)
		resp := Message(false, "fail")
		resp["data"] = msg
		Respond(w, resp)
		return
	}

	resp := Message(true, "success")
	resp["data"] = table
	Respond(w, resp)
}
