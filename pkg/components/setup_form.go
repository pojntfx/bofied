package components

import "github.com/maxence-charriere/go-app/v9/pkg/app"

type SetupForm struct {
	app.Compo

	Error        error
	ErrorMessage string

	BackendURL    string
	SetBackendURL func(string, app.Context)

	OIDCIssuer    string
	SetOIDCIssuer func(string, app.Context)

	OIDCClientID    string
	SetOIDCClientID func(string, app.Context)

	OIDCRedirectURL    string
	SetOIDCRedirectURL func(string, app.Context)

	Submit func(app.Context)
}

const (
	// Names and IDs
	backendURLName      = "backendURLName"
	oidcIssuerName      = "oidcIssuer"
	oidcClientIDName    = "oidcClientID"
	oidcRedirectURLName = "oidcRedirectURL"

	// Placeholders
	backendURLPlaceholder      = "ws://localhost:15124"
	oidcIssuerPlaceholder      = "https://pojntfx.eu.auth0.com/"
	oidcRedirectURLPlaceholder = "http://localhost:15125/"
)

func (c *SetupForm) Render() app.UI {
	return app.Form().
		Class("pf-c-form").
		Body(
			// Error display
			app.If(c.Error != nil, app.P().
				Class("pf-c-form__helper-text pf-m-error").
				Aria("live", "polite").
				Body(
					app.Span().
						Class("pf-c-form__helper-text-icon").
						Body(
							app.I().
								Class("fas fa-exclamation-circle").
								Aria("hidden", true),
						),
					app.Text(c.ErrorMessage),
				),
			),
			// Backend URL Input
			&FormGroup{
				Label: app.
					Label().
					For(backendURLName).
					Class("pf-c-form__label").
					Body(
						app.
							Span().
							Class("pf-c-form__label-text").
							Text("Backend URL"),
					),
				Input: app.
					Input().
					Name(backendURLName).
					ID(backendURLName).
					Type("url").
					Required(true).
					Placeholder(backendURLPlaceholder).
					Class("pf-c-form-control").
					Aria("invalid", c.Error != nil).
					Value(c.BackendURL).
					OnInput(func(ctx app.Context, e app.Event) {
						c.SetBackendURL(ctx.JSSrc().Get("value").String(), ctx)
					}),
				Required: true,
			},
			// OIDC Issuer Input
			&FormGroup{
				Label: app.
					Label().
					For(oidcIssuerName).
					Class("pf-c-form__label").
					Body(
						app.
							Span().
							Class("pf-c-form__label-text").
							Text("OIDC Issuer"),
					),
				Input: app.
					Input().
					Name(oidcIssuerName).
					ID(oidcIssuerName).
					Type("url").
					Required(true).
					Placeholder(oidcIssuerPlaceholder).
					Class("pf-c-form-control").
					Aria("invalid", c.Error != nil).
					Value(c.OIDCIssuer).
					OnInput(func(ctx app.Context, e app.Event) {
						c.SetOIDCIssuer(ctx.JSSrc().Get("value").String(), ctx)
					}),
				Required: true,
			},
			// OIDC Client ID
			&FormGroup{
				Label: app.
					Label().
					For(oidcClientIDName).
					Class("pf-c-form__label").
					Body(
						app.
							Span().
							Class("pf-c-form__label-text").
							Text("OIDC Client ID"),
					),
				Input: app.
					Input().
					Name(oidcClientIDName).
					ID(oidcClientIDName).
					Type("text").
					Required(true).
					Class("pf-c-form-control").
					Aria("invalid", c.Error != nil).
					Value(c.OIDCClientID).
					OnInput(func(ctx app.Context, e app.Event) {
						c.SetOIDCClientID(ctx.JSSrc().Get("value").String(), ctx)
					}),
				Required: true,
			},
			// OIDC Redirect URL
			&FormGroup{
				Label: app.
					Label().
					For(oidcRedirectURLName).
					Class("pf-c-form__label").
					Body(
						app.
							Span().
							Class("pf-c-form__label-text").
							Text("OIDC Redirect URL"),
					),
				Input: app.
					Input().
					Name(oidcRedirectURLName).
					ID(oidcRedirectURLName).
					Type("url").
					Required(true).
					Placeholder(oidcRedirectURLPlaceholder).
					Class("pf-c-form-control").
					Aria("invalid", c.Error != nil).
					Value(c.OIDCRedirectURL).
					OnInput(func(ctx app.Context, e app.Event) {
						c.SetOIDCRedirectURL(ctx.JSSrc().Get("value").String(), ctx)
					}),
				Required: true,
			},
			// Configuration Apply Trigger
			app.Div().
				Class("pf-c-form__group pf-m-action").
				Body(
					app.
						Button().
						Type("submit").
						Class("pf-c-button pf-m-primary pf-m-block").
						Text("Log in"),
				),
		).OnSubmit(func(ctx app.Context, e app.Event) {
		e.PreventDefault()

		c.Submit(ctx)
	})
}
