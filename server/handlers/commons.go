package handlers

import "github.com/gin-gonic/gin"

func Respond(g *gin.Context, status int, payload interface{}, errs ...error) {
	if len(errs) != 0 {
		if len(errs) == 1 {
			g.AbortWithStatusJSON(status, map[string]string{"error": errs[0].Error()})
		} else {
			errorStrings := []string{}
			for _, err := range errs {
				errorStrings = append(errorStrings, err.Error())
			}
			g.AbortWithStatusJSON(status, map[string][]string{"errors": errorStrings})
		}
	} else {
		g.JSON(status, payload)
	}
}
