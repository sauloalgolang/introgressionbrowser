package endpoints

import (
	// "github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// var ROUTER *mux.Router
var databases DbDb

func init() {
	log.Trace("web :: endpoints :: init()")
	databases = *NewDbDb()
}

// func SetRouter(router *mux.Router) {
// 	ROUTER = router
// }
