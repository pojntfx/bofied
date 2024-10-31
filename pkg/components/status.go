package components

import "github.com/maxence-charriere/go-app/v10/pkg/app"

type Status struct {
	app.Compo

	Error       error
	ErrorText   string
	Recover     func(app.Context)
	RecoverText string
	Ignore      func()
}

func (c *Status) Render() app.UI {
	// Display the error message if error != nil
	errorMessage := ""
	if c.Error != nil {
		errorMessage = c.Error.Error()
	}

	return app.If(c.Error != nil, func() app.UI {
		return app.Div().
			Class("pf-v6-c-alert pf-m-danger").
			Aria("label", c.ErrorText).
			Body(
				app.Div().
					Class("pf-v6-c-alert__icon").
					Body(
						app.I().
							Class("fas fa-fw fa-exclamation-circle").
							Aria("hidden", true),
					),
				app.P().
					Class("pf-v6-c-alert__title").
					Body(
						app.Strong().Body(
							app.Span().
								Class("pf-screen-reader").
								Text(c.ErrorText),
						),
						app.Text(c.ErrorText),
					),
				app.Div().
					Class("pf-v6-c-alert__action").
					Body(
						app.Button().
							Class("pf-v6-c-button pf-m-plain").
							Aria("label", "Ignore error").
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
					Class("pf-v6-c-alert__description").
					Body(
						app.P().Body(
							app.Code().
								Text(errorMessage),
						),
					),
				app.If(c.Recover != nil,
					func() app.UI {
						return app.Div().
							Class("pf-v6-c-alert__action-group").
							Body(
								app.Button().
									Class("pf-v6-c-button pf-m-link pf-m-inline").
									Type("button").
									OnClick(func(ctx app.Context, e app.Event) {
										c.Recover(ctx)
									}).
									Text(c.RecoverText),
							)
					},
				),
			)
	}).Else(func() app.UI {
		return app.Span()
	})
}
