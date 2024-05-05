package login

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"

	"01-Login/platform/authenticator"
)

// Handler for our login.
func Handler(auth *authenticator.Authenticator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		state, err := generateRandomState()
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		// Save the state inside the session, which is a cookie-session
		// https://github.com/gorilla/sessions/blob/main/store.go#L105
		// TODO: learn how cookies work in the browser
		session := sessions.Default(ctx)
		session.Set("state", state)
		if err := session.Save(); err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.Redirect(http.StatusTemporaryRedirect, auth.AuthCodeURL(
			state,
			// https://community.auth0.com/t/why-is-my-access-token-not-a-jwt-opaque-token/31028
			// originally I was missing this, and I was getting an opaque token, so I couldn't verify the access token
			// I was getting an opaque token because I was missing the audience
			// by adding an api named `https://api.tb-sb.com` at https://manage.auth0.com/dashboard/us/tg-sb/apis
			// i was able to set the audience param in the authorization request to the
			// authorization endpoint (docs found here https://auth0.com/docs/api/authentication#authorize47)
			// you can obtain this information by looking at the `authorization_endpoint` of my tennant's well known config
			// https://tg-sb.us.auth0.com/.well-known/openid-configuration
			oauth2.SetAuthURLParam("audience", "https://api.tg-sb.com"),
		))
	}
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}
