package components

import (
	"net/url"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/pojntfx/bofied/pkg/providers"
	"github.com/studio-b12/gowebdav"
)

type DataShell struct {
	app.Compo

	// Config file editor
	ConfigFile    string
	SetConfigFile func(string)

	FormatConfigFile  func()
	RefreshConfigFile func()
	SaveConfigFile    func()

	ConfigFileError       error
	IgnoreConfigFileError func()

	// File explorer
	CurrentPath    string
	SetCurrentPath func(string)

	Index        []os.FileInfo
	RefreshIndex func()
	WriteToPath  func(string, []byte)

	HTTPShareLink url.URL
	TFTPShareLink url.URL
	SharePath     func(string)

	CreatePath      func(string)
	CreateEmptyFile func(string)
	DeletePath      func(string)
	MovePath        func(string, string)
	CopyPath        func(string, string)

	EditPathContents    string
	SetEditPathContents func(string)
	EditPath            func(string)

	WebDAVAddress  url.URL
	WebDAVUsername string
	WebDAVPassword string

	OperationIndex []os.FileInfo

	OperationCurrentPath    string
	OperationSetCurrentPath func(string)

	FileExplorerError        error
	RecoverFileExplorerError func(app.Context)
	IgnoreFileExplorerError  func()

	Events []providers.Event

	EventsError        error
	RecoverEventsError func(app.Context)
	IgnoreEventsError  func()

	// Identity
	UserInfo oidc.UserInfo
	Logout   func(app.Context)

	// Metadata
	UseAdvertisedIP    bool
	SetUseAdvertisedIP func(bool)

	UseAdvertisedIPForWebDAV    bool
	SetUseAdvertisedIPForWebDAV func(bool)

	SetUseHTTPS func(bool)
	SetUseDavs  func(bool)

	// Internal state
	aboutDialogOpen         bool
	notificationsDrawerOpen bool
	overflowMenuExpanded    bool
	userMenuExpanded        bool
}

