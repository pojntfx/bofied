package components

import "github.com/maxence-charriere/go-app/v10/pkg/app"

type ExpandableSection struct {
	app.Compo

	Open        bool
	OnToggle    func()
	Title       string
	ClosedTitle string
	OpenTitle   string
	Body        []app.UI
}

func (c *ExpandableSection) Render() app.UI {
	return app.Div().
		Class(func() string {
			classes := "pf-v6-c-expandable-section"

			if c.Open {
				classes += " pf-m-expanded"
			}

			return classes
		}()).
		Body(
			app.Div().
				Class("pf-v6-c-expandable-section__toggle").
				Body(
					app.Button().
						Type("button").
						Class("pf-v6-c-button pf-m-link").
						Aria("label", func() string {
							message := c.ClosedTitle

							if c.Open {
								message = c.OpenTitle
							}

							return message
						}()).
						Aria("expanded", c.Open).
						OnClick(func(ctx app.Context, e app.Event) {
							c.OnToggle()
						}).
						Body(
							app.Span().
								Class("pf-v6-c-button__icon pf-m-start").
								Body(
									app.I().
										Class("fas fa-angle-right").
										Aria("hidden", true),
								),
							app.Span().
								Class("pf-v6-c-button__text").
								Text(c.Title),
						),
				),
			app.Div().
				Class("pf-v6-c-expandable-section__content").
				Hidden(!c.Open).
				Body(c.Body...),
		)
}
