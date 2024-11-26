package components

import (
	"path"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type Breadcrumbs struct {
	app.Compo

	PathComponents []string

	CurrentPath    string
	SetCurrentPath func(string)

	SelectedPath    string
	SetSelectedPath func(string)

	ItemClass string
}

func (c *Breadcrumbs) Render() app.UI {
	return app.Nav().
		Class("pf-v6-c-breadcrumb").
		Aria("label", "Current path").
		Body(
			app.Ol().
				Class("pf-v6-c-breadcrumb__list pf-v6-u-font-weight-bold").
				Body(
					app.Li().
						Class("pf-v6-c-breadcrumb__item", c.ItemClass).
						Body(
							app.Span().
								Class("pf-v6-c-breadcrumb__item-divider").
								Body(
									app.I().
										Class("fas fa-angle-right").
										Aria("hidden", true),
								),
							app.Button().
								Type("button").
								Class("pf-v6-c-breadcrumb__link").
								TabIndex(0).
								OnClick(func(ctx app.Context, e app.Event) {
									c.SetCurrentPath("/")

									c.SetSelectedPath("")
								}).
								Text("Files"),
						),
					app.Range(c.PathComponents).Slice(func(i int) app.UI {
						link := path.Join(append([]string{"/"}, c.PathComponents[:i+1]...)...)

						// The last path part shouldn't be marked as a link
						classes := "pf-v6-c-breadcrumb__link"
						if i == len(c.PathComponents)-1 {
							classes += " pf-m-current"
						}

						return app.Li().
							Class("pf-v6-c-breadcrumb__item", c.ItemClass).
							Body(
								app.Span().
									Class("pf-v6-c-breadcrumb__item-divider").
									Body(
										app.I().
											Class("fas fa-angle-right").
											Aria("hidden", true),
									),
								app.If(
									// The last path part shouldn't be an action
									i == len(c.PathComponents)-1,
									func() app.UI {
										return app.A().
											Class(classes).
											Text(c.PathComponents[i])
									},
								).Else(
									func() app.UI {
										return app.Button().
											Type("button").
											Class(classes).
											OnClick(func(ctx app.Context, e app.Event) {
												c.SetCurrentPath(link)

												c.SetSelectedPath("")
											}).
											Text(c.PathComponents[i])
									},
								),
							)
					}),
				),
		)
}
