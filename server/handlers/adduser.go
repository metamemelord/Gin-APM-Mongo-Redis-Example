package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/metamemelord/Gin-APM-Mongo-Redis-Example/db"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

var addUserModule = fx.Options(
	fx.Invoke(configureAddUserHandler),
)

func configureAddUserHandler(g *gin.Engine, database db.DB) {
	g.POST("/users", getAddUserHandler(database))
}

func getAddUserHandler(database db.DB) func(*gin.Context) {
	return func(g *gin.Context) {
		user := new(db.User)
		err := json.NewDecoder(g.Request.Body).Decode(user)
		if err != nil {
			log.Println(err)
			Respond(g, http.StatusBadRequest, nil, err)
			return
		}

		errs := validateUser(user)
		if len(errs) != 0 {
			Respond(g, http.StatusPreconditionFailed, nil, errs...)
			return
		}
		user.Active = true
		user, err = database.AddUser(g.Request.Context(), user)
		if err != nil {
			log.Println(err)
			Respond(g, http.StatusBadRequest, nil, err)
			return
		}
		Respond(g, http.StatusCreated, user)
	}
}

func validateUser(user *db.User) []error {
	var errs []error
	if user.FirstName == "" {
		errs = append(errs, fmt.Errorf("First name cannot be empty"))
	}
	if user.LastName == "" {
		errs = append(errs, fmt.Errorf("Last name cannot be empty"))
	}
	if user.Age <= 0 {
		errs = append(errs, fmt.Errorf("Age must be a positive number"))
	}
	return errs
}
