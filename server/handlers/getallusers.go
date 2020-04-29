package handlers

import (
	"encoding/json"
	"github.com/metamemelord/Gin-APM-Mongo-Redis-Example/db"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

var allUsersModule = fx.Options(
	fx.Invoke(configureAllUsersHandlers),
)

func configureAllUsersHandlers(g *gin.Engine, database db.DB) {
	g.GET("/users", getAllUsers(database))
}

func getAllUsers(database db.DB) func(*gin.Context) {
	return func(g *gin.Context) {
		filters := map[string]interface{}{}
		query := g.Query("q")
		_ = json.Unmarshal([]byte(query), &filters)

		if len(filters) == 0 {
			users, err := database.Find(g.Request.Context())
			if err != nil {
				log.Println("ERROR:", err)
				Respond(g, http.StatusCreated, nil, err)
				return
			}
			Respond(g, http.StatusCreated, users)
		} else {
			users, err := database.FindByFilters(g.Request.Context(), filters)
			if err != nil {
				log.Println("ERROR:", err)
				Respond(g, http.StatusCreated, nil, err)
				return
			}
			Respond(g, http.StatusCreated, users)
		}

	}
}
