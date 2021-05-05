package components

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type Modal struct {
	app.Compo

	Open  bool
	Close func()

	ID      string
	Classes string

	Title  string
	Body   []app.UI
	Footer []app.UI

	Large   bool
	Overlay bool

	Nested bool
}

func (c *Modal) Render() app.UI {
	return app.Div().
		Class(func() string {
			classes := "pf-c-backdrop"

			if c.Classes != "" {
				classes += " " + c.Classes
			}

			if !c.Open {
				classes += " pf-u-display-none"
			}

			if c.Overlay {
				classes += " pf-x-m-modal-overlay"
			}

			if c.Nested {
				classes += " pf-x-c-backdrop--nested"
			}

			return classes
		}()).
		Body(
			app.Div().
				Class("pf-l-bullseye").
				Body(
					app.Div().
						Class(func() string {
							classes := "pf-c-modal-box"
							if c.Large {
								classes += " pf-m-lg"
							} else {
								classes += " pf-m-sm"
							}

							return classes
						}()).
						Aria("modal", true).
						Aria("labelledby", c.ID).
						Body(
							app.Button().
								Class("pf-c-button pf-m-plain").
								Type("button").
								Aria("label", "Close dialog").
								OnClick(func(ctx app.Context, e app.Event) {
									c.Close()
								}).
								Body(
									app.I().
										Class("fas fa-times").
										Aria("hidden", true),
								),
							app.Header().
								Class("pf-c-modal-box__header").
								Body(
									app.H1().
										Class("pf-c-modal-box__title").
										ID(c.ID).
										Text(c.Title),
								),
							app.Div().
								Class("pf-c-modal-box__body").
								Body(c.Body...),
							app.If(
								c.Footer != nil,
								app.Footer().
									Class("pf-c-modal-box__footer").
									Body(c.Footer...),
							),
						),
				),
		)
}

func (c *Modal) OnMount(ctx app.Context) {
	app.Window().AddEventListener("keyup", func(ctx app.Context, e app.Event) {
		if e.Get("key").String() == "Escape" {
			c.Close()
		}
	})
}
