package servers

import "net/http"

type HTTPServer struct {
	FileServer
}

func NewHTTPServer(workingDir string, listenAddress string) *HTTPServer {
	return &HTTPServer{
		FileServer: FileServer{
			workingDir:    workingDir,
			listenAddress: listenAddress,
		},
	}
}

func (s *HTTPServer) ListenAndServe() error {
	h := http.FileServer(
		http.Dir(s.workingDir),
	)

	return http.ListenAndServe(s.listenAddress, h)
}
