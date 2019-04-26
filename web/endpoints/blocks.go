package endpoints

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

// router.HandleFunc("/databases/{database}/chromosomes/{chromosome}/blocks", endpoints.Blocks).Methods("GET")                      //.HeadersRegexp("Content-Type", "application/json")
// router.HandleFunc("/databases/{database}/chromosomes/{chromosome}/blocks/{block}/block", endpoints.BlocksBlock).Methods("GET")   //.HeadersRegexp("Content-Type", "application/json")
// router.HandleFunc("/databases/{database}/chromosomes/{chromosome}/blocks/{block}/matrix", endpoints.BlocksMatrix).Methods("GET") //.HeadersRegexp("Content-Type", "application/json")

func Blocks(w http.ResponseWriter, r *http.Request) {
	log.Tracef("Blocks %#v", r)

	params := mux.Vars(r)
	database := params["database"]
	chromosome := params["chromosome"]

	blocks, ok := databases.GetBlocks(database, chromosome)

	if !ok {
		resp := Message(false, "fail")
		resp["data"] = "No such chromosome: " + chromosome + " in database " + database
		Respond(w, resp)
		return
	}

	resp := Message(true, "success")
	resp["data"] = blocks

	Respond(w, resp)
}

func BlocksBlock(w http.ResponseWriter, r *http.Request) {
	log.Tracef("BlocksBlock %#v", r)

	params := mux.Vars(r)
	database := params["database"]
	chromosome := params["chromosome"]
	blockNumS := params["blockNum"]

	blockNum, cok := strconv.ParseUint(blockNumS, 10, 64)

	if cok != nil {
		resp := Message(false, "fail")
		resp["data"] = "Invalid blockNum: " + blockNumS + ". Not a number"
		Respond(w, resp)
		return
	}

	block, ok := databases.GetBlock(database, chromosome, blockNum)

	if !ok {
		resp := Message(false, "fail")
		resp["data"] = "No such blockNum: " + blockNumS + " in chromosome: " + chromosome + " in database " + database
		Respond(w, resp)
		return
	}

	resp := Message(true, "success")
	resp["data"] = block

	Respond(w, resp)
}

func BlocksMatrix(w http.ResponseWriter, r *http.Request) {
	log.Tracef("BlocksMatrix %#v", r)

	params := mux.Vars(r)
	database := params["database"]
	chromosome := params["chromosome"]
	blockNumS := params["blockNum"]

	blockNum, cok := strconv.ParseUint(blockNumS, 10, 64)

	if cok != nil {
		resp := Message(false, "fail")
		resp["data"] = "Invalid blockNum: " + blockNumS + ". Not a number"
		Respond(w, resp)
		return
	}

	matrix, ok := databases.GetBlockMatrix(database, chromosome, blockNum)

	if !ok {
		resp := Message(false, "fail")
		resp["data"] = "No such blockNum: " + blockNumS + " in chromosome: " + chromosome + " in database " + database
		Respond(w, resp)
		return
	}

	resp := Message(true, "success")
	resp["data"] = matrix

	Respond(w, resp)
}
