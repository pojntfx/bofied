package providers

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type SetupProviderChildrenProps struct {
	BackendURL      string
	OIDCIssuer      string
	OIDCClientID    string
	OIDCRedirectURL string
	Ready           bool

	SetBackendURL,
	SetOIDCIssuer,
	SetOIDCClientID,
	SetOIDCRedirectURL func(string, app.Context)
	ApplyConfig func(app.Context)

	Error error
}

type SetupProvider struct {
	app.Compo

	StoragePrefix       string
	StateQueryParameter string
	CodeQueryParameter  string
	Children            func(SetupProviderChildrenProps) app.UI

	backendURL      string
	oidcIssuer      string
	oidcClientID    string
	oidcRedirectURL string
	ready           bool

	err error
}

const (
	backendURLKey      = "backendURL"
	oidcIssuerKey      = "oidcIssuer"
	oidcClientIDKey    = "oidcClientID"
	oidcRedirectURLKey = "oidcRedirectURL"
)

func (c *SetupProvider) Render() app.UI {
	return c.Children(SetupProviderChildrenProps{
		BackendURL:      c.backendURL,
		OIDCIssuer:      c.oidcIssuer,
		OIDCClientID:    c.oidcClientID,
		OIDCRedirectURL: c.oidcRedirectURL,
		Ready:           c.ready,

		SetBackendURL: func(s string, ctx app.Context) {
			c.ready = false
			c.backendURL = s
		},
		SetOIDCIssuer: func(s string, ctx app.Context) {
			c.ready = false
			c.oidcIssuer = s
		},
		SetOIDCClientID: func(s string, ctx app.Context) {
			c.ready = false
			c.oidcClientID = s
		},
		SetOIDCRedirectURL: func(s string, ctx app.Context) {
			c.ready = false
			c.oidcRedirectURL = s
		},
		ApplyConfig: func(ctx app.Context) {
			c.validate(ctx)
		},

		Error: c.err,
	})
}

func (c *SetupProvider) invalidate(err error) {
	// Set the error state
	c.err = err
	c.ready = false
}

func (c *SetupProvider) validate(ctx app.Context) {
	// Validate fields
	if c.oidcClientID == "" {
		c.invalidate(errors.New("invalid OIDC client ID"))

		return
	}

	if _, err := url.ParseRequestURI(c.oidcIssuer); err != nil {
		c.invalidate(fmt.Errorf("invalid OIDC issuer: %v", err))

		return
	}

	if _, err := url.ParseRequestURI(c.backendURL); err != nil {
		c.invalidate(fmt.Errorf("invalid backend URL: %v", err))

		return
	}

	if _, err := url.ParseRequestURI(c.oidcRedirectURL); err != nil {
		c.invalidate(fmt.Errorf("invalid OIDC redirect URL: %v", err))

		return
	}

	// Persist state
	if err := c.persist(ctx); err != nil {
		c.invalidate(err)

		return
	}

	// If all are valid, set ready state
	c.err = nil
	c.ready = true
}

func (c *SetupProvider) persist(ctx app.Context) error {
	// Write state to storage
	if err := ctx.LocalStorage().Set(c.getKey(backendURLKey), c.backendURL); err != nil {
		return err
	}
	if err := ctx.LocalStorage().Set(c.getKey(oidcIssuerKey), c.oidcIssuer); err != nil {
		return err
	}
	if err := ctx.LocalStorage().Set(c.getKey(oidcClientIDKey), c.oidcClientID); err != nil {
		return err
	}

	return ctx.LocalStorage().Set(c.getKey(oidcRedirectURLKey), c.oidcRedirectURL)
}

func (c *SetupProvider) rehydrateFromURL(ctx app.Context) bool {
	// Read state from URL
	query := app.Window().URL().Query()

	backendURL := query.Get(backendURLKey)
	oidcIssuer := query.Get(oidcIssuerKey)
	oidcClientID := query.Get(oidcClientIDKey)
	oidcRedirectURL := query.Get(oidcRedirectURLKey)

	// If all values are set, set them in the data provider
	if backendURL != "" && oidcIssuer != "" && oidcClientID != "" && oidcRedirectURL != "" {
		c.backendURL = backendURL
		c.oidcIssuer = oidcIssuer
		c.oidcClientID = oidcClientID
		c.oidcRedirectURL = oidcRedirectURL

		return true
	}

	return false
}

func (c *SetupProvider) rehydrateFromStorage(ctx app.Context) bool {
	// Read state from storage
	backendURL := ""
	oidcIssuer := ""
	oidcClientID := ""
	oidcRedirectURL := ""

	if err := ctx.LocalStorage().Get(c.getKey(backendURLKey), &backendURL); err != nil {
		c.invalidate(err)

		return false
	}
	if err := ctx.LocalStorage().Get(c.getKey(oidcIssuerKey), &oidcIssuer); err != nil {
		c.invalidate(err)

		return false
	}
	if err := ctx.LocalStorage().Get(c.getKey(oidcClientIDKey), &oidcClientID); err != nil {
		c.invalidate(err)

		return false
	}
	if err := ctx.LocalStorage().Get(c.getKey(oidcRedirectURLKey), &oidcRedirectURL); err != nil {
		c.invalidate(err)

		return false
	}

	// If all values are set, set them in the data provider
	if backendURL != "" && oidcIssuer != "" && oidcClientID != "" && oidcRedirectURL != "" {
		c.backendURL = backendURL
		c.oidcIssuer = oidcIssuer
		c.oidcClientID = oidcClientID
		c.oidcRedirectURL = oidcRedirectURL

		return true
	}

	return false
}

func (c *SetupProvider) rehydrateAuthenticationFromURL() bool {
	// Read state from URL
	query := app.Window().URL().Query()

	state := query.Get(c.StateQueryParameter)
	code := query.Get(c.CodeQueryParameter)

	// If all values are set, set them in the data provider
	if state != "" && code != "" {
		return true
	}

	return false
}

func (c *SetupProvider) getKey(key string) string {
	// Get a prefixed key
	return fmt.Sprintf("%v.%v", c.StoragePrefix, key)
}

func (c *SetupProvider) OnMount(ctx app.Context) {
	// Initialize state
	c.backendURL = ""
	c.oidcIssuer = ""
	c.oidcClientID = ""
	c.oidcRedirectURL = ""
	c.ready = false

	// If rehydrated from URL, validate & apply
	if c.rehydrateFromURL(ctx) {
		// Auto-apply if configured
		// Disabled until a flow for handling wrong input details has been implemented
		// c.validate()
	}

	// If rehydrated from storage, validate & apply
	if c.rehydrateFromStorage(ctx) {
		// Auto-apply if configured
		// Disabled until a flow for handling wrong input details has been implemented
		// c.validate()
	}

	// If rehydrated authentication from URL, continue
	if c.rehydrateAuthenticationFromURL() {
		// Auto-apply if configured; set ready state
		c.err = nil
		c.ready = true
	}
}
