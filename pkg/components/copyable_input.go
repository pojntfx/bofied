package components

import "github.com/maxence-charriere/go-app/v10/pkg/app"

type CopyableInput struct {
	app.Compo

	Component app.UI
	ID        string
}

func (c *CopyableInput) Render() app.UI {
	return app.Div().
		Class("pf-v6-c-clipboard-copy").
		Body(
			app.Div().
				Class("pf-v6-c-clipboard-copy__group").
				Body(
					app.Span().
						Class("pf-v6-c-form-control").
						Body(
							c.Component,
						),
					app.Button().
						Class("pf-v6-c-button pf-m-control").
						Type("button").
						Aria("label", "Copy to clipboard").
						Aria("labelledby", c.ID).
						OnClick(func(ctx app.Context, e app.Event) {
							app.Window().JSValue().Get("document").Call("getElementById", c.ID).Call("select")

							app.Window().JSValue().Get("document").Call("execCommand", "copy")
						}).
						Body(
							app.I().
								Class("fas fa-copy").
								Aria("hidden", true),
						),
				),
		)
}
