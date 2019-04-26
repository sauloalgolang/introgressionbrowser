package endpoints

import (
	log "github.com/sirupsen/logrus"
)

var databases DbDb

func init() {
	log.Trace("web :: endpoints :: init()")
	databases = *NewDbDb()
}
