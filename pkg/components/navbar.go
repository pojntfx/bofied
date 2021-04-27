package components

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type Navbar struct {
	app.Compo

	NotificationsDrawerOpen       bool
	ToggleNotificationsDrawerOpen func()

	ToggleAbout func()

	OverflowMenuExpanded       bool
	ToggleOverflowMenuExpanded func()

	UserMenuExpanded       bool
	ToggleUserMenuExpanded func()

	UserEmail string
	Logout    func(app.Context)
}

func (c *Navbar) Render() app.UI {
	// Get the MD5 hash for the user's gravatar
	avatarHash := md5.Sum([]byte(c.UserEmail))

	return app.Header().
		Class("pf-c-page__header").
		Body(
			app.Div().
				Class("pf-c-page__header-brand").
				Body(
					app.A().
						Href("#").
						Class("pf-c-page__header-brand-link").
						Body(
							app.Img().
								Class("pf-c-brand pf-x-c-brand--nav").
								Src("/web/logo.svg").
								Alt("Logo"),
						),
				),
			app.Div().
				Class("pf-c-page__header-tools").
				Body(
					app.Div().
						Class("pf-c-page__header-tools-group").
						Body(
							app.Div().
								Class("pf-c-page__header-tools-group").
								Body(
									app.Div().
										Class(func() string {
											classes := "pf-c-page__header-tools-item"

											if c.NotificationsDrawerOpen {
												classes += " pf-m-selected"
											}

											return classes
										}()).
										Body(
											app.Button().
												Class("pf-c-button pf-m-plain").
												Type("button").
												Aria("label", "Unread notifications").
												Aria("expanded", false).
												OnClick(func(ctx app.Context, e app.Event) {
													c.ToggleNotificationsDrawerOpen()
												}).
												Body(
													app.Span().
														Class("pf-c-notification-badge").
														Body(
															app.I().
																Class("pf-icon-bell").
																Aria("hidden", true),
														),
												),
										),
									app.Div().Class("pf-c-page__header-tools-item").
										Body(
											app.Div().
												Class(func() string {
													classes := "pf-c-dropdown"

													if c.OverflowMenuExpanded {
														classes += " pf-m-expanded"
													}

													return classes
												}()).
												Body(
													app.Button().
														Class("pf-c-dropdown__toggle pf-m-plain").
														ID("page-default-nav-example-dropdown-kebab-1-button").
														Aria("expanded", c.OverflowMenuExpanded).Type("button").
														Aria("label", "Actions").
														Body(
															app.I().
																Class("fas fa-ellipsis-v pf-u-display-none-on-lg").
																Aria("hidden", true),
															app.I().
																Class("fas fa-question-circle pf-u-display-none pf-u-display-inline-block-on-lg").
																Aria("hidden", true),
														).OnClick(func(ctx app.Context, e app.Event) {
														c.ToggleOverflowMenuExpanded()
													}),
													app.Ul().
														Class("pf-c-dropdown__menu pf-m-align-right").
														Aria("aria-labelledby", "page-default-nav-example-dropdown-kebab-1-button").
														Hidden(!c.OverflowMenuExpanded).
														Body(
															app.Li().
																Body(
																	app.A().
																		Class("pf-c-dropdown__menu-item").
																		Href("https://github.com/pojntfx/bofied#Usage").
																		Text("Documentation").
																		Target("_blank"),
																),
															app.Li().
																Body(
																	app.Button().
																		Class("pf-c-button pf-c-dropdown__menu-item").
																		Type("button").
																		OnClick(func(ctx app.Context, e app.Event) {
																			c.ToggleAbout()
																		}).
																		Text("About"),
																),
															app.Li().
																Class("pf-c-divider pf-u-display-none-on-md").
																Aria("role", "separator"),
															app.Li().
																Class("pf-u-display-none-on-md").
																Body(
																	app.Button().
																		Class("pf-c-button pf-c-dropdown__menu-item").
																		Type("button").
																		Body(
																			app.Span().
																				Class("pf-c-button__icon pf-m-start").
																				Body(
																					app.I().
																						Class("fas fa-sign-out-alt").
																						Aria("hidden", true),
																				),
																			app.Text("Logout"),
																		).
																		OnClick(func(ctx app.Context, e app.Event) {
																			c.Logout(ctx)
																		}),
																),
														),
												),
										),
									app.Div().
										Class("pf-c-page__header-tools-item pf-m-hidden pf-m-visible-on-md").
										Body(
											app.Div().
												Class(func() string {
													classes := "pf-c-dropdown"

													if c.UserMenuExpanded {
														classes += " pf-m-expanded"
													}

													return classes
												}()).
												Body(
													app.Button().
														Class("pf-c-dropdown__toggle pf-m-plain").
														ID("page-layout-horizontal-nav-dropdown-kebab-2-button").
														Aria("expanded", c.UserMenuExpanded).
														Type("button").
														Body(
															app.Span().
																Class("pf-c-dropdown__toggle-text").
																Text(c.UserEmail),
															app.
																Span().
																Class("pf-c-dropdown__toggle-icon").
																Body(
																	app.I().
																		Class("fas fa-caret-down").
																		Aria("hidden", true),
																),
														).OnClick(func(ctx app.Context, e app.Event) {
														c.ToggleUserMenuExpanded()
													}),
													app.Ul().
														Class("pf-c-dropdown__menu").
														Aria("labelledby", "page-layout-horizontal-nav-dropdown-kebab-2-button").
														Hidden(!c.UserMenuExpanded).
														Body(
															app.Li().Body(
																app.Button().
																	Class("pf-c-button pf-c-dropdown__menu-item").
																	Type("button").
																	Body(
																		app.Span().
																			Class("pf-c-button__icon pf-m-start").
																			Body(
																				app.I().
																					Class("fas fa-sign-out-alt").
																					Aria("hidden", true),
																			),
																		app.Text("Logout"),
																	).
																	OnClick(func(ctx app.Context, e app.Event) {
																		c.Logout(ctx)
																	}),
															),
														),
												),
										),
								),
							app.Img().Class("pf-c-avatar").Src(fmt.Sprintf("https://www.gravatar.com/avatar/%v?s=150", hex.EncodeToString(avatarHash[:]))).Alt("Avatar image"),
						),
				),
		)
}
