package main

import (
	"go.uber.org/fx"
	"github.com/metamemelord/Gin-APM-Mongo-Redis-Example/configuration"
	"github.com/metamemelord/Gin-APM-Mongo-Redis-Example/server"
	"github.com/metamemelord/Gin-APM-Mongo-Redis-Example/db"
)

func main() {
	app := fx.New(
		fx.Provide(configuration.Get),
		server.Module,
		db.Module,
	)
	app.Run()
}
