package servers

import (
	"net/http"

	"golang.org/x/net/webdav"
)

type WebDAVServer struct {
	FileServer
}

func NewWebDAVServer(workingDir string, listenAddress string) *WebDAVServer {
	return &WebDAVServer{
		FileServer: FileServer{
			workingDir:    workingDir,
			listenAddress: listenAddress,
		},
	}
}

func (s *WebDAVServer) ListenAndServe() error {
	h := &webdav.Handler{
		FileSystem: webdav.Dir(s.workingDir),
		LockSystem: webdav.NewMemLS(),
	}

	return http.ListenAndServe(s.listenAddress, h)
}
