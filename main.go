package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/shemanaev/updtr/internal/meta"
	"github.com/shemanaev/updtr/internal/server"
)

func init() {
	log.SetLevel(log.InfoLevel)
}

func main() {
	log.Infof("Starting updtr server %s (%v)", meta.Version, meta.BuildDate)
	if err := server.RunServer(); err != nil {
		log.Fatal(err)
	}
}
