package components

import (
	"github.com/maxence-charriere/go-app/v8/pkg/app"
	"github.com/pojntfx/liwasc/pkg/components"
)

type ConfigFileEditor struct {
	app.Compo

	ConfigFile    string
	SetConfigFile func(string)

	FormatConfigFile  func()
	RefreshConfigFile func()
	SaveConfigFile    func()

	Error   error
	Recover func()
	Ignore  func()
}

func (c *ConfigFileEditor) Render() app.UI {
	return app.Div().
		Body(
			app.H2().
				Text("Config"),
			app.Div().
				Body(
					app.Button().
						OnClick(func(ctx app.Context, e app.Event) {
							c.FormatConfigFile()
						}).
						Text("Format"),
					app.Button().
						OnClick(func(ctx app.Context, e app.Event) {
							c.RefreshConfigFile()
						}).
						Text("Refresh"),
					app.Button().
						OnClick(func(ctx app.Context, e app.Event) {
							c.SaveConfigFile()
						}).
						Text("Save"),
				),
			app.Div().
				Body(
					&components.Controlled{
						Component: app.Textarea().
							OnInput(func(ctx app.Context, e app.Event) {
								c.SetConfigFile(ctx.JSSrc.Get("value").String())
							}).
							Text(c.ConfigFile),
						Properties: map[string]interface{}{
							"value": c.ConfigFile,
						},
					},
				),
			app.If(
				c.Error != nil,
				app.Div().
					Body(
						app.H3().
							Text("Error"),
						app.Code().
							Text(c.Error),
						app.Button().
							OnClick(func(ctx app.Context, e app.Event) {
								c.Ignore()
							}).
							Text("Ignore"),
						app.Button().
							OnClick(func(ctx app.Context, e app.Event) {
								c.Recover()
							}).
							Text("Recover"),
					),
			),
		)
}
