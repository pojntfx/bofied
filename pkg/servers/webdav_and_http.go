package servers

import (
	"net/http"

	"github.com/pojntfx/bofied/pkg/authorization"
	"github.com/pojntfx/bofied/pkg/constants"
	"github.com/pojntfx/liwasc/pkg/validators"
	"github.com/rs/cors"
	"golang.org/x/net/webdav"
)

const (
	WebDAVRealmDescription = "bofied protected area. Please enter `" + constants.OIDCOverBasicAuthUsername + "` as the username and a OpenID Connect token (i.e. from the frontend) as the password"
	HTTPPrefix             = "/public"
	WebDAVPrefix           = "/private"
)

type WebDAVAndHTTPServer struct {
	FileServer

	oidcValidator *validators.OIDCValidator
}

func NewWebDAVAndHTTPServer(workingDir string, listenAddress string, oidcValidator *validators.OIDCValidator) *WebDAVAndHTTPServer {
	return &WebDAVAndHTTPServer{
		FileServer: FileServer{
			workingDir:    workingDir,
			listenAddress: listenAddress,
		},

		oidcValidator: oidcValidator,
	}
}

func (s *WebDAVAndHTTPServer) GetWebDAVHandler(prefix string) webdav.Handler {
	return webdav.Handler{
		Prefix:     prefix,
		FileSystem: webdav.Dir(s.workingDir),
		LockSystem: webdav.NewMemLS(),
	}
}

func (s *WebDAVAndHTTPServer) GetHTTPHandler() http.Handler {
	return http.FileServer(
		http.Dir(s.workingDir),
	)
}

func (s *WebDAVAndHTTPServer) ListenAndServe() error {
	webDAVHandler := s.GetWebDAVHandler(WebDAVPrefix)
	httpHandler := s.GetHTTPHandler()

	mux := http.NewServeMux()

	mux.Handle(
		HTTPPrefix+"/",
		http.StripPrefix("/public", httpHandler),
	)
	mux.Handle(
		WebDAVPrefix+"/",
		cors.New(cors.Options{
			AllowedMethods: []string{
				"GET",
				"PUT",
				"PROPFIND",
			},
			AllowCredentials: true,
			AllowedHeaders: []string{
				"Authorization",
				"Content-Type",
				"Depth",
			},
		}).Handler(
			authorization.OIDCOverBasicAuth(
				&webDAVHandler,
				constants.OIDCOverBasicAuthUsername,
				s.oidcValidator,
				WebDAVRealmDescription,
			),
		),
	)

	return http.ListenAndServe(
		s.listenAddress,
		mux,
	)
}
