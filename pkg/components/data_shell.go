package components

import "github.com/maxence-charriere/go-app/v8/pkg/app"

type DataShell struct {
	app.Compo

	GetIDToken func() string
}

func (c *DataShell) Render() app.UI {
	return app.Text(c.GetIDToken())
}
