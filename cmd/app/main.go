package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/geeeeorge/AWS-app-runner/config"
	"github.com/geeeeorge/AWS-app-runner/infra/database"
	"github.com/geeeeorge/AWS-app-runner/infra/server"
)

func main() {
	conf := config.Load()
	db, err := database.NewMysqlDB(conf.RDBConfig.ConnectionString())
	if err != nil {
		log.Fatalf(err.Error())
	}
	s := server.NewServer(conf.AppPort(), conf.AppHost(), db, nil)
	s.Start()
}
