package components

import "github.com/maxence-charriere/go-app/v8/pkg/app"

type DataShell struct {
	app.Compo

	GetAuthorizedWebDAVURL func() string

	Error   error
	Recover func()
	Ignore  func()
}

func (c *DataShell) Render() app.UI {
	return app.Input().Value(c.GetAuthorizedWebDAVURL())
}
