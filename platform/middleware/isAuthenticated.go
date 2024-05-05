package middleware

import (
	"01-Login/platform/authenticator"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// IsAuthenticated is a middleware that checks if
// the user has already been authenticated previously.
func IsAuthenticated(auth *authenticator.Authenticator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken := sessions.Default(ctx).Get("access_token")
		if accessToken == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		accessTokenString, ok := accessToken.(string)
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		fmt.Printf("Access Token: %s\n", accessTokenString)
		tok, err := auth.VerifyToken(ctx.Request.Context(), accessTokenString)
		if err != nil {
			fmt.Printf("Error verifying access token: %s\n", err.Error())
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if !tok.Valid {
			fmt.Printf("Token is invalid\n")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		fmt.Printf("Token: %+v\n", tok)
		claims, ok := tok.Claims.(jwt.MapClaims)
		if !ok {
			fmt.Printf("Error getting claims\n")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		fmt.Printf("Claims: %+v\n", claims)
		fmt.Printf("Scope: %+v\n", claims["scope"])
		profile := sessions.Default(ctx).Get("profile")
		if profile == nil {
			ctx.Redirect(http.StatusSeeOther, "/")
		} else {
			fmt.Printf("Profile: %+v\n", profile)
			ctx.Next()
		}
	}
}
