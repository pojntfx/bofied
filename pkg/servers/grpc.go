package servers

import (
	"net"
	"sync"

	api "github.com/pojntfx/bofied/pkg/api/proto/v1"
	"github.com/pojntfx/bofied/pkg/services"
	"github.com/pojntfx/bofied/pkg/websocketproxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	listenAddress string

	eventsService   *services.EventsService
	metadataService *services.MetadataService

	proxy *websocketproxy.WebSocketProxyServer
}

func NewGRPCServer(listenAddress string, eventsService *services.EventsService, metadataService *services.MetadataService) (*GRPCServer, *websocketproxy.WebSocketProxyServer) {
	proxy := websocketproxy.NewWebSocketProxyServer()

	return &GRPCServer{
		listenAddress:   listenAddress,
		eventsService:   eventsService,
		metadataService: metadataService,
		proxy:           proxy,
	}, proxy
}

func (s *GRPCServer) ListenAndServe() error {
	listener, err := net.Listen("tcp", s.listenAddress)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	reflection.Register(server)

	api.RegisterEventsServiceServer(server, s.eventsService)
	api.RegisterMetadataServiceServer(server, s.metadataService)

	doneChan := make(chan struct{})
	errChan := make(chan error)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		wg.Wait()

		close(doneChan)
	}()

	go func() {
		if err := server.Serve(listener); err != nil {
			errChan <- err
		}

		wg.Done()
	}()

	go func() {
		if err := server.Serve(s.proxy); err != nil {
			errChan <- err
		}

		wg.Done()
	}()

	select {
	case <-doneChan:
		return nil
	case <-errChan:
		return err
	}
}
