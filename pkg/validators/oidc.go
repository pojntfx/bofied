package validators

import (
	"context"

	"github.com/coreos/go-oidc/v3/oidc"
)

type OIDCValidator struct {
	Issuer   string
	ClientID string

	verifier *oidc.IDTokenVerifier
}

func NewOIDCValidator(issuer string, clientID string) *OIDCValidator {
	return &OIDCValidator{ClientID: clientID, Issuer: issuer}
}

func (v *OIDCValidator) Open() error {
	provider, err := oidc.NewProvider(context.Background(), v.Issuer)
	if err != nil {
		return err
	}

	v.verifier = provider.Verifier(&oidc.Config{ClientID: v.ClientID})

	return nil
}

func (v *OIDCValidator) Validate(token string) (*oidc.IDToken, error) {
	idToken, err := v.verifier.Verify(context.Background(), token)
	if err != nil {
		return nil, err
	}

	return idToken, nil
}
