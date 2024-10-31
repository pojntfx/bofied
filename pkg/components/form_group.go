package components

import "github.com/maxence-charriere/go-app/v10/pkg/app"

type FormGroup struct {
	app.Compo

	Required     bool
	Label        app.UI
	Input        app.UI
	NoTopPadding bool
}

func (c *FormGroup) Render() app.UI {
	return app.Div().
		Class("pf-v6-c-form__group").
		Body(
			app.Div().
				Class(func() string {
					classes := "pf-v6-c-form__group-label"
					if c.NoTopPadding {
						classes += " pf-m-no-padding-top"
					}

					return classes
				}()).
				Body(
					c.Label,
					app.If(c.Required,
						func() app.UI {
							return app.Span().
								Class("pf-v6-c-form__label-required").
								Aria("hidden", true).
								Text("*")
						},
					),
				),
			app.Div().
				Class("pf-v6-c-form__group-control").
				Body(
					app.
						Span().
						Class("pf-v6-c-form-control").
						Body(
							c.Input,
						),
				),
		)
}
