package components

import "github.com/maxence-charriere/go-app/v10/pkg/app"

type Notification struct {
	Message   string
	CreatedAt string
}

type NotificationDrawer struct {
	app.Compo

	Notifications []Notification
	EmptyState    app.UI
}

func (c *NotificationDrawer) Render() app.UI {
	return app.Div().
		Class("pf-v6-c-notification-drawer").
		Body(
			app.Div().
				Class("pf-v6-c-notification-drawer__header").
				Body(
					app.H1().
						Class("pf-v6-c-notification-drawer__header-title").
						Text("Events"),
				),
			app.Div().Class("pf-v6-c-notification-drawer__body").Body(
				app.If(
					len(c.Notifications) > 0,
					func() app.UI {
						return app.Ul().Class("pf-v6-c-notification-drawer__list").Body(
							app.Range(c.Notifications).Slice(func(i int) app.UI {
								return app.Li().Class("pf-v6-c-notification-drawer__list-item pf-m-read pf-m-info").Body(
									app.Div().Class("pf-v6-c-notification-drawer__list-item-description").Text(
										c.Notifications[len(c.Notifications)-1-i].Message,
									),
									app.Div().Class("pf-v6-c-notification-drawer__list-item-timestamp").Text(
										c.Notifications[len(c.Notifications)-1-i].CreatedAt,
									),
								)
							}),
						)
					},
				).Else(func() app.UI {
					return c.EmptyState
				}),
			),
		)
}
