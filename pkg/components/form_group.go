package components

import "github.com/maxence-charriere/go-app/v9/pkg/app"

type FormGroup struct {
	app.Compo

	Required     bool
	Label        app.UI
	Input        app.UI
	NoTopPadding bool
}

func (c *FormGroup) Render() app.UI {
	return app.Div().
		Class("pf-c-form__group").
		Body(
			app.Div().
				Class(func() string {
					classes := "pf-c-form__group-label"
					if c.NoTopPadding {
						classes += " pf-m-no-padding-top"
					}

					return classes
				}()).
				Body(
					c.Label,
					app.If(c.Required,
						app.Span().
							Class("pf-c-form__label-required").
							Aria("hidden", true).
							Text("*"),
					),
				),
			app.Div().
				Class("pf-c-form__group-control").
				Body(c.Input),
		)
}
