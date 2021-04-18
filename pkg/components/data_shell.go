package components

import (
	"github.com/maxence-charriere/go-app/v8/pkg/app"
	"github.com/pojntfx/liwasc/pkg/components"
)

type DataShell struct {
	app.Compo

	AuthorizedWebDAVURL string
	ConfigFile          string

	SetConfigFile      func(string)
	ValidateConfigFile func()
	SaveConfigFile     func()

	Error   error
	Recover func()
	Ignore  func()
}

func (c *DataShell) Render() app.UI {
	if c.Error != nil {
		return app.Div().
			Body(
				app.Code().Text(c.Error),
				app.Button().
					OnClick(func(ctx app.Context, e app.Event) {
						c.Recover()
					}).
					Text("Recover"),
				app.Button().
					OnClick(func(ctx app.Context, e app.Event) {
						c.Ignore()
					}).
					Text("Ignore"),
			)
	}

	return app.Div().Body(
		app.Input().
			ReadOnly(true).
			Value(
				c.AuthorizedWebDAVURL,
			),
		&components.Controlled{
			Component: app.Textarea().
				OnInput(func(ctx app.Context, e app.Event) {
					c.SetConfigFile(ctx.JSSrc.Get("value").String())
				}).
				Text(
					c.ConfigFile,
				),
			Properties: map[string]interface{}{
				"value": c.ConfigFile,
			},
		},
		app.Button().
			OnClick(func(ctx app.Context, e app.Event) {
				c.ValidateConfigFile()
			}).
			Text("Validate"),
		app.Button().
			OnClick(func(ctx app.Context, e app.Event) {
				c.SaveConfigFile()
			}).
			Text("Save"),
	)
}
