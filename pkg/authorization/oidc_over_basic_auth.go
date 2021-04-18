package authorization

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pojntfx/liwasc/pkg/validators"
)

func OIDCOverBasicAuth(next http.Handler, username string, oidcValidator *validators.OIDCValidator, description string) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// Validate the OIDC token which is passed as a HTTP basic auth password due to client limitations
		user, pass, ok := r.BasicAuth()
		if _, err := oidcValidator.Validate(pass); err != nil || !ok || user != username {
			// Unauthorized, log and redirect
			log.Println("could not authorized user, redirecting")

			rw.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm="%v"`, description))
			rw.WriteHeader(401)
			rw.Write([]byte("could not authorize: " + err.Error()))

			return
		}

		// Authorized, continue
		next.ServeHTTP(rw, r)
	})
}
