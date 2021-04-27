package components

import (
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/pojntfx/bofied/pkg/providers"
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

	HTTPShareLink string
	TFTPShareLink string
	SharePath     func(string)

	CreatePath func(string)
	DeletePath func(string)
	MovePath   func(string, string)
	CopyPath   func(string, string)

	WebDAVAddress  string
	WebDAVUsername string
	WebDAVPassword string

	OperationIndex []os.FileInfo

	OperationCurrentPath    string
	OperationSetCurrentPath func(string)

	FileExplorerError        error
	RecoverFileExplorerError func()
	IgnoreFileExplorerError  func()

	Events []providers.Event

	EventsError        error
	RecoverEventsError func(app.Context)
	IgnoreEventsError  func()

	// Identity
	UserInfo oidc.UserInfo
	Logout   func(app.Context)

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

	return app.Div().
		Class("pf-u-h-100").
		Body(
			app.Div().
				Class("pf-c-page").
				ID("page-layout-horizontal-nav").
				Aria("hidden", c.aboutDialogOpen).
				Body(
					app.A().
						Class("pf-c-skip-to-content pf-c-button pf-m-primary").
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
						Class("pf-c-page__drawer").
						Body(
							app.Div().
								Class(func() string {
									classes := "pf-c-drawer"

									if c.notificationsDrawerOpen {
										classes += " pf-m-expanded"
									}

									return classes
								}()).
								Body(
									app.Div().
										Class("pf-c-drawer__main").
										Body(
											app.Div().
												Class("pf-c-drawer__content").
												Body(
													app.Div().Class("pf-c-drawer__body").Body(
														app.Main().
															Class("pf-c-page__main pf-u-h-100").
															ID("main-content-page-layout-horizontal-nav").
															TabIndex(-1).
															Body(
																app.Section().
																	Class("pf-c-page__main-section").
																	Body(
																		&ConfigFileEditor{
																			ConfigFile:    c.ConfigFile,
																			SetConfigFile: c.SetConfigFile,

																			FormatConfigFile:  c.FormatConfigFile,
																			RefreshConfigFile: c.RefreshConfigFile,
																			SaveConfigFile:    c.SaveConfigFile,

																			Error:  c.ConfigFileError,
																			Ignore: c.IgnoreConfigFileError,
																		},

																		&FileExplorer{
																			CurrentPath:    c.CurrentPath,
																			SetCurrentPath: c.SetCurrentPath,

																			Index:        c.Index,
																			RefreshIndex: c.RefreshIndex,
																			WriteToPath:  c.WriteToPath,

																			HTTPShareLink: c.HTTPShareLink,
																			TFTPShareLink: c.TFTPShareLink,
																			SharePath:     c.SharePath,

																			CreatePath: c.CreatePath,
																			DeletePath: c.DeletePath,
																			MovePath:   c.MovePath,
																			CopyPath:   c.CopyPath,

																			WebDAVAddress:  c.WebDAVAddress,
																			WebDAVUsername: c.WebDAVUsername,
																			WebDAVPassword: c.WebDAVPassword,

																			OperationIndex: c.OperationIndex,

																			OperationCurrentPath:    c.OperationCurrentPath,
																			OperationSetCurrentPath: c.OperationSetCurrentPath,
																		},
																	),
															),
													),
												),
											app.Div().
												Class("pf-c-drawer__panel").
												Body(
													app.Div().
														Class("pf-c-drawer__body pf-m-no-padding").
														Body(
															&NotificationDrawer{
																Notifications: notifications,
															},
														),
												),
										),
								),
						),

					app.Ul().
						Class("pf-c-alert-group pf-m-toast").
						Body(
							app.If(
								c.FileExplorerError != nil,
								app.Li().
									Class("pf-c-alert-group__item").
									Body(
										&Status{
											Error:       c.FileExplorerError,
											ErrorText:   "Fatal Error",
											Recover:     c.RecoverFileExplorerError,
											RecoverText: "Reconnect",
											Ignore:      c.IgnoreFileExplorerError,
										},
									),
							),
						),

					&AboutModal{
						Open: c.aboutDialogOpen,
						Close: func() {
							c.aboutDialogOpen = false
						},

						ID: "about-modal-title",

						LogoSrc: "/web/logo.svg",
						LogoAlt: "Logo",
						Title:   "bofied",

						Body: app.Dl().
							Body(
								app.Dt().Text("Frontend version"),
								app.Dd().Text("main"),
								app.Dt().Text("Backend version"),
								app.Dd().Text("main"),
							),
						Footer: "Copyright © 2021 Felix Pojtinger and contributors (SPDX-License-Identifier: AGPL-3.0)",
					},
				),
		)
}
