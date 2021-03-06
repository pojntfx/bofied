package components

import "github.com/maxence-charriere/go-app/v9/pkg/app"

type AboutModal struct {
	app.Compo

	Open  bool
	Close func()

	ID string

	LogoSrc string
	LogoAlt string
	Title   string

	Body   app.UI
	Footer string
}

func (c *AboutModal) Render() app.UI {
	return app.Div().
		Class(func() string {
			classes := "pf-c-backdrop"

			if !c.Open {
				classes += " pf-u-display-none"
			}

			return classes
		}()).
		Body(
			app.Div().
				Class("pf-l-bullseye").
				Body(
					app.Div().
						Class("pf-c-about-modal-box").
						Aria("role", "dialog").
						Aria("modal", true).
						Aria("labelledby", c.ID).
						Body(
							app.Div().
								Class("pf-c-about-modal-box__brand").
								Body(
									app.Img().
										Class("pf-c-about-modal-box__brand-image").
										Src(c.LogoSrc).
										Alt(c.LogoAlt),
								),
							app.Div().
								Class("pf-c-about-modal-box__close").
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
								),
							app.Div().
								Class("pf-c-about-modal-box__header").
								Body(
									app.H1().
										Class("pf-c-title pf-m-4xl").
										ID(c.ID).
										Text(c.Title),
								),
							app.Div().Class("pf-c-about-modal-box__hero"),
							app.Div().
								Class("pf-c-about-modal-box__content").
								Body(
									app.Div().
										Class("pf-c-content").
										Body(
											app.Dl().
												Class("pf-c-content").
												Body(c.Body),
										),
									app.P().
										Class("pf-c-about-modal-box__strapline").
										Text(c.Footer),
								),
						),
				),
		)
}

func (c *AboutModal) OnMount(ctx app.Context) {
	app.Window().AddEventListener("keyup", func(ctx app.Context, e app.Event) {
		if e.Get("key").String() == "Escape" {
			c.Close()
		}
	})
}
