package endpoints

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func Update(w http.ResponseWriter, r *http.Request) {
	log.Tracef("Update %#v", r)

	ListDatabases()

	params := mux.Vars(r)

	resp := Message(true, "success")
	resp["data"] = params

	Respond(w, resp)
}
