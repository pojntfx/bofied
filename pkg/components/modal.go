package components

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
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

	Large        bool
	PaddedBottom bool
	Overlay      bool

	Nested bool
}

func (c *Modal) Render() app.UI {
	return app.Div().
		Class(func() string {
			classes := "pf-v6-c-backdrop"

			if c.Classes != "" {
				classes += " " + c.Classes
			}

			if !c.Open {
				classes += " pf-v6-u-display-none"
			}

			if c.Overlay {
				classes += " pf-v6-x-m-modal-overlay"
			}

			if c.Nested {
				classes += " pf-v6-x-c-backdrop--nested"
			}

			return classes
		}()).
		Body(
			app.Div().
				Class("pf-v6-l-bullseye").
				Body(
					app.Div().
						Class(func() string {
							classes := "pf-v6-c-modal-box"
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
							app.Div().
								Class("pf-v6-c-modal-box__close").
								Body(
									app.Button().
										Class("pf-v6-c-button pf-m-plain").
										Type("button").
										Aria("label", "Close dialog").
										OnClick(func(ctx app.Context, e app.Event) {
											c.Close()
										}).
										Body(
											app.Span().
												Class("pf-v6-c-button__icon").
												Body(
													app.I().
														Class("fas fa-times").
														Aria("hidden", true),
												),
										),
								),
							app.Header().
								Class("pf-v6-c-modal-box__header").
								Body(
									app.H1().
										Class("pf-v6-c-modal-box__title").
										ID(c.ID).
										Text(c.Title),
								),
							app.Div().
								Class(func() string {
									classes := "pf-v6-c-modal-box__body"
									if c.PaddedBottom {
										classes += " pf-v6-u-pb-md"
									}

									return classes
								}()).
								Body(c.Body...),
							app.If(
								c.Footer != nil,
								func() app.UI {
									return app.Footer().
										Class("pf-v6-c-modal-box__footer").
										Body(c.Footer...)
								},
							),
						),
				),
		)
}

func (c *Modal) OnMount(ctx app.Context) {
	app.Window().Call("addEventListener", "keyup", app.FuncOf(func(this app.Value, args []app.Value) any {
		ctx.Async(func() {
			if len(args) > 0 && args[0].Get("key").String() == "Escape" {
				c.Close()

				ctx.Update()
			}
		})

		return nil
	}))
}
