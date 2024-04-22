package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// IsAuthenticated is a middleware that checks if
// the user has already been authenticated previously.
func IsAuthenticated(ctx *gin.Context) {
	profile := sessions.Default(ctx).Get("profile")
	if profile == nil {
		ctx.Redirect(http.StatusSeeOther, "/")
	} else {
		fmt.Printf("Profile: %+v\n", profile)
		ctx.Next()
	}
}
