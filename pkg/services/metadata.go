package services

//go:generate sh -c "mkdir -p ../api/proto/v1 && protoc --go_out=paths=source_relative,plugins=grpc:../api/proto/v1 -I=../../api/proto/v1 ../../api/proto/v1/*.proto"

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	api "github.com/pojntfx/bofied/pkg/api/proto/v1"
	"github.com/pojntfx/liwasc/pkg/validators"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MetadataService struct {
	api.UnimplementedMetadataServiceServer

	advertisedIP string

	contextValidator *validators.ContextValidator
}

func NewMetadataService(advertisedIP string, contextValidator *validators.ContextValidator) *MetadataService {
	return &MetadataService{
		advertisedIP: advertisedIP,

		contextValidator: contextValidator,
	}
}

func (s *MetadataService) GetMetadata(ctx context.Context, _ *empty.Empty) (*api.MetadataMessage, error) {
	// Authorize
	valid, err := s.contextValidator.Validate(ctx)
	if err != nil || !valid {
		return nil, status.Errorf(codes.Unauthenticated, "could not authorize: %v", err)
	}

	// Return the constructed message
	return &api.MetadataMessage{
		AdvertisedIP: s.advertisedIP,
	}, nil
}
