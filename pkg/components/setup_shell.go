package components

import (
	"github.com/maxence-charriere/go-app/v8/pkg/app"
)

type SetupShell struct {
	app.Compo

	LogoSrc          string
	Title            string
	ShortDescription string
	LongDescription  string
	HelpLink         string
	Links            map[string]string

	BackendURL      string
	OIDCIssuer      string
	OIDCClientID    string
	OIDCRedirectURL string

	SetBackendURL,
	SetOIDCIssuer,
	SetOIDCClientID,
	SetOIDCRedirectURL func(string)
	ApplyConfig func()

	Error error
}

func (c *SetupShell) Render() app.UI {
	// Display the error message if error != nil
	errorMessage := ""
	if c.Error != nil {
		errorMessage = c.Error.Error()
	}

	return app.Div().
		Class("pf-u-h-100").
		Body(
			app.Div().
				Class("pf-c-background-image").
				Body(
					app.Raw(`<svg
  xmlns="http://www.w3.org/2000/svg"
  class="pf-c-background-image__filter"
  width="0"
  height="0"
>
  <filter id="image_overlay">
    <feColorMatrix
      type="matrix"
      values="1 0 0 0 0 1 0 0 0 0 1 0 0 0 0 0 0 0 1 0"
    ></feColorMatrix>
    <feComponentTransfer
      color-interpolation-filters="sRGB"
      result="duotone"
    >
      <feFuncR
        type="table"
        tableValues="0.086274509803922 0.43921568627451"
      ></feFuncR>
      <feFuncG
        type="table"
        tableValues="0.086274509803922 0.43921568627451"
      ></feFuncG>
      <feFuncB
        type="table"
        tableValues="0.086274509803922 0.43921568627451"
      ></feFuncB>
      <feFuncA type="table" tableValues="0 1"></feFuncA>
    </feComponentTransfer>
  </filter>
</svg>`),
				),
			app.Div().Class("pf-c-login").Body(
				app.Div().Class("pf-c-login__container").Body(
					app.Header().Class("pf-c-login__header").Body(
						app.Img().
							Class("pf-c-brand pf-x-c-brand--main").
							Src(c.LogoSrc).
							Alt("Logo"),
					),
					app.Main().Class("pf-c-login__main").Body(
						app.Header().Class("pf-c-login__main-header").Body(
							app.H1().Class("pf-c-title pf-m-3xl").Text(
								c.Title,
							),
							app.P().Class("pf-c-login__main-header-desc").Text(
								c.ShortDescription,
							),
						),
						app.Div().Class("pf-c-login__main-body").Body(
							&SetupForm{
								Error:        c.Error,
								ErrorMessage: errorMessage,

								BackendURL:    c.BackendURL,
								SetBackendURL: c.SetBackendURL,

								OIDCIssuer:    c.OIDCIssuer,
								SetOIDCIssuer: c.SetOIDCIssuer,

								OIDCClientID:    c.OIDCClientID,
								SetOIDCClientID: c.SetOIDCClientID,

								OIDCRedirectURL:    c.OIDCRedirectURL,
								SetOIDCRedirectURL: c.SetOIDCRedirectURL,

								Submit: c.ApplyConfig,
							},
						),
						app.Footer().Class("pf-c-login__main-footer").Body(
							app.Div().Class("pf-c-login__main-footer-band").Body(
								app.P().Class("pf-c-login__main-footer-band-item").Body(
									app.Text("Not sure what to do? "),
									app.A().
										Href(c.HelpLink).
										Target("_blank").
										Text("Get help."),
								),
							),
						),
					),
					app.Footer().Class("pf-c-login__footer").Body(
						app.P().Text(
							c.LongDescription,
						),
						app.Ul().Class("pf-c-list pf-m-inline").Body(
							app.Range(c.Links).Map(func(s string) app.UI {
								return app.Li().Body(
									app.
										A().
										Target("_blank").
										Href(c.Links[s]).
										Text(s),
								)
							}),
						),
					),
				),
			),
		)

}
