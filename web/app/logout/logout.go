package logout

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

// Handler for our logout.
func Handler(ctx *gin.Context) {
	logoutUrl, err := url.Parse("https://" + os.Getenv("AUTH0_DOMAIN") + "/v2/logout")
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	scheme := "http"
	if ctx.Request.TLS != nil {
		scheme = "https"
	}

	returnTo, err := url.Parse(scheme + "://" + ctx.Request.Host)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// https://auth0.com/docs/authenticate/login/logout/log-users-out-of-auth0#oidc-logout-endpoint-parameters
	parameters := url.Values{}
	parameters.Add("post_logout_redirect_uri", returnTo.String())
	parameters.Add("client_id", os.Getenv("AUTH0_CLIENT_ID"))
	logoutUrl.RawQuery = parameters.Encode()

	fmt.Printf("Logout URL: %s\n", logoutUrl.String())
	ctx.Redirect(http.StatusTemporaryRedirect, logoutUrl.String())
}
