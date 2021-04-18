package components

import "github.com/maxence-charriere/go-app/v8/pkg/app"

type DataShell struct {
	app.Compo

	AuthorizedWebDAVURL string
	ConfigFile          string

	ValidateConfigFile func(string)

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
			Value(
				c.AuthorizedWebDAVURL,
			),
		app.Textarea().Text(
			c.ConfigFile,
		),
		app.Button().
			OnClick(func(ctx app.Context, e app.Event) {
				// TODO: Use local value of textarea, this doesn't change yet
				c.ValidateConfigFile(c.ConfigFile)
			}).
			Text("Validate"),
	)
}
