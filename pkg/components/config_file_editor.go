package components

import (
	"github.com/maxence-charriere/go-app/v8/pkg/app"
)

type ConfigFileEditor struct {
	app.Compo

	ConfigFile    string
	SetConfigFile func(string)

	FormatConfigFile  func()
	RefreshConfigFile func()
	SaveConfigFile    func()

	Error  error
	Ignore func()
}

func (c *ConfigFileEditor) Render() app.UI {
	return app.Div().
		Class("pf-c-card").
		Body(
			app.Div().
				Class("pf-c-card__title").
				Text("Config"),
			app.Div().
				Class("pf-c-card__body").
				Body(
					app.Div().
						Class("pf-c-code-editor").
						Body(
							app.Div().
								Class("pf-c-code-editor__header").
								Body(
									app.Div().
										Class("pf-c-code-editor__controls").
										Body(
											app.Button().
												Class("pf-c-button pf-m-control").
												Type("button").
												Aria("label", "Format").
												Title("Format").
												OnClick(func(ctx app.Context, e app.Event) {
													c.FormatConfigFile()
												}).
												Body(
													app.I().
														Class("fas fa-align-left").
														Aria("hidden", true),
												),
											app.Button().
												Class("pf-c-button pf-m-control").
												Type("button").
												Aria("label", "Refresh").
												Title("Refresh").
												OnClick(func(ctx app.Context, e app.Event) {
													c.RefreshConfigFile()
												}).
												Body(
													app.I().
														Class("fas fas fa-sync").
														Aria("hidden", true),
												),
											app.Button().
												Class("pf-c-button pf-m-control").
												Type("button").
												Aria("label", "Save").
												Title("Save").
												OnClick(func(ctx app.Context, e app.Event) {
													c.SaveConfigFile()
												}).
												Body(
													app.I().
														Class("fas fas fa-save").
														Aria("hidden", true),
												),
										),
									app.Div().
										Class("pf-c-code-editor__tab").
										Body(
											app.Span().
												Class("pf-c-code-editor__tab-icon").
												Body(
													app.I().
														Class("fas fa-code").
														Aria("hidden", true),
												),
											app.Span().
												Class("pf-c-code-editor__tab-text").
												Text("Go"),
										),
								),
							&Controlled{
								Component: app.Textarea().
									Class("pf-c-code-editor__main pf-u-w-100 pf-x-u-resize-vertical").
									Rows(25).
									OnInput(func(ctx app.Context, e app.Event) {
										c.SetConfigFile(ctx.JSSrc.Get("value").String())
									}).
									Text(c.ConfigFile),
								Properties: map[string]interface{}{
									"value": c.ConfigFile,
								},
							},
						),
				),
			app.If(
				c.Error != nil,
				app.Div().
					Class("pf-c-card__footer").
					Body(
						app.Div().
							Class("pf-c-alert pf-m-danger pf-m-inline").
							Aria("label", "Syntax error alert").
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
													Text("Syntax error:"),
												app.Text("Syntax Error"),
											),
									),
								app.Div().
									Class("pf-c-alert__action").
									Body(
										app.Button().
											Class("pf-c-button pf-m-plain").
											Type("button").
											Aria("label", "Button to ignore the syntax error").
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
