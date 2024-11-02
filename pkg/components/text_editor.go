package components

import (
	"strings"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type TextEditor struct {
	app.Compo

	Content    string
	SetContent func(string)

	Format  func()
	Refresh func()
	Save    func()

	Language       string
	VariableHeight bool
}

func (c *TextEditor) Render() app.UI {
	return app.Div().
		Class(func() string {
			classes := "pf-v6-c-code-editor pf-v6-u-h-100 pf-v6-u-display-flex pf-v6-u-flex-direction-column"
			if c.SetContent == nil {
				classes += " pf-m-read-only"
			}

			return classes
		}()).
		Body(
			app.Div().
				Class("pf-v6-c-code-editor__header").
				Body(
					app.Div().
						Class("pf-v6-c-code-editor__header-content").
						Body(
							app.Div().
								Class("pf-v6-c-code-editor__controls").
								Body(
									app.If(
										c.Format != nil,
										func() app.UI {
											return app.Button().
												Class("pf-v6-c-button pf-m-plain").
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
												)
										},
									),
									app.If(
										c.Refresh != nil,
										func() app.UI {
											return app.Button().
												Class("pf-v6-c-button pf-m-plain").
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
												)
										},
									),
									app.If(
										c.Save != nil,
										func() app.UI {
											return app.Button().
												Class("pf-v6-c-button pf-m-plain").
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
												)
										},
									),
								),
							app.Div().
								Class("pf-v6-c-code-editor__header-main"),
						),
					app.If(
						c.Language != "",
						func() app.UI {
							return app.Div().
								Class("pf-v6-c-code-editor__tab").
								Body(
									app.Span().
										Class("pf-v6-c-code-editor__tab-icon").
										Body(
											app.I().
												Class("fas fa-code").
												Aria("hidden", true),
										),
									app.Span().
										Class("pf-v6-c-code-editor__tab-text").
										Text(c.Language),
								)
						},
					),
				),
			app.Textarea().
				Class(func() string {
					classes := "pf-v6-c-code-editor__main pf-v6-u-w-100 pf-v6-x-u-resize-none pf-v6-u-p-sm pf-v6-u-p-sm pf-v6-u-flex-fill"
					if c.VariableHeight {
						classes += " pf-v6-x-m-overflow-y-hidden"
					}

					return classes
				}()).
				Rows(func() int {
					if c.VariableHeight {
						return strings.Count(c.Content, "\n") + 1 // Trailing newline
					}

					return 25
				}()).
				ReadOnly(c.SetContent == nil).
				OnInput(func(ctx app.Context, e app.Event) {
					c.SetContent(ctx.JSSrc().Get("value").String())
				}).
				Text(c.Content),
		)
}
