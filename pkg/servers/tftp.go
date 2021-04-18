package servers

import (
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pin/tftp"
)

type TFTPServer struct {
	FileServer
}

func NewTFTPServer(workingDir string, listenAddress string) *TFTPServer {
	return &TFTPServer{
		FileServer: FileServer{
			workingDir:    workingDir,
			listenAddress: listenAddress,
		},
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
				log.Printf(`could not send file: get request to file "%v" by client "%v" blocked because it is located outside the working directory "%v"`, fullFilename, raddr.String(), s.workingDir)

				return errors.New("unauthorized: tried to access file outside working directory")
			}

			// Open file to send
			file, err := os.Open(fullFilename)
			if err != nil {
				log.Printf(`could not open file "%v" for client "%v": %v`, fullFilename, raddr.String(), err)

				return err
			}

			// Send the file to the client
			n, err := rf.ReadFrom(file)
			log.Printf(`sent file "%v" (%v bytes) to client "%v"`, fullFilename, n, raddr.String())

			return err
		},
		nil,
	)

	return h.ListenAndServe(s.listenAddress)
}
