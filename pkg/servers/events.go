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

type EventsServer struct {
	listenAddress string

	service *services.EventsService
	proxy   *websocketproxy.WebSocketProxyServer
}

func NewEventsServer(listenAddress string, service *services.EventsService) (*EventsServer, *websocketproxy.WebSocketProxyServer) {
	proxy := websocketproxy.NewWebSocketProxyServer()

	return &EventsServer{
		listenAddress: listenAddress,
		service:       service,
		proxy:         proxy,
	}, proxy
}

func (s *EventsServer) ListenAndServe() error {
	listener, err := net.Listen("tcp", s.listenAddress)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	reflection.Register(server)

	api.RegisterEventsServiceServer(server, s.service)

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
