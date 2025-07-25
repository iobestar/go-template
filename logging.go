package template

import log "github.com/sirupsen/logrus"

func SetupLogging() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: false,
		DisableColors: true,
	})
}
