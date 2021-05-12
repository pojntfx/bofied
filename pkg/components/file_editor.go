package components

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type FileEditor struct {
	app.Compo

	Content    string
	SetContent func(string)

	Format  func()
	Refresh func()
	Save    func()

	Language string
	HelpLink string

	Error            error
	ErrorDescription string
	Ignore           func()
}

func (c *FileEditor) Render() app.UI {
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
						Text("Config"),
				),
			app.Div().
				Class("pf-c-card__body").
				Body(
					app.Div().
						Class("pf-c-code-editor pf-u-h-100 pf-u-display-flex pf-u-flex-direction-column").
						Body(
							app.Div().
								Class("pf-c-code-editor__header").
								Body(
									app.Div().
										Class("pf-c-code-editor__controls").
										Body(
											app.If(
												c.Format != nil,
												app.Button().
													Class("pf-c-button pf-m-control").
													Type("button").
													Aria("label", "Format").
													Title("Format").
													OnClick(func(ctx app.Context, e app.Event) {
														c.Format()
													}).
													Body(
														app.I().
															Class("fas fa-align-left").
															Aria("hidden", true),
													),
											),
											app.If(
												c.Refresh != nil,
												app.Button().
													Class("pf-c-button pf-m-control").
													Type("button").
													Aria("label", "Refresh").
													Title("Refresh").
													OnClick(func(ctx app.Context, e app.Event) {
														c.Refresh()
													}).
													Body(
														app.I().
															Class("fas fas fa-sync").
															Aria("hidden", true),
													),
											),
											app.If(
												c.Save != nil,
												app.Button().
													Class("pf-c-button pf-m-control").
													Type("button").
													Aria("label", "Save").
													Title("Save").
													OnClick(func(ctx app.Context, e app.Event) {
														c.Save()
													}).
													Body(
														app.I().
															Class("fas fas fa-save").
															Aria("hidden", true),
													),
											),
										),
									app.If(
										c.Language != "",
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
													Text(c.Language),
											),
									),
								),
							&Controlled{
								Component: app.Textarea().
									Class("pf-c-code-editor__main pf-u-w-100 pf-x-u-resize-none pf-u-p-sm pf-u-p-sm pf-u-flex-fill").
									Rows(25).
									OnInput(func(ctx app.Context, e app.Event) {
										c.SetContent(ctx.JSSrc().Get("value").String())
									}).
									Text(c.Content),
								Properties: map[string]interface{}{
									"value": c.Content,
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
