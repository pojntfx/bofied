package components

import "github.com/maxence-charriere/go-app/v9/pkg/app"

type Controlled struct {
	app.Compo

	Component  app.UI
	Properties map[string]interface{}
}

func (c *Controlled) Render() app.UI {
	return c.Component
}

func (c *Controlled) OnUpdate(ctx app.Context) {
	ctx.Defer(func(_ app.Context) {
		for key, value := range c.Properties {
			c.JSValue().Set(key, value)
		}
	})
}
