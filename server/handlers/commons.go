package handlers

import "github.com/gin-gonic/gin"

func Respond(g *gin.Context, status int, payload interface{}, errs []error) {
	if errs != nil || len(errs) != 0 {
		errorStrings := []string{}
		for _, err := range errs {
			errorStrings = append(errorStrings, err.Error())
		}
		g.AbortWithStatusJSON(status, map[string][]string{"errors": errorStrings})
	} else {
		g.JSON(status, payload)
	}
}
