package services

//go:generate sh -c "mkdir -p ../api/proto/v1 && protoc --go_out=paths=source_relative:../api/proto/v1 --go-grpc_out=paths=source_relative:../api/proto/v1 -I=../../api/proto/v1 ../../api/proto/v1/*.proto"

import (
	"fmt"
	"log"
	"time"

	api "github.com/pojntfx/bofied/pkg/api/proto/v1"
	"github.com/pojntfx/bofied/pkg/eventing"
	"github.com/pojntfx/bofied/pkg/validators"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	AuthorizationMetadataKey = "X-Bofied-Authorization"
)

type EventsService struct {
	api.UnimplementedEventsServiceServer

	eventsHandler *eventing.EventHandler

	contextValidator *validators.ContextValidator
}

func NewEventsService(eventsHandler *eventing.EventHandler, contextValidator *validators.ContextValidator) *EventsService {
	return &EventsService{
		eventsHandler: eventsHandler,

		contextValidator: contextValidator,
	}
}

func (s *EventsService) SubscribeToEvents(_ *api.Empty, stream api.EventsService_SubscribeToEventsServer) error {
	// Authorize
	valid, err := s.contextValidator.Validate(stream.Context())
	if err != nil || !valid {
		return status.Errorf(codes.Unauthenticated, "could not authorize: %v", err)
	}

	// Subscribe to events
	events, err := s.eventsHandler.Sub()
	if err != nil {
		msg := fmt.Sprintf("could not get events from messenger: %v", err)

		log.Println(msg)

		return status.Error(codes.Unknown, msg)
	}
	defer s.eventsHandler.Unsub(events)

	// Receive event from bus
	for event := range events {
		e := event.(eventing.Event)

		// Send event to client
		if err := stream.Send(&api.EventMessage{
			CreatedAt: e.CreatedAt.Format(time.RFC3339),
			Message:   e.Message,
		}); err != nil {
			log.Printf("could send event to client: %v\n", err)

			return err
		}
	}

	return nil
}
