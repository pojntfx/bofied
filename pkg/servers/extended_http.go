package servers

import (
	"net/http"

	"github.com/pojntfx/bofied/pkg/authorization"
	"github.com/pojntfx/bofied/pkg/constants"
	"github.com/pojntfx/bofied/pkg/eventing"
	"github.com/pojntfx/bofied/pkg/validators"
	"github.com/rs/cors"
	"golang.org/x/net/webdav"
)

const (
	WebDAVRealmDescription = `bofied protected area. You can find your credentials (username and password/token) with the "Mount Folder" option in the frontend.`
	HTTPPrefix             = "/public"
	WebDAVPrefix           = "/private"
	GRPCPrefix             = "/grpc"
)

type ExtendedHTTPServer struct {
	FileServer

	eventsServerHandler http.Handler
	eventHandler        *eventing.EventHandler

	oidcValidator *validators.OIDCValidator
}

func NewExtendedHTTPServer(
	workingDir string,
	listenAddress string,
	oidcValidator *validators.OIDCValidator,
	eventsServerHandler http.Handler,
	eventHandler *eventing.EventHandler,
) *ExtendedHTTPServer {
	return &ExtendedHTTPServer{
		FileServer: FileServer{
			workingDir:    workingDir,
			listenAddress: listenAddress,
		},

		eventsServerHandler: eventsServerHandler,
		eventHandler:        eventHandler,

		oidcValidator: oidcValidator,
	}
}

func (s *ExtendedHTTPServer) GetWebDAVHandler(prefix string) webdav.Handler {
	return webdav.Handler{
		Prefix:     prefix,
		FileSystem: webdav.Dir(s.workingDir),
		LockSystem: webdav.NewMemLS(),
	}
}

func (s *ExtendedHTTPServer) GetHTTPHandler() http.Handler {
	return eventing.LogRequestHandler(
		http.FileServer(
			http.Dir(s.workingDir),
		),
		s.eventHandler,
	)
}

func (s *ExtendedHTTPServer) ListenAndServe() error {
	webDAVHandler := s.GetWebDAVHandler(WebDAVPrefix)
	httpHandler := s.GetHTTPHandler()

	mux := http.NewServeMux()

	mux.Handle(
		HTTPPrefix+"/",
		http.StripPrefix(HTTPPrefix, httpHandler),
	)
	mux.Handle(
		WebDAVPrefix+"/",
		cors.New(cors.Options{
			AllowedMethods: []string{
				"GET",
				"PUT",
				"PROPFIND",
				"MKCOL",
				"MOVE",
				"COPY",
				"DELETE",
			},
			AllowCredentials: true,
			AllowedHeaders: []string{
				"*",
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
	mux.Handle(
		GRPCPrefix,
		s.eventsServerHandler,
	)

	return http.ListenAndServe(
		s.listenAddress,
		mux,
	)
}
