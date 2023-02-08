package logs

import (
	"os"

	log "gopkg.in/inconshreveable/log15.v2"
)

var Logger log.Logger

func CreateLogger() {
	l := log.New()
	l.SetHandler(log.StreamHandler(os.Stderr, log.JsonFormat()))

	Logger = l
}
