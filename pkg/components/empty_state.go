package components

import "github.com/maxence-charriere/go-app/v10/pkg/app"

type EmptyState struct {
	app.Compo

	Action app.UI
}

func (c *EmptyState) Render() app.UI {
	return app.Div().
		Class("pf-c-empty-state").
		Body(
			app.Div().
				Class("pf-c-empty-state__content").
				Body(
					app.I().
						Class("fas fa-folder-open pf-c-empty-state__icon").
						Aria("hidden", true),
					app.H1().
						Class("pf-c-title pf-m-lg").
						Text("No files or directories here yet"),
					app.Div().
						Class("pf-c-empty-state__body").
						Text("You can add a file or directory to make it available for nodes."),
					app.If(
						c.Action != nil,
						func() app.UI {
							return c.Action
						},
					),
				),
		)
}
