package components

import "github.com/maxence-charriere/go-app/v9/pkg/app"

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
			classes := "pf-c-expandable-section"

			if c.Open {
				classes += " pf-m-expanded"
			}

			return classes
		}()).
		Body(
			app.Button().
				Type("button").
				Class("pf-c-expandable-section__toggle").
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
						Class("pf-c-expandable-section__toggle-icon").
						Body(
							app.I().
								Class("fas fa-angle-right").
								Aria("hidden", true),
						),
					app.Span().
						Class("pf-c-expandable-section__toggle-text").
						Text(c.Title),
				),
			&Controlled{
				Component: app.Div().
					Class("pf-c-expandable-section__content").
					Hidden(!c.Open).
					Body(c.Body...),
				Properties: map[string]interface{}{
					"hidden": !c.Open,
				},
			},
		)
}
