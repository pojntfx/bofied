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
)

type WebDAVServer struct {
	FileServer

	oidcValidator *validators.OIDCValidator
}

func NewWebDAVServer(workingDir string, listenAddress string, oidcValidator *validators.OIDCValidator) *WebDAVServer {
	return &WebDAVServer{
		FileServer: FileServer{
			workingDir:    workingDir,
			listenAddress: listenAddress,
		},

		oidcValidator: oidcValidator,
	}
}

func (s *WebDAVServer) ListenAndServe() error {
	h := &webdav.Handler{
		FileSystem: webdav.Dir(s.workingDir),
		LockSystem: webdav.NewMemLS(),
	}

	return http.ListenAndServe(
		s.listenAddress,
		cors.AllowAll().Handler(
			authorization.OIDCOverBasicAuth(
				h,
				constants.OIDCOverBasicAuthUsername,
				s.oidcValidator,
				WebDAVRealmDescription,
			),
		),
	)
}
