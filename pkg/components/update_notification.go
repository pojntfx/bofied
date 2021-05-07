package components

import "github.com/maxence-charriere/go-app/v9/pkg/app"

type UpdateNotification struct {
	app.Compo

	UpdateTitle       string
	UpdateDescription string

	StartUpdateText  string
	IgnoreUpdateText string

	updateAvailable bool
	updateIgnored   bool
}

func (c *UpdateNotification) Render() app.UI {
	return app.If(
		c.updateAvailable && !c.updateIgnored,
		app.Li().
			Class("pf-c-alert-group__item").
			Body(
				app.Div().
					Class("pf-c-alert pf-m-info").
					Aria("label", c.UpdateTitle).
					Body(
						app.Div().
							Class("pf-c-alert__icon").
							Body(
								app.I().
									Class("fas fa-fw fa-bell").
									Aria("hidden", true),
							),
						app.P().
							Class("pf-c-alert__title").
							Body(
								app.Strong().Body(
									app.Span().
										Class("pf-screen-reader").
										Text(c.UpdateTitle),
								),
								app.Text(c.UpdateTitle),
							),
						app.Div().
							Class("pf-c-alert__action").
							Body(
								app.Button().
									Class("pf-c-button pf-m-plain").
									Aria("label", c.IgnoreUpdateText).
									OnClick(func(ctx app.Context, e app.Event) {
										c.updateIgnored = true
									}).
									Body(
										app.I().
											Class("fas fa-times").
											Aria("hidden", true),
									),
							),
						app.If(
							c.UpdateDescription != "",
							app.Div().
								Class("pf-c-alert__description").
								Body(
									app.P().Text(c.UpdateDescription),
								),
						),
						app.Div().
							Class("pf-c-alert__action-group").
							Body(
								app.Button().
									Class("pf-c-button pf-m-link pf-m-inline").
									Type("button").
									OnClick(func(ctx app.Context, e app.Event) {
										ctx.Reload()
									}).
									Body(
										app.Span().Class("pf-c-button__icon pf-m-start").Body(
											app.I().Class("fas fas fa-arrow-up").Aria("hidden", true),
										),
										app.Text(c.StartUpdateText),
									),
								app.Button().
									Class("pf-c-button pf-m-link pf-m-inline").
									Type("button").
									OnClick(func(ctx app.Context, e app.Event) {
										c.updateIgnored = true
									}).
									Body(
										app.Span().Class("pf-c-button__icon pf-m-start").Body(
											app.I().Class("fas fa-ban").Aria("hidden", true),
										),
										app.Text(c.IgnoreUpdateText),
									),
							),
					),
			)).Else(app.Span())
}

func (c *UpdateNotification) OnMount(ctx app.Context) {
	c.updateAvailable = ctx.AppUpdateAvailable()
}
