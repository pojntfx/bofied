package components

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

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
			classes := "pf-v6-c-backdrop"

			if !c.Open {
				classes += " pf-v6-u-display-none"
			}

			return classes
		}()).
		Body(
			app.Div().
				Class("pf-v6-l-bullseye").
				Body(
					app.Div().
						Class("pf-v6-c-modal-box pf-m-lg").
						Aria("role", "dialog").
						Aria("modal", true).
						Aria("labelledby", c.ID).
						Body(
							app.Div().
								Class("pf-v6-c-about-modal-box").
								Body(
									app.Div().
										Class("pf-v6-c-about-modal-box__brand").
										Body(
											app.Img().
												Class("pf-v6-c-about-modal-box__brand-image").
												Src(c.LogoSrc).
												Alt(c.LogoAlt),
										),
									app.Div().
										Class("pf-v6-c-about-modal-box__close").
										Body(
											app.Button().
												Class("pf-v6-c-button pf-m-plain").
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
										Class("pf-v6-c-about-modal-box__header").
										Body(
											app.H1().
												Class("pf-v6-c-title pf-m-4xl").
												ID(c.ID).
												Text(c.Title),
										),
									app.Div().Class("pf-v6-c-about-modal-box__hero"),
									app.Div().
										Class("pf-v6-c-about-modal-box__content").
										Body(
											app.Div().
												Class("pf-v6-c-content").
												Body(
													c.Body,
												),
											app.P().
												Class("pf-v6-c-about-modal-box__strapline").
												Text(c.Footer),
										),
								),
						),
				),
		)
}

func (c *AboutModal) OnMount(ctx app.Context) {
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
