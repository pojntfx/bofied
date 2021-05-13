package components

import "github.com/maxence-charriere/go-app/v9/pkg/app"

type TextEditorWrapper struct {
	app.Compo

	Title    string
	HelpLink string

	Children app.UI

	Error            error
	ErrorDescription string
	Ignore           func()
}

func (c *TextEditorWrapper) Render() app.UI {
	return app.Div().
		Class("pf-c-card pf-u-h-100").
		Body(
			app.Div().
				Class("pf-c-card__header pf-x-m-gap-md").
				Body(
					app.If(
						c.HelpLink != "",
						app.Div().
							Class("pf-c-card__actions").
							Body(
								app.A().
									Class("pf-c-button pf-m-plain").
									Aria("label", "Help").
									Target("_blank").
									Href(c.HelpLink).
									Body(
										app.I().
											Class("fas fa-question-circle").
											Aria("hidden", true),
									),
							),
					),
					app.Div().
						Class("pf-c-card__title").
						Text(c.Title),
				),
			app.Div().
				Class("pf-c-card__body").
				Body(c.Children),
			app.If(
				c.Error != nil,
				app.Div().
					Class("pf-c-card__footer").
					Body(
						app.Div().
							Class("pf-c-alert pf-m-danger pf-m-inline").
							Aria("label", "Error alert").
							Body(
								app.Div().
									Class("pf-c-alert__icon").
									Body(
										app.I().
											Class("fas fa-fw fa-exclamation-circle").
											Aria("hidden", true),
									),
								app.P().
									Class("pf-c-alert__title").
									Body(
										app.
											Strong().
											Body(
												app.Span().
													Class("pf-screen-reader").
													Text(c.ErrorDescription+":"),
												app.Text(c.ErrorDescription),
											),
									),
								app.Div().
									Class("pf-c-alert__action").
									Body(
										app.Button().
											Class("pf-c-button pf-m-plain").
											Type("button").
											Aria("label", "Button to ignore the error").
											OnClick(func(ctx app.Context, e app.Event) {
												c.Ignore()
											}).
											Body(
												app.I().
													Class("fas fa-times").
													Aria("hidden", true),
											),
									),
								app.Div().
									Class("pf-c-alert__description").
									Body(
										app.Code().Text(c.Error),
									),
							),
					),
			),
		)
}
