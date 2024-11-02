package components

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
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
		Class("pf-v6-c-switch").
		For(c.ID).
		Body(
			app.Input().
				Class("pf-v6-c-switch__input").
				Type("checkbox").
				ID(c.ID).
				Aria("labelledby", c.ID+"-on").
				Name(c.ID).
				Checked(c.Open).
				OnInput(func(ctx app.Context, e app.Event) {
					c.ToggleOpen()
				}),
			app.Span().
				Class("pf-v6-c-switch__toggle"),
			app.If(
				c.OnMessage != "" && c.Open,
				func() app.UI {
					return app.Span().
						Class("pf-v6-c-switch__label pf-m-on").
						ID(c.ID+"-on").
						Aria("hidden", true).
						Text(c.OnMessage)
				},
			),
			app.If(
				c.OffMessage != "" && !c.Open,
				func() app.UI {
					return app.Span().
						Class("pf-v6-c-switch__label pf-m-off").
						ID(c.ID+"-off").
						Aria("hidden", true).
						Text(c.OffMessage)
				},
			),
		)
}
