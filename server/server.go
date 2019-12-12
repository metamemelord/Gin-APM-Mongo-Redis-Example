package server

import (
	"context"
	"log"
	"strconv"

	"github.com/metamemelord/Gin-APM-Mongo-Redis-Example/configuration"
	"github.com/metamemelord/Gin-APM-Mongo-Redis-Example/server/handlers"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewServer),
	handlers.Module,
	fx.Invoke(InitialiseServer),
)

func NewServer() (*gin.Engine, error) {
	return gin.New(), nil
}

func InitialiseServer(g *gin.Engine, lc fx.Lifecycle, config *configuration.Configuration) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("Starting the application...")
			go func() {
				err := g.Run(":" + strconv.Itoa(config.Port))
				if err != nil {
					log.Fatal(err)
				}
			}()
			return nil
		},
	})
}
