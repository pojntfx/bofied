package servers

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	tftp "github.com/pin/tftp/v3"
	"github.com/pojntfx/bofied/pkg/eventing"
)

type TFTPServer struct {
	FileServer

	eventHandler *eventing.EventHandler
}

func NewTFTPServer(workingDir string, listenAddress string, eventHandler *eventing.EventHandler) *TFTPServer {
	return &TFTPServer{
		FileServer: FileServer{
			workingDir:    workingDir,
			listenAddress: listenAddress,
		},

		eventHandler: eventHandler,
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
				s.eventHandler.Emit(`could not send file: get request to file "%v" by client "%v" blocked because it is located outside the working directory "%v"`, fullFilename, raddr.String(), s.workingDir)

				return errors.New("unauthorized: tried to access file outside working directory")
			}

			// Open file to send
			file, err := os.Open(fullFilename)
			if err != nil {
				s.eventHandler.Emit(`could not open file "%v" for client "%v": %v`, fullFilename, raddr.String(), err)

				return err
			}

			// Send the file to the client
			n, err := rf.ReadFrom(file)
			if err != nil {
				s.eventHandler.Emit(`could not sent file "%v" to client "%v": %v`, fullFilename, raddr.String(), err)

				return err
			}

			s.eventHandler.Emit(`sent file "%v" (%v bytes) to client "%v"`, fullFilename, n, raddr.String())

			return nil
		},
		nil,
	)

	return h.ListenAndServe(s.listenAddress)
}
