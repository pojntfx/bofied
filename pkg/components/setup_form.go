package components

import "github.com/maxence-charriere/go-app/v10/pkg/app"

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
	backendURLPlaceholder      = "http://localhost:15256"
	oidcIssuerPlaceholder      = "https://pojntfx.eu.auth0.com/"
	oidcRedirectURLPlaceholder = "http://localhost:15255/"
)

func (c *SetupForm) Render() app.UI {
	return app.Form().
		Class("pf-v6-c-form").
		Body(
			// Error display
			app.If(c.Error != nil, func() app.UI {
				return app.P().
					Class("pf-v6-c-form__helper-text pf-m-error").
					Aria("live", "polite").
					Body(
						app.Span().
							Class("pf-v6-c-form__helper-text-icon").
							Body(
								app.I().
									Class("fas fa-exclamation-circle").
									Aria("hidden", true),
							),
						app.Text(c.ErrorMessage),
					)
			},
			),
			// Backend URL Input
			&FormGroup{
				Label: app.
					Label().
					For(backendURLName).
					Class("pf-v6-c-form__label").
					Body(
						app.
							Span().
							Class("pf-v6-c-form__label-text").
							Text("Backend URL"),
					),
				Input: app.
					Input().
					Name(backendURLName).
					ID(backendURLName).
					Type("url").
					Required(true).
					Placeholder(backendURLPlaceholder).
					Class("pf-v6-c-form-control").
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
					Class("pf-v6-c-form__label").
					Body(
						app.
							Span().
							Class("pf-v6-c-form__label-text").
							Text("OIDC Issuer"),
					),
				Input: app.
					Input().
					Name(oidcIssuerName).
					ID(oidcIssuerName).
					Type("url").
					Required(true).
					Placeholder(oidcIssuerPlaceholder).
					Class("pf-v6-c-form-control").
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
					Class("pf-v6-c-form__label").
					Body(
						app.
							Span().
							Class("pf-v6-c-form__label-text").
							Text("OIDC Client ID"),
					),
				Input: app.
					Input().
					Name(oidcClientIDName).
					ID(oidcClientIDName).
					Type("text").
					Required(true).
					Class("pf-v6-c-form-control").
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
					Class("pf-v6-c-form__label").
					Body(
						app.
							Span().
							Class("pf-v6-c-form__label-text").
							Text("OIDC Redirect URL"),
					),
				Input: app.
					Input().
					Name(oidcRedirectURLName).
					ID(oidcRedirectURLName).
					Type("url").
					Required(true).
					Placeholder(oidcRedirectURLPlaceholder).
					Class("pf-v6-c-form-control").
					Aria("invalid", c.Error != nil).
					Value(c.OIDCRedirectURL).
					OnInput(func(ctx app.Context, e app.Event) {
						c.SetOIDCRedirectURL(ctx.JSSrc().Get("value").String(), ctx)
					}),
				Required: true,
			},
			// Configuration Apply Trigger
			app.Div().
				Class("pf-v6-c-form__group pf-m-action").
				Body(
					app.
						Button().
						Type("submit").
						Class("pf-v6-c-button pf-m-primary pf-m-block").
						Text("Log in"),
				),
		).OnSubmit(func(ctx app.Context, e app.Event) {
		e.PreventDefault()

		c.Submit(ctx)
	})
}
