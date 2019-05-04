package endpoints

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

// router.HandleFunc(DATABASE_ENDPOINT+"/{database}/chromosome/{chromosome}/block", endpoints.Blocks).Methods("GET")                                           //.HeadersRegexp("Content-Type", "application/json")
// router.HandleFunc(DATABASE_ENDPOINT+"/{database}/chromosome/{chromosome}/block/{blockNum:[0-9]+}", endpoints.Block).Methods("GET")                          //.HeadersRegexp("Content-Type", "application/json")
// router.HandleFunc(DATABASE_ENDPOINT+"/{database}/chromosome/{chromosome}/block/{blockNum:[0-9]+}/matrix", endpoints.BlockMatrix).Methods("GET")             //.HeadersRegexp("Content-Type", "application/json")
// router.HandleFunc(DATABASE_ENDPOINT+"/{database}/chromosome/{chromosome}/block/{blockNum:[0-9]+}/matrix/table", endpoints.BlocksMatrixTable).Methods("GET") //.HeadersRegexp("Content-Type", "application/json")

// Blocks handle calls for blocks
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

func blockParams(r *http.Request) (database string, chromosome string, blockNum uint64, msg string, ok bool) {
	params := mux.Vars(r)
	database = params["database"]
	chromosome = params["chromosome"]
	blockNumS := params["blockNum"]
	ok = true
	err := errors.New("")

	blockNum, err = strconv.ParseUint(blockNumS, 10, 64)

	if err != nil {
		msg = "Invalid blockNum: " + blockNumS + ". Not a number"
		ok = false
		return
	}
	return
}

func getBlock(w http.ResponseWriter, r *http.Request) (database string, chromosome string, blockNum uint64, msg string, ok bool) {
	database, chromosome, blockNum, msg, ok = blockParams(r)

	if !ok {
		resp := Message(false, "fail")
		resp["data"] = msg
		Respond(w, resp)
		return
	}

	return
}

// Block handles calls for block
func Block(w http.ResponseWriter, r *http.Request) {
	log.Tracef("Block %#v", r)

	database, chromosome, blockNum, msg, ok := getBlock(w, r)

	if !ok {
		return
	}

	block, bOk := databases.GetBlock(database, chromosome, blockNum)

	if !bOk {
		msg = fmt.Sprintf("No such blockNum: %d in chromosome: %s in database %s", blockNum, chromosome, database)

		resp := Message(false, "fail")
		resp["data"] = msg
		Respond(w, resp)
		return
	}

	resp := Message(true, "success")
	resp["data"] = block

	Respond(w, resp)
}

// BlockMatrix handles calls for block matrix
func BlockMatrix(w http.ResponseWriter, r *http.Request) {
	log.Tracef("BlockMatrix %#v", r)

	database, chromosome, blockNum, msg, ok := getBlock(w, r)

	if !ok {
		return
	}

	matrix, bOk := databases.GetBlockMatrix(database, chromosome, blockNum)

	if !bOk {
		msg = fmt.Sprintf("No such blockNum: %d in chromosome: %s in database %s", blockNum, chromosome, database)

		resp := Message(false, "fail")
		resp["data"] = msg
		Respond(w, resp)
		return
	}

	resp := Message(true, "success")
	resp["data"] = matrix

	Respond(w, resp)
}

// BlocksMatrixTable handles calls for block matrix table
func BlocksMatrixTable(w http.ResponseWriter, r *http.Request) {
	log.Tracef("BlocksMatrixTable %#v", r)

	database, chromosome, blockNum, msg, ok := getBlock(w, r)

	if !ok {
		return
	}

	table, bOk := databases.GetBlockMatrixTable(database, chromosome, blockNum)

	if !bOk {
		msg = fmt.Sprintf("No such blockNum: %d in chromosome: %s in database %s", blockNum, chromosome, database)

		resp := Message(false, "fail")
		resp["data"] = msg
		Respond(w, resp)
		return
	}

	resp := Message(true, "success")
	resp["data"] = table

	Respond(w, resp)
}
