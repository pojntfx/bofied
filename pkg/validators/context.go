package validators

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/grpc/metadata"
)

type ContextValidator struct {
	metadataKey   string
	oidcValidator *OIDCValidator
}

func NewContextValidator(metadataKey string, oidcValidator *OIDCValidator) *ContextValidator {
	return &ContextValidator{
		metadataKey:   metadataKey,
		oidcValidator: oidcValidator,
	}
}

func (v *ContextValidator) Validate(ctx context.Context) (bool, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return false, errors.New("could not parse metadata")
	}

	token := md.Get(v.metadataKey)
	if len(token) <= 0 {
		return false, errors.New("could not parse metadata")
	}

	idToken, err := v.oidcValidator.Validate(token[0])
	if err != nil || idToken == nil {
		return false, fmt.Errorf("invalid token: %v", err)
	}

	return true, nil
}
