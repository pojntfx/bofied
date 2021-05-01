package components

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type Switch struct {
	app.Compo

	ID string

	Open       bool
	ToggleOpen func()

	OnMessage  string
	OffMessage string
}

func (c *Switch) Render() app.UI {
	return app.Label().
		Class("pf-c-switch").
		For(c.ID).
		Body(
			&Controlled{
				Component: app.Input().
					Class("pf-c-switch__input").
					Type("checkbox").
					ID(c.ID).
					Aria("labelledby", c.ID+"-on").
					Name(c.ID).
					Checked(c.Open).
					OnInput(func(ctx app.Context, e app.Event) {
						c.ToggleOpen()
					}),
				Properties: map[string]interface{}{
					"checked": c.Open,
				},
			},
			app.Span().
				Class("pf-c-switch__toggle"),
			app.If(
				c.OnMessage != "",
				app.Span().
					Class("pf-c-switch__label pf-m-on").
					ID(c.ID+"-on").
					Aria("hidden", true).
					Text(c.OnMessage),
			),
			app.If(
				c.OffMessage != "",
				app.Span().
					Class("pf-c-switch__label pf-m-off").
					ID(c.ID+"-off").
					Aria("hidden", true).
					Text(c.OffMessage),
			),
		)
}
