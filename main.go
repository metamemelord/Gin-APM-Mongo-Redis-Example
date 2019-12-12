package main

import (
	"github.com/metamemelord/Gin-APM-Mongo-Redis-Example/configuration"
	"github.com/metamemelord/Gin-APM-Mongo-Redis-Example/db"
	"github.com/metamemelord/Gin-APM-Mongo-Redis-Example/server"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.Provide(configuration.Get),
		server.Module,
		db.Module,
	)
	app.Run()
}