func (c *DataShell) Render() app.UI {
	// Gather notifications
	notifications := []Notification{}
	for _, event := range c.Events {
		notifications = append(notifications, Notification{
			CreatedAt: event.CreatedAt.String(),
			Message:   event.Message,
		})
	}

	// Reduce errors to global error
	globalError := c.FileExplorerError
	if c.EventsError != nil {
		globalError = c.EventsError
	}

	recoverGlobalError := c.RecoverFileExplorerError
	if c.EventsError != nil {
		recoverGlobalError = c.RecoverEventsError
	}

	ignoreGlobalError := c.IgnoreFileExplorerError
	if c.EventsError != nil {
		ignoreGlobalError = c.IgnoreEventsError
	}

	return app.Div().
		Class("pf-v6-u-h-100").
		Body(
			app.Div().
				Class("pf-v6-c-page").
				ID("page-layout-horizontal-nav").
				Body(
					app.A().
						Class("pf-v6-c-skip-to-content pf-v6-c-button pf-m-primary").
						Href("#main-content-page-layout-horizontal-nav").
						Text(
							"Skip to content",
						),
					&Navbar{
						NotificationsDrawerOpen: c.notificationsDrawerOpen,
						ToggleNotificationsDrawerOpen: func() {
							c.notificationsDrawerOpen = !c.notificationsDrawerOpen
							c.overflowMenuExpanded = false
						},

						ToggleAbout: func() {
							c.aboutDialogOpen = true
							c.overflowMenuExpanded = false
						},

						OverflowMenuExpanded: c.overflowMenuExpanded,
						ToggleOverflowMenuExpanded: func() {
							c.overflowMenuExpanded = !c.overflowMenuExpanded
							c.userMenuExpanded = false
						},

						UserMenuExpanded: c.userMenuExpanded,
						ToggleUserMenuExpanded: func() {
							c.userMenuExpanded = !c.userMenuExpanded
							c.overflowMenuExpanded = false
						},

						UserEmail: c.UserInfo.Email,
						Logout: func(ctx app.Context) {
							c.Logout(ctx)
						},
					},
					app.Div().
						Class("pf-v6-c-page__drawer").
						Body(
							app.Div().
								Class(func() string {
									classes := "pf-v6-c-drawer"

									if c.notificationsDrawerOpen {
										classes += " pf-m-expanded"
									}

									return classes
								}()).
								Body(
									app.Div().
										Class("pf-v6-c-drawer__main").
										Body(
											app.Div().
												Class("pf-v6-c-drawer__content").
												Body(
													app.Div().
														Class("pf-v6-c-page__main-container").
														TabIndex(-1).
														Body(
															app.Div().Class("pf-v6-c-drawer__body").Body(
																app.Main().
																	Class("pf-v6-c-page__main pf-v6-u-h-100").
																	ID("main-content-page-layout-horizontal-nav").
																	TabIndex(-1).
																	Body(
																		app.Section().
																			Class("pf-v6-c-page__main-section").
																			Body(
																				app.Div().
																					Class("pf-v6-l-grid pf-m-gutter pf-v6-u-h-100").
																					Body(
																						app.Div().
																							Class("pf-v6-l-grid__item pf-m-12-col pf-m-12-col-on-md pf-m-5-col-on-xl").
																							Body(
																								&TextEditorWrapper{
																									Title: "Config",

																									HelpLink: "https://github.com/pojntfx/bofied#config-script",

																									Error:            c.ConfigFileError,
																									ErrorDescription: "Syntax Error",
																									Ignore:           c.IgnoreConfigFileError,

																									Children: &TextEditor{
																										Content:    c.ConfigFile,
																										SetContent: c.SetConfigFile,

																										Format:  c.FormatConfigFile,
																										Refresh: c.RefreshConfigFile,
																										Save:    c.SaveConfigFile,

																										Language: "Go",
																									},
																								},
																							),

																						app.Div().
																							Class("pf-v6-l-grid__item pf-m-12-col pf-m-12-col-on-md pf-m-7-col-on-xl").
																							Body(
																								&FileExplorer{
																									CurrentPath:    c.CurrentPath,
																									SetCurrentPath: c.SetCurrentPath,

																									Index:        c.Index,
																									RefreshIndex: c.RefreshIndex,
																									WriteToPath:  c.WriteToPath,

																									HTTPShareLink: c.HTTPShareLink,
																									TFTPShareLink: c.TFTPShareLink,
																									SharePath:     c.SharePath,

																									CreatePath:      c.CreatePath,
																									CreateEmptyFile: c.CreateEmptyFile,
																									DeletePath:      c.DeletePath,
																									MovePath:        c.MovePath,
																									CopyPath:        c.CopyPath,

																									EditPathContents:    c.EditPathContents,
																									SetEditPathContents: c.SetEditPathContents,
																									EditPath:            c.EditPath,

																									WebDAVAddress:  c.WebDAVAddress,
																									WebDAVUsername: c.WebDAVUsername,
																									WebDAVPassword: c.WebDAVPassword,

																									OperationIndex: c.OperationIndex,

																									OperationCurrentPath:    c.OperationCurrentPath,
																									OperationSetCurrentPath: c.OperationSetCurrentPath,

																									UseAdvertisedIP:    c.UseAdvertisedIP,
																									SetUseAdvertisedIP: c.SetUseAdvertisedIP,

																									UseAdvertisedIPForWebDAV:    c.UseAdvertisedIPForWebDAV,
																									SetUseAdvertisedIPForWebDAV: c.SetUseAdvertisedIPForWebDAV,

																									SetUseHTTPS: c.SetUseHTTPS,
																									SetUseDavs:  c.SetUseDavs,

																									Nested: true,

																									GetContentType: func(fi os.FileInfo) string {
																										return fi.(gowebdav.File).ContentType()
																									},
																								},
																							),
																					),
																			),
																	),
															),
														),
												),
											app.Div().
												Class("pf-v6-c-drawer__panel").
												Body(
													app.Div().
														Class("pf-v6-c-drawer__body pf-m-no-padding").
														Body(
															&NotificationDrawer{
																Notifications: notifications,
																EmptyState: app.Div().
																	Class("pf-v6-c-empty-state").
																	Body(
																		app.Div().
																			Class("pf-v6-c-empty-state__content").
																			Body(
																				app.Div().
																					Class("pf-v6-c-empty-state__header").
																					Body(
																						app.I().
																							Class("fas fa-inbox pf-v6-c-empty-state__icon").
																							Aria("hidden", true),
																						app.Div().
																							Class("pf-v6-c-empty-state__title").Body(
																							app.H2().
																								Class("pf-v6-c-empty-state__title-text").
																								Text("No events yet"),
																						),
																					),
																				app.Div().
																					Class("pf-v6-c-empty-state__body").
																					Text("Network boot a node to see events here."),
																			),
																	),
															},
														),
												),
										),
								),
						),

					app.Ul().
						Class("pf-v6-c-alert-group pf-m-toast").
						Body(
							&UpdateNotification{
								UpdateTitle: "An update for bofied is available",

								StartUpdateText:  "Upgrade now",
								IgnoreUpdateText: "Maybe later",
							},
							app.If(
								globalError != nil,
								func() app.UI {
									return app.Li().
										Class("pf-v6-c-alert-group__item").
										Body(
											&Status{
												Error:       globalError,
												ErrorText:   "Fatal Error",
												Recover:     recoverGlobalError,
												RecoverText: "Reconnect",
												Ignore:      ignoreGlobalError,
											},
										)
								},
							),
						),

					&AboutModal{
						Open: c.aboutDialogOpen,
						Close: func() {
							c.aboutDialogOpen = false
						},

						ID: "about-modal-title",

						LogoDarkSrc: "/web/logo-dark.png",
						LogoDarkAlt: "bofied Logo (dark variant)",

						LogoLightSrc: "/web/logo-light.png",
						LogoLightAlt: "bofied Logo (light variant)",

						Title: "bofied",

						Body: app.Dl().
							Body(
								app.Dt().Text("Frontend version"),
								app.Dd().Text("main"),
								app.Dt().Text("Backend version"),
								app.Dd().Text("main"),
							),
						Footer: "Copyright © 2024 Felicitas Pojtinger and contributors (SPDX-License-Identifier: AGPL-3.0)",
					},
				),
		)
}
