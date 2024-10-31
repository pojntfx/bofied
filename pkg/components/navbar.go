package components

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
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
		Class("pf-v6-c-masthead pf-m-display-stack pf-m-display-inline-on-sm").
		Body(
			app.Div().
				Class("pf-v6-c-masthead__main").
				Body(
					app.Div().
						Class("pf-v6-c-masthead__brand").
						Body(
							app.A().
								Href("#").
								Class("pf-v6-c-page__header-brand-link").
								Body(
									app.Img().
										Class("pf-v6-c-brand pf-v6-x-c-brand--nav").
										Src("/web/logo.svg").
										Alt("Logo"),
								),
						),
				),
			app.Div().
				Class("pf-v6-c-masthead__content").
				Body(
					app.Div().
						Class("pf-v6-c-toolbar").
						Body(
							app.Div().
								Class("pf-v6-c-toolbar__content").
								Body(
									app.Div().
										Class("pf-v6-c-toolbar__content-section").
										Body(
											app.Div().
												Class("pf-v6-c-toolbar__group pf-m-align-end").
												Body(
													app.Div().
														Class(func() string {
															classes := "pf-v6-c-toolbar__item"

															if c.NotificationsDrawerOpen {
																classes += " pf-m-selected"
															}

															return classes
														}()).
														Body(
															app.Button().
																Class(
																	func() string {
																		classes := "pf-v6-c-button pf-m-plain"

																		if c.NotificationsDrawerOpen {
																			classes += " pf-m-read pf-m-stateful pf-m-clicked pf-m-expanded"
																		}

																		return classes
																	}()).
																Type("button").
																Aria("label", "Unread notifications").
																Aria("expanded", c.NotificationsDrawerOpen).
																OnClick(func(ctx app.Context, e app.Event) {
																	c.ToggleNotificationsDrawerOpen()
																}).
																Body(
																	app.Span().
																		Class("pf-v6-c-button__icon").
																		Body(
																			app.I().
																				Class("fas fa-bell").
																				Aria("hidden", true),
																		),
																),
														),
													app.Div().Class("pf-v6-c-toolbar__item").
														Body(
															app.Div().
																Class(func() string {
																	classes := "pf-v6-c-dropdown"

																	if c.OverflowMenuExpanded {
																		classes += " pf-m-expanded"
																	}

																	return classes
																}()).
																Body(
																	app.Button().
																		Class("pf-v6-c-menu-toggle pf-m-plain").
																		Type("button").
																		Aria("expanded", c.OverflowMenuExpanded).
																		Aria("label", "Actions").
																		Body(
																			app.Span().
																				Class("pf-v6-c-menu-toggle__text").
																				Body(
																					app.I().
																						Class("fas fa-ellipsis-v pf-v6-u-display-none-on-lg").
																						Aria("hidden", true),
																					app.I().
																						Class("fas fa-question-circle pf-v6-u-display-none pf-v6-u-display-inline-block-on-lg").
																						Aria("hidden", true),
																				),
																		).
																		OnClick(func(ctx app.Context, e app.Event) {
																			c.ToggleOverflowMenuExpanded()
																		}),

																	app.Div().
																		Class("pf-v6-c-menu pf-v6-x-u-position-absolute").
																		Hidden(!c.OverflowMenuExpanded).
																		Body(
																			app.Div().
																				Class("pf-v6-c-menu__content").
																				Body(
																					app.Ul().
																						Role("menu").
																						Class("pf-v6-c-menu__list").
																						Body(
																							app.Li().
																								Class("pf-v6-c-menu__list-item").
																								Role("none").
																								Body(
																									app.Button().
																										Class("pf-v6-c-menu__item").
																										Type("button").
																										Aria("role", "menuitem").
																										Body(
																											app.Span().
																												Class("pf-v6-c-menu__item-main").
																												Body(
																													app.Span().
																														Class("pf-v6-c-menu__item-text").
																														Text("About"),
																												),
																										).
																										OnClick(func(ctx app.Context, e app.Event) {
																											c.ToggleAbout()
																										}),
																								),
																							app.Li().
																								Class("pf-v6-c-menu__list-item").
																								Role("none").
																								Body(
																									app.A().
																										Class("pf-v6-c-menu__item").
																										Target("_blank").
																										Href("https://github.com/pojntfx/bofied#Usage").
																										Aria("role", "menuitem").
																										Body(
																											app.Span().
																												Class("pf-v6-c-menu__item-main").
																												Body(
																													app.Span().
																														Class("pf-v6-c-menu__item-text").
																														Text("Documentation"),
																												),
																										),
																								),
																							app.Li().
																								Class("pf-v6-c-divider pf-v6-u-display-inherit pf-v6-u-display-none-on-md").
																								Role("separator"),
																							app.Li().
																								Class("pf-v6-c-menu__list-item pf-v6-u-display-inherit pf-v6-u-display-none-on-md").
																								Role("none").
																								Body(
																									app.Button().
																										Class("pf-v6-c-menu__item").
																										Type("button").
																										Aria("role", "menuitem").
																										Body(
																											app.Span().
																												Class("pf-v6-c-menu__item-main").
																												Body(
																													app.Span().
																														Class("pf-v6-c-menu__item-icon").
																														Body(
																															app.I().
																																Class("fas fa-sign-out-alt").
																																Aria("hidden", true),
																														),
																													app.Span().
																														Class("pf-v6-c-menu__item-text").
																														Text("Logout"),
																												),
																										).
																										OnClick(func(ctx app.Context, e app.Event) {
																											c.Logout(ctx)
																										}),
																								),
																						),
																				),
																		),
																),
														),
													app.Div().
														Class("pf-v6-c-toolbar__item pf-m-hidden pf-m-visible-on-md").
														Body(
															app.Div().
																Class(func() string {
																	classes := "pf-v6-c-dropdown"

																	if c.UserMenuExpanded {
																		classes += " pf-m-expanded"
																	}

																	return classes
																}()).
																Body(
																	app.Button().
																		Class("pf-v6-c-dropdown__toggle pf-m-plain").
																		ID("page-layout-horizontal-nav-dropdown-kebab-2-button").
																		Aria("expanded", c.UserMenuExpanded).
																		Type("button").
																		Body(
																			app.Span().
																				Class("pf-v6-c-dropdown__toggle-text").
																				Text(c.UserEmail),
																			app.
																				Span().
																				Class("pf-v6-c-dropdown__toggle-icon").
																				Body(
																					app.I().
																						Class("fas fa-caret-down").
																						Aria("hidden", true),
																				),
																		).OnClick(func(ctx app.Context, e app.Event) {
																		c.ToggleUserMenuExpanded()
																	}),
																	app.Ul().
																		Class("pf-v6-c-dropdown__menu").
																		Aria("labelledby", "page-layout-horizontal-nav-dropdown-kebab-2-button").
																		Hidden(!c.UserMenuExpanded).
																		Body(
																			app.Li().Body(
																				app.Button().
																					Class("pf-v6-c-button pf-v6-c-dropdown__menu-item").
																					Type("button").
																					Body(
																						app.Span().
																							Class("pf-v6-c-button__icon pf-m-start").
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
														Class("pf-v6-c-toolbar__item").
														Body(
															app.Img().
																Class("pf-v6-c-avatar").
																Src(fmt.Sprintf("https://www.gravatar.com/avatar/%v?s=150", hex.EncodeToString(avatarHash[:]))).
																Alt("Avatar image"),
														),
												),
										),
								),
						),
				),
		)
}
