package servers

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pin/tftp"
	"github.com/pojntfx/bofied/pkg/eventing"
)

type TFTPServer struct {
	FileServer

	EventHandler *eventing.EventHandler
}

func NewTFTPServer(workingDir string, listenAddress string, eventHandler *eventing.EventHandler) *TFTPServer {
	return &TFTPServer{
		FileServer: FileServer{
			workingDir:    workingDir,
			listenAddress: listenAddress,
		},

		EventHandler: eventHandler,
	}
}

func (s *TFTPServer) ListenAndServe() error {
	h := tftp.NewServer(
		func(filename string, rf io.ReaderFrom) error {
			// Get remote IP
			raddr := rf.(tftp.OutgoingTransfer).RemoteAddr()

			// Prevent accessing any parent directories
			fullFilename := filepath.Join(s.workingDir, filename)
			if strings.Contains(filename, "..") {
				s.EventHandler.Emit(`could not send file: get request to file "%v" by client "%v" blocked because it is located outside the working directory "%v"`, fullFilename, raddr.String(), s.workingDir)

				return errors.New("unauthorized: tried to access file outside working directory")
			}

			// Open file to send
			file, err := os.Open(fullFilename)
			if err != nil {
				s.EventHandler.Emit(`could not open file "%v" for client "%v": %v`, fullFilename, raddr.String(), err)

				return err
			}

			// Send the file to the client
			n, err := rf.ReadFrom(file)
			s.EventHandler.Emit(`sent file "%v" (%v bytes) to client "%v"`, fullFilename, n, raddr.String())

			return err
		},
		nil,
	)

	return h.ListenAndServe(s.listenAddress)
}
