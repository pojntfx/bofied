package components

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type Autofocused struct {
	app.Compo

	Component app.UI
}

func (c *Autofocused) Render() app.UI {
	return c.Component
}

func (c *Autofocused) OnUpdate(ctx app.Context) {
	ctx.Defer(func(_ app.Context) {
		c.JSValue().Call("focus")
	})
}
