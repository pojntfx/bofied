package services

//go:generate sh -c "mkdir -p ../api/proto/v1 && protoc --go_out=paths=source_relative:../api/proto/v1 --go-grpc_out=paths=source_relative:../api/proto/v1 -I=../../api/proto/v1 ../../api/proto/v1/*.proto"

import (
	"context"

	api "github.com/pojntfx/bofied/pkg/api/proto/v1"
	"github.com/pojntfx/bofied/pkg/validators"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MetadataService struct {
	api.UnimplementedMetadataServiceServer

	advertisedIP string
	tftpPort     int32
	httpPort     int32

	contextValidator *validators.ContextValidator
}

func NewMetadataService(
	advertisedIP string,
	tftpPort int32,
	httpPort int32,
	contextValidator *validators.ContextValidator,
) *MetadataService {
	return &MetadataService{
		advertisedIP: advertisedIP,
		tftpPort:     tftpPort,
		httpPort:     httpPort,

		contextValidator: contextValidator,
	}
}

func (s *MetadataService) GetMetadata(ctx context.Context, _ *api.Empty) (*api.MetadataMessage, error) {
	// Authorize
	valid, err := s.contextValidator.Validate(ctx)
	if err != nil || !valid {
		return nil, status.Errorf(codes.Unauthenticated, "could not authorize: %v", err)
	}

	// Return the constructed message
	return &api.MetadataMessage{
		AdvertisedIP: s.advertisedIP,
		TFTPPort:     s.tftpPort,
		HTTPPort:     s.httpPort,
	}, nil
}
