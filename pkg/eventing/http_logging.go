package eventing

import (
	"net/http"
)

func LogRequestHandler(h http.Handler, eventHandler *EventHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			eventHandler.Emit(`sending file "%v" to client "%v" with user agent "%v"`, r.URL.Path, r.RemoteAddr, r.UserAgent())
		}

		h.ServeHTTP(w, r)
	})
}
