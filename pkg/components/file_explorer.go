package components

import (
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type FileExplorer struct {
	app.Compo

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

	UseAdvertisedIP    bool
	SetUseAdvertisedIP func(bool)

	UseAdvertisedIPForWebDAV    bool
	SetUseAdvertisedIPForWebDAV func(bool)

	SetUseHTTPS func(bool)
	SetUseDavs  func(bool)

	Nested bool

	GetContentType func(os.FileInfo) string

	selectedPath     string
	newDirectoryName string
	newEmptyFilename string
	newFileName      string

	overflowMenuOpen bool

	mountFolderModalOpen     bool
	sharePathModalOpen       bool
	createDirectoryModalOpen bool
	createEmptyFileModalOpen bool
	deletionConfirmModalOpen bool
	renamePathModalOpen      bool
	movePathModalOpen        bool
	copyPathModalOpen        bool
	uploadModalOpen          bool
	editModalOpen            bool
	discardEditsModalOpen    bool

	// So that the "dirty" state in the text editor can be tracked
	cleanEditPathContents string

	operationSelectedPath string

	shareExpandableSectionOpen bool
	mountExpandableSectionOpen bool

	// True if accepting to discard edits leads to refresh the selected file
	discardEditsModalTargetsRefresh bool
}

func (c *FileExplorer) Render() app.UI {
	// Parse the current path
	pathComponents := []string{}
	for _, pathPart := range strings.Split(c.CurrentPath, string(os.PathSeparator)) {
		// Ignore empty paths
		if pathPart != "" {
			pathComponents = append(pathComponents, pathPart)
		}
	}

	// Parse the current operation path
	operationPathComponents := []string{}
	for _, pathPart := range strings.Split(c.OperationCurrentPath, string(os.PathSeparator)) {
		// Ignore empty paths
		if pathPart != "" {
			operationPathComponents = append(operationPathComponents, pathPart)
		}
	}

	// Check if we are using a secure protocols
	useHTTPS := false
	if c.HTTPShareLink.Scheme == "https" {
		useHTTPS = true
	}

	useDavs := false
	if c.WebDAVAddress.Scheme == "davs" {
		useDavs = true
	}

	// Check if the current path can be edited as text
	selectedPathContentType := ""
	for _, candidate := range c.Index {
		if !candidate.IsDir() && filepath.Join(c.CurrentPath, candidate.Name()) == c.selectedPath {
			ctype := c.GetContentType(candidate)

			if ctype == "application/json" || strings.HasPrefix(ctype, "text/") {
				selectedPathContentType = ctype
			}

			break
		}
	}

	return app.Div().
		Class("pf-v6-u-h-100").
		Body(
			app.Div().
				Class("pf-v6-c-card pf-m-plain pf-v6-u-h-100").
				Body(
					app.Div().Class("pf-v6-c-card__header").Body(
						app.Div().Class("pf-v6-c-card__actions").Body(
							app.Div().
								Class("pf-v6-c-toolbar pf-v6-u-py-0").
								Body(
									app.Div().
										Class("pf-v6-c-toolbar__content").
										Body(
											app.Div().
												Class("pf-v6-c-toolbar__content-section pf-v6-x-m-gap-md pf-v6-u-align-items-center").
												Body(
													app.Div().
														Class("pf-v6-c-toolbar__group pf-m-align-end pf-v6-u-align-items-center").
														Body(
															app.If(
																c.selectedPath != "",
																func() app.UI {
																	return app.Div().
																		Class("pf-v6-c-toolbar__item pf-v6-u-display-none pf-v6-u-display-flex-on-lg").
																		Body(
																			app.Button().
																				Type("button").
																				Aria("label", "Share file").
																				Title("Share file").
																				Class("pf-v6-c-button pf-m-plain").
																				OnClick(func(ctx app.Context, e app.Event) {
																					c.sharePath()
																				}).
																				Body(
																					app.Span().
																						Class("pf-v6-c-button__icon").
																						Body(
																							app.I().
																								Class("fas fa-share-alt").
																								Aria("hidden", true),
																						),
																				),
																		)
																},
															),

															app.If(
																c.selectedPath != "",
																func() app.UI {
																	return app.Div().
																		Class("pf-v6-c-toolbar__item pf-v6-u-display-none pf-v6-u-display-flex-on-lg").
																		Body(
																			app.Button().
																				Type("button").
																				Aria("label", "Delete file").
																				Title("Delete file").
																				Class("pf-v6-c-button pf-m-plain").
																				OnClick(func(ctx app.Context, e app.Event) {
																					c.deleteFile()
																				}).
																				Body(
																					app.Span().
																						Class("pf-v6-c-button__icon").
																						Body(
																							app.I().
																								Class("fas fa-trash").
																								Aria("hidden", true),
																						),
																				),
																		)
																},
															),

															app.Div().Class("pf-v6-c-toolbar__item").
																Body(
																	app.Div().
																		Class(func() string {
																			classes := "pf-v6-c-dropdown"

																			if c.overflowMenuOpen {
																				classes += " pf-m-expanded"
																			}

																			return classes
																		}()).
																		Body(
																			app.If(
																				c.selectedPath != "",
																				func() app.UI {
																					return app.Button().
																						Class("pf-v6-c-menu-toggle pf-m-plain pf-v6-u-display-none pf-v6-u-display-flex-on-lg").
																						Type("button").
																						Aria("expanded", c.overflowMenuOpen).
																						Aria("label", "Actions").
																						Body(
																							app.Span().
																								Class("pf-v6-c-menu-toggle__text pf-v6-u-display-flex pf-v6-u-display-block-on-md").
																								Body(
																									app.I().
																										Class("fas fa-ellipsis-v").
																										Aria("hidden", true),
																								),
																						).
																						OnClick(func(ctx app.Context, e app.Event) {
																							c.overflowMenuOpen = !c.overflowMenuOpen
																						})
																				},
																			),

																			app.Button().
																				Class("pf-v6-c-menu-toggle pf-m-plain pf-v6-u-display-inherit pf-v6-u-display-none-on-lg").
																				Type("button").
																				Aria("expanded", c.overflowMenuOpen).
																				Aria("label", "Actions").
																				Body(
																					app.Span().
																						Class("pf-v6-c-menu-toggle__text pf-v6-u-display-flex pf-v6-u-display-block-on-md").
																						Body(
																							app.I().
																								Class("fas fa-ellipsis-v").
																								Aria("hidden", true),
																						),
																				).
																				OnClick(func(ctx app.Context, e app.Event) {
																					c.overflowMenuOpen = !c.overflowMenuOpen
																				}),

																			app.Div().
																				Class("pf-v6-c-menu pf-v6-x-u-position-absolute pf-v6-x-dropdown-menu").
																				Hidden(!c.overflowMenuOpen).
																				Body(
																					app.Div().
																						Class("pf-v6-c-menu__content").
																						Body(
																							app.Ul().
																								Role("menu").
																								Class("pf-v6-c-menu__list").
																								Body(
																									app.If(
																										c.selectedPath != "" && selectedPathContentType != "",
																										func() app.UI {
																											return app.Li().
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
																																		Class("pf-v6-c-menu__item-icon").
																																		Body(
																																			app.I().
																																				Class("fas fa-pencil-alt").
																																				Aria("hidden", true),
																																		),
																																	app.Span().
																																		Class("pf-v6-c-menu__item-text").
																																		Text("Edit file"),
																																),
																														).
																														OnClick(func(ctx app.Context, e app.Event) {
																															c.editPath()
																														}),
																												)
																										},
																									),
																									app.If(
																										c.selectedPath != "",
																										func() app.UI {
																											return app.Li().
																												Class("pf-v6-c-menu__list-item pf-v6-u-display-inherit pf-v6-u-display-none-on-lg").
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
																																				Class("fas fa-share-alt").
																																				Aria("hidden", true),
																																		),
																																	app.Span().
																																		Class("pf-v6-c-menu__item-text").
																																		Text("Share file"),
																																),
																														).
																														OnClick(func(ctx app.Context, e app.Event) {
																															c.sharePath()
																														}),
																												)
																										},
																									),
																									app.If(
																										c.selectedPath != "",
																										func() app.UI {
																											return app.Li().
																												Class("pf-v6-c-menu__list-item pf-v6-u-display-inherit pf-v6-u-display-none-on-lg").
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
																																				Class("fas fa-trash").
																																				Aria("hidden", true),
																																		),
																																	app.Span().
																																		Class("pf-v6-c-menu__item-text").
																																		Text("Delete file"),
																																),
																														).
																														OnClick(func(ctx app.Context, e app.Event) {
																															c.deleteFile()
																														}),
																												)
																										},
																									),

																									app.If(
																										c.selectedPath != "",
																										func() app.UI {
																											return app.Li().
																												Class(func() string {
																													base := "pf-v6-c-divider"

																													// Always show the divider if we have an editable file selected, else make it responsive
																													if selectedPathContentType == "" {
																														base += " pf-v6-u-display-inherit pf-v6-u-display-none-on-lg"
																													}

																													return base
																												}()).
																												Role("separator")
																										},
																									),

																									app.If(
																										c.selectedPath != "",
																										func() app.UI {
																											return app.Li().
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
																																		Text("Move to ..."),
																																),
																														).
																														OnClick(func(ctx app.Context, e app.Event) {
																															c.moveTo()
																														}),
																												)
																										},
																									),
																									app.If(
																										c.selectedPath != "",
																										func() app.UI {
																											return app.Li().
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
																																		Text("Copy to ..."),
																																),
																														).
																														OnClick(func(ctx app.Context, e app.Event) {
																															c.copyTo()
																														}),
																												)
																										},
																									),

																									app.If(
																										c.selectedPath != "",
																										func() app.UI {
																											return app.Li().
																												Class("pf-v6-c-divider").
																												Role("separator")
																										},
																									),

																									app.If(
																										c.selectedPath != "",
																										func() app.UI {
																											return app.Li().
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
																																		Text("Rename"),
																																),
																														).
																														OnClick(func(ctx app.Context, e app.Event) {
																															c.rename()
																														}),
																												)
																										},
																									),

																									app.If(
																										c.selectedPath != "",
																										func() app.UI {
																											return app.Li().
																												Class("pf-v6-c-divider pf-v6-u-display-inherit pf-v6-u-display-none-on-lg").
																												Role("separator")
																										},
																									),

																									app.Li().
																										Class("pf-v6-c-menu__list-item pf-v6-u-display-inherit pf-v6-u-display-none-on-lg").
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
																																		Class("fas fa-folder-plus").
																																		Aria("hidden", true),
																																),
																															app.Span().
																																Class("pf-v6-c-menu__item-text").
																																Text("Create directory"),
																														),
																												).
																												OnClick(func(ctx app.Context, e app.Event) {
																													c.createDirectory()
																												}),
																										),
																									app.Li().
																										Class("pf-v6-c-menu__list-item pf-v6-u-display-inherit pf-v6-u-display-none-on-lg").
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
																																		Class("fas fa-pen-square").
																																		Aria("hidden", true),
																																),
																															app.Span().
																																Class("pf-v6-c-menu__item-text").
																																Text("Create empty file"),
																														),
																												).
																												OnClick(func(ctx app.Context, e app.Event) {
																													c.createEmptyFile()
																												}),
																										),
																									app.Li().
																										Class("pf-v6-c-menu__list-item pf-v6-u-display-inherit pf-v6-u-display-none-on-lg").
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
																																		Class("fas fa-cloud-upload-alt").
																																		Aria("hidden", true),
																																),
																															app.Span().
																																Class("pf-v6-c-menu__item-text").
																																Text("Upload file"),
																														),
																												).
																												OnClick(func(ctx app.Context, e app.Event) {
																													c.uploadFile()
																												}),
																										),

																									app.Li().
																										Class("pf-v6-c-divider pf-v6-u-display-inherit pf-v6-u-display-none-on-lg").
																										Role("separator"),

																									app.Li().
																										Class("pf-v6-c-menu__list-item pf-v6-u-display-inherit pf-v6-u-display-none-on-lg").
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
																																		Class("fas fas fa-sync").
																																		Aria("hidden", true),
																																),
																															app.Span().
																																Class("pf-v6-c-menu__item-text").
																																Text("Refresh"),
																														),
																												).
																												OnClick(func(ctx app.Context, e app.Event) {
																													c.refresh()
																												}),
																										),
																								),
																						),
																				),
																		),
																),

															app.If(
																c.selectedPath != "",
																func() app.UI {
																	return app.Hr().Class("pf-v6-c-divider pf-m-vertical pf-v6-u-display-none pf-v6-u-display-flex-on-lg")
																},
															),

															app.Div().
																Class("pf-v6-c-toolbar__item pf-v6-u-display-none pf-v6-u-display-flex-on-lg").
																Body(
																	app.Button().
																		Type("button").
																		Aria("label", "Create directory").
																		Title("Create directory").
																		Class("pf-v6-c-button pf-m-plain").
																		OnClick(func(ctx app.Context, e app.Event) {
																			c.createDirectory()
																		}).
																		Body(
																			app.Span().
																				Class("pf-v6-c-button__icon").
																				Body(
																					app.I().
																						Class("fas fa-folder-plus").
																						Aria("hidden", true),
																				),
																		),
																),
															app.Div().
																Class("pf-v6-c-toolbar__item pf-v6-u-display-none pf-v6-u-display-flex-on-lg").
																Body(
																	app.Button().
																		Type("button").
																		Aria("label", "Create empty file").
																		Title("Create empty file").
																		Class("pf-v6-c-button pf-m-plain").
																		OnClick(func(ctx app.Context, e app.Event) {
																			c.createEmptyFile()
																		}).
																		Body(
																			app.Span().
																				Class("pf-v6-c-button__icon").
																				Body(
																					app.I().
																						Class("fas fa-pen-square").
																						Aria("hidden", true),
																				),
																		),
																),
															app.Div().
																Class("pf-v6-c-toolbar__item pf-v6-u-display-none pf-v6-u-display-flex-on-lg").
																Body(
																	app.Button().
																		Type("button").
																		Aria("label", "Upload file").
																		Title("Upload file").
																		Class("pf-v6-c-button pf-m-plain").
																		OnClick(func(ctx app.Context, e app.Event) {
																			c.uploadFile()
																		}).
																		Body(
																			app.Span().
																				Class("pf-v6-c-button__icon").
																				Body(
																					app.I().
																						Class("fas fa-cloud-upload-alt").
																						Aria("hidden", true),
																				),
																		),
																),

															app.Hr().Class("pf-v6-c-divider pf-m-vertical pf-v6-u-display-none pf-v6-u-display-flex-on-lg"),

															app.Div().
																Class("pf-v6-c-toolbar__item pf-v6-u-display-none pf-v6-u-display-flex-on-lg").
																Body(
																	app.Button().
																		Type("button").
																		Aria("label", "Refresh").
																		Title("Refresh").
																		Class("pf-v6-c-button pf-m-plain").
																		OnClick(func(ctx app.Context, e app.Event) {
																			c.refresh()
																		}).
																		Body(
																			app.Span().
																				Class("pf-v6-c-button__icon").
																				Body(
																					app.I().
																						Class("fas fas fa-sync").
																						Aria("hidden", true),
																				),
																		),
																),
															app.Div().
																Class("pf-v6-c-toolbar__item").
																Body(
																	app.Button().
																		Class("pf-v6-c-button pf-m-control").
																		Type("button").
																		OnClick(func(ctx app.Context, e app.Event) {
																			c.mountDirectory()
																		}).
																		Body(
																			app.Span().
																				Class("pf-v6-c-button__icon pf-m-start").
																				Body(
																					app.I().
																						Class("fas fa-hdd").
																						Aria("hidden", true),
																				),
																			app.Text("Mount directory"),
																		),
																),
														),
												),
										),
								),
						),
						app.Div().Class("pf-v6-c-card__header-main").Body(
							app.Div().Class("pf-v6-c-card__title").Body(
								&Breadcrumbs{
									PathComponents: pathComponents,

									CurrentPath:    c.CurrentPath,
									SetCurrentPath: c.SetCurrentPath,

									SelectedPath: c.selectedPath,
									SetSelectedPath: func(s string) {
										c.selectedPath = s
									},

									ItemClass: "pf-v6-c-card__title-text",
								},
							),
						),
					),
					app.Div().Class("pf-v6-c-card__body").Body(
						app.If(
							len(c.Index) > 0,
							func() app.UI {
								return &FileGrid{
									Index: c.Index,

									SelectedPath: c.selectedPath,
									SetSelectedPath: func(s string) {
										c.selectedPath = s
									},

									CurrentPath:    c.CurrentPath,
									SetCurrentPath: c.SetCurrentPath,
								}
							},
						).Else(
							func() app.UI {
								return &EmptyState{
									Action: app.Button().
										Class("pf-v6-c-button pf-m-primary").
										Type("button").
										OnClick(func(ctx app.Context, e app.Event) {
											c.uploadModalOpen = true
										}).
										Body(
											app.Span().
												Class("pf-v6-c-button__icon pf-m-start").
												Body(
													app.I().
														Class("fas fa-cloud-upload-alt").
														Aria("hidden", true),
												),
											app.Text("Upload File"),
										),
								}
							},
						),
					),
				),

			&Modal{
				Open: c.mountFolderModalOpen,
				Close: func() {
					c.mountFolderModalOpen = false
					c.mountExpandableSectionOpen = false
				},

				ID:     "mount-folder-modal-title",
				Nested: c.Nested,

				Title: "Mount Folder",
				Body: []app.UI{
					app.Div().
						Class("pf-v6-c-content").
						Body(
							app.P().
								Text(`You can mount this folder as a WebDAV share using your system's file explorer. To do so, please use the following credentials:`),
							app.Form().
								Class("pf-v6-c-form").
								OnSubmit(func(ctx app.Context, e app.Event) {
									e.PreventDefault()
								}).
								Body(
									app.Div().
										Class("pf-v6-c-form__group").
										Body(
											app.Div().
												Class("pf-v6-c-form__group-label").
												Body(
													app.Label().
														Class("pf-v6-c-form__label").
														For("webdav-address").
														Body(
															app.Span().
																Class("pf-v6-c-form__label-text").
																Text("Address"),
														),
												),
											app.Div().
												Class("pf-v6-c-form__group-control").
												Body(
													&CopyableInput{
														Component: app.Input().
															Class("pf-v6-c-form-control").
															ReadOnly(true).
															Type("text").
															Value(c.WebDAVAddress.String()).
															Aria("label", "WebDAV server address").
															Name("webdav-address").
															ID("webdav-address"),
														ID: "webdav-address",
													},
												),
										),
									app.Div().
										Class("pf-v6-c-form__group").
										Body(
											app.Div().
												Class("pf-v6-c-form__group-label").
												Body(
													app.Label().
														Class("pf-v6-c-form__label").
														For("webdav-username").
														Body(
															app.Span().
																Class("pf-v6-c-form__label-text").
																Text("Username"),
														),
												),
											app.Div().
												Class("pf-v6-c-form__group-control").
												Body(
													&CopyableInput{
														Component: app.Input().
															Class("pf-v6-c-form-control").
															ReadOnly(true).
															Type("text").
															Value(c.WebDAVUsername).
															Aria("label", "WebDAV username").
															Name("webdav-username").
															ID("webdav-username"),
														ID: "webdav-username",
													},
												),
										),
									app.Div().
										Class("pf-v6-c-form__group").
										Body(
											app.Div().
												Class("pf-v6-c-form__group-label").
												Body(
													app.Label().
														Class("pf-v6-c-form__label").
														For("webdav-password").
														Body(
															app.Span().
																Class("pf-v6-c-form__label-text").
																Text("Password (One-Time Token)"),
														),
												),
											app.Div().
												Class("pf-v6-c-form__group-control").
												Body(
													&CopyableInput{
														Component: app.Input().
															Class("pf-v6-c-form-control").
															ReadOnly(true).
															Type("text").
															Value(c.WebDAVPassword).
															Aria("label", "WebDAV password").
															Name("webdav-password").
															ID("webdav-password"),
														ID: "webdav-password",
													},
												),
										),
									&ExpandableSection{
										Open: c.mountExpandableSectionOpen,
										OnToggle: func() {
											c.mountExpandableSectionOpen = !c.mountExpandableSectionOpen
										},
										Title:       "Advanced",
										ClosedTitle: "Show advanced options",
										OpenTitle:   "Hide advanced options",
										Body: []app.UI{
											app.Form().
												NoValidate(true).
												Class("pf-v6-c-form pf-m-horizontal pf-v6-u-mb-md").
												Body(
													&FormGroup{
														NoTopPadding:     true,
														NoControlWrapper: true,
														Label: app.Label().
															For("use-advertised-ip-for-webdav").
															Class("pf-v6-c-form__label").
															Body(
																app.
																	Span().
																	Class("pf-v6-c-form__label-text").
																	Text("Use advertised IP"),
															),
														Input: &Switch{
															ID: "use-advertised-ip-for-webdav",

															Open: c.UseAdvertisedIPForWebDAV,
															ToggleOpen: func() {
																c.SetUseAdvertisedIPForWebDAV(!c.UseAdvertisedIPForWebDAV)
															},

															OnMessage:  "Using advertised IP",
															OffMessage: "Not using advertised IP",
														},
													},
													&FormGroup{
														NoTopPadding:     true,
														NoControlWrapper: true,
														Label: app.Label().
															For("use-davs").
															Class("pf-v6-c-form__label").
															Body(
																app.
																	Span().
																	Class("pf-v6-c-form__label-text").
																	Text("Use secure protocol"),
															),
														Input: &Switch{
															ID: "use-davs",

															Open: useDavs,
															ToggleOpen: func() {
																c.SetUseDavs(!useDavs)
															},

															OnMessage:  "Using secure protocol",
															OffMessage: "Not using secure protocol",
														},
													},
												),
										},
									},
								),
						),
				},
				Footer: []app.UI{
					app.Button().
						Class("pf-v6-c-button pf-m-primary").
						OnClick(func(ctx app.Context, e app.Event) {
							c.mountFolderModalOpen = false
						}).
						Text("OK"),
				},
			},

			&Modal{
				Open: c.sharePathModalOpen,
				Close: func() {
					c.sharePathModalOpen = false
					c.shareExpandableSectionOpen = false
				},

				ID:     "share-path-modal-title",
				Nested: c.Nested,

				Title: `Share "` + path.Base(c.selectedPath) + `"`,
				Body: []app.UI{
					app.Div().
						Class("pf-v6-c-content").
						Body(
							app.P().
								Text(`You should be able to access this path using the following addresses. Please note that, especially for the TFTP address, access is probably only possible from the LAN on which bofied is running:`),
							app.Form().
								Class("pf-v6-c-form").
								OnSubmit(func(ctx app.Context, e app.Event) {
									e.PreventDefault()
								}).
								Body(
									app.Div().
										Class("pf-v6-c-form__group").
										Body(
											app.Div().
												Class("pf-v6-c-form__group-label").
												Body(
													app.Label().
														Class("pf-v6-c-form__label").
														For("http-address").
														Body(
															app.Span().
																Class("pf-v6-c-form__label-text").
																Text("HTTP Address"),
														),
												),
											app.Div().
												Class("pf-v6-c-form__group-control").
												Body(
													&CopyableInput{
														Component: app.Input().
															Class("pf-v6-c-form-control").
															ReadOnly(true).
															Type("text").
															Value(c.HTTPShareLink.String()).
															Aria("label", "HTTP address").
															Name("http-address").
															ID("http-address"),
														ID: "http-address",
													},
												),
										),
									app.Div().
										Class("pf-v6-c-form__group").
										Body(
											app.Div().
												Class("pf-v6-c-form__group-label").
												Body(
													app.Label().
														Class("pf-v6-c-form__label").
														For("tftp-address").
														Body(
															app.Span().
																Class("pf-v6-c-form__label-text").
																Text("TFTP Address"),
														),
												),
											app.Div().
												Class("pf-v6-c-form__group-control").
												Body(
													&CopyableInput{
														Component: app.Input().
															Class("pf-v6-c-form-control").
															ReadOnly(true).
															Type("text").
															Value(c.TFTPShareLink.String()).
															Aria("label", "TFTP address").
															Name("tftp-address").
															ID("tftp-address"),
														ID: "tftp-address",
													},
												),
										),
									&ExpandableSection{
										Open: c.shareExpandableSectionOpen,
										OnToggle: func() {
											c.shareExpandableSectionOpen = !c.shareExpandableSectionOpen
										},
										Title:       "Advanced",
										ClosedTitle: "Show advanced options",
										OpenTitle:   "Hide advanced options",
										Body: []app.UI{
											app.Form().
												NoValidate(true).
												Class("pf-v6-c-form pf-m-horizontal pf-v6-u-mb-md").
												Body(
													&FormGroup{
														NoTopPadding:     true,
														NoControlWrapper: true,
														Label: app.Label().
															For("use-advertised-ip").
															Class("pf-v6-c-form__label").
															Body(
																app.
																	Span().
																	Class("pf-v6-c-form__label-text").
																	Text("Use advertised IP"),
															),
														Input: &Switch{
															ID: "use-advertised-ip",

															Open: c.UseAdvertisedIP,
															ToggleOpen: func() {
																c.SetUseAdvertisedIP(!c.UseAdvertisedIP)
															},

															OnMessage:  "Using advertised IP",
															OffMessage: "Not using advertised IP",
														},
													},
													&FormGroup{
														NoTopPadding:     true,
														NoControlWrapper: true,
														Label: app.Label().
															For("use-https").
															Class("pf-v6-c-form__label").
															Body(
																app.
																	Span().
																	Class("pf-v6-c-form__label-text").
																	Text("Use secure protocol"),
															),
														Input: &Switch{
															ID: "use-https",

															Open: useHTTPS,
															ToggleOpen: func() {
																c.SetUseHTTPS(!useHTTPS)
															},

															OnMessage:  "Using secure protocol",
															OffMessage: "Not using secure protocol",
														},
													},
												),
										},
									},
								),
						),
				},
				Footer: []app.UI{
					app.Button().
						Class("pf-v6-c-button pf-m-primary").
						OnClick(func(ctx app.Context, e app.Event) {
							c.sharePathModalOpen = false
						}).
						Text("OK"),
				},
			},

			&Modal{
				Open: c.createDirectoryModalOpen,
				Close: func() {
					c.createDirectoryModalOpen = false
					c.newDirectoryName = ""
				},

				ID:      "create-directory-modal-title",
				Overlay: true, // It should be possible to create a directory from another modal
				Nested:  true,

				Title: "Create Directory",
				Body: []app.UI{
					app.Form().
						Class("pf-v6-c-form").
						ID("create-directory").
						OnSubmit(func(ctx app.Context, e app.Event) {
							e.PreventDefault()

							// Switch the base path if in move or copy operations
							basePath := c.CurrentPath
							if c.movePathModalOpen || c.copyPathModalOpen {
								basePath = c.OperationCurrentPath
							}

							c.CreatePath(filepath.Join(basePath, c.newDirectoryName))

							c.newDirectoryName = ""
							c.createDirectoryModalOpen = false
						}).
						Body(
							&FormGroup{
								Label: app.Label().
									For("directory-name-input").
									Class("pf-v6-c-form__label").
									Body(
										app.
											Span().
											Class("pf-v6-c-form__label-text").
											Text("Directory name"),
									),
								Input: &Autofocused{
									Component: app.Input().
										Name("directory-name-input").
										ID("directory-name-input").
										Type("text").
										Required(true).
										Class("pf-v6-c-form-control").
										Value(c.newDirectoryName).
										OnInput(func(ctx app.Context, e app.Event) {
											c.newDirectoryName = ctx.JSSrc().Get("value").String()
										}),
								},
								Required: true,
							},
						),
				},
				Footer: []app.UI{
					app.Button().
						Class("pf-v6-c-button pf-m-primary").
						Type("submit").
						Form("create-directory").
						Text("Create"),
					app.Button().
						Class("pf-v6-c-button pf-m-link").
						Type("button").
						OnClick(func(ctx app.Context, e app.Event) {
							c.newDirectoryName = ""
							c.createDirectoryModalOpen = false
						}).
						Text("Cancel"),
				},
			},

			&Modal{
				Open: c.createEmptyFileModalOpen,
				Close: func() {
					c.createEmptyFileModalOpen = false
					c.newEmptyFilename = ""
				},

				ID:     "create-empty-file-modal-title",
				Nested: true,

				Title: "Create Empty File",
				Body: []app.UI{
					app.Form().
						Class("pf-v6-c-form").
						ID("create-empty-file").
						OnSubmit(func(ctx app.Context, e app.Event) {
							e.PreventDefault()

							emptyFilePath := filepath.Join(c.CurrentPath, c.newEmptyFilename)
							c.CreateEmptyFile(emptyFilePath)

							c.newEmptyFilename = ""
							c.selectedPath = emptyFilePath
							c.createEmptyFileModalOpen = false
						}).
						Body(
							&FormGroup{
								Label: app.Label().
									For("new-filename-input").
									Class("pf-v6-c-form__label").
									Body(
										app.
											Span().
											Class("pf-v6-c-form__label-text").
											Text("Filename"),
									),
								Input: &Autofocused{
									Component: app.Input().
										Name("new-filename-input").
										ID("new-filename-input").
										Type("text").
										Required(true).
										Class("pf-v6-c-form-control").
										Value(c.newEmptyFilename).
										OnInput(func(ctx app.Context, e app.Event) {
											c.newEmptyFilename = ctx.JSSrc().Get("value").String()
										}),
								},
								Required: true,
							},
						),
				},
				Footer: []app.UI{
					app.Button().
						Class("pf-v6-c-button pf-m-primary").
						Type("submit").
						Form("create-empty-file").
						Text("Create"),
					app.Button().
						Class("pf-v6-c-button pf-m-link").
						Type("button").
						OnClick(func(ctx app.Context, e app.Event) {
							c.newEmptyFilename = ""
							c.createEmptyFileModalOpen = false
						}).
						Text("Cancel"),
				},
			},

			&Modal{
				Open: c.deletionConfirmModalOpen,
				Close: func() {
					c.deletionConfirmModalOpen = false
				},

				ID:     "deletion-confirm-modal-title",
				Nested: c.Nested,

				Title: `Permanently delete "` + path.Base(c.selectedPath) + `"?`,
				Body: []app.UI{
					app.P().Text(`If you delete an item, it will be permanently lost.`),
				},
				Footer: []app.UI{
					app.Button().
						Class("pf-v6-c-button pf-m-danger").
						Type("button").
						OnClick(func(ctx app.Context, e app.Event) {
							c.DeletePath(c.selectedPath)

							c.selectedPath = ""
							c.deletionConfirmModalOpen = false
						}).
						Text("Delete"),
					app.Button().
						Class("pf-v6-c-button pf-m-link").
						Type("button").
						OnClick(func(ctx app.Context, e app.Event) {
							c.deletionConfirmModalOpen = false
						}).
						Text("Cancel"),
				},
			},

			&Modal{
				Open: c.renamePathModalOpen,
				Close: func() {
					c.renamePathModalOpen = false
				},

				ID:     "rename-path-modal-title",
				Nested: c.Nested,

				Title: `Rename "` + path.Base(c.selectedPath) + `"`,
				Body: []app.UI{
					app.Form().
						Class("pf-v6-c-form").
						ID("rename-path").
						OnSubmit(func(ctx app.Context, e app.Event) {
							e.PreventDefault()

							c.MovePath(c.selectedPath, filepath.Join(c.CurrentPath, c.newFileName))

							c.newFileName = ""
							c.renamePathModalOpen = false
						}).
						Body(
							&FormGroup{
								Label: app.Label().
									For("path-rename-input").
									Class("pf-v6-c-form__label").
									Body(
										app.
											Span().
											Class("pf-v6-c-form__label-text").
											Text("New name"),
									),
								Input: &Autofocused{
									Component: app.Input().
										Name("path-rename-input").
										ID("path-rename-input").
										Type("text").
										Required(true).
										Class("pf-v6-c-form-control").
										Value(c.newFileName).
										OnInput(func(ctx app.Context, e app.Event) {
											c.newFileName = ctx.JSSrc().Get("value").String()
										}),
								},
								Required: true,
							},
						),
				},
				Footer: []app.UI{
					app.Button().
						Class("pf-v6-c-button pf-m-primary").
						Type("submit").
						Form("rename-path").
						Text("Rename"),
					app.Button().
						Class("pf-v6-c-button pf-m-link").
						Type("button").
						OnClick(func(ctx app.Context, e app.Event) {
							c.newFileName = ""
							c.renamePathModalOpen = false
						}).
						Text("Cancel"),
				},
			},

			&Modal{
				Open: c.movePathModalOpen,
				Close: func() {
					c.movePathModalOpen = false
				},

				ID:           "move-path-modal-title",
				Large:        true,
				PaddedBottom: true,
				Nested:       c.Nested,

				Title: `Move "` + path.Base(c.selectedPath) + `"`,
				Body: []app.UI{
					&PathPickerToolbar{
						Index:        c.OperationIndex,
						RefreshIndex: c.RefreshIndex,

						PathComponents: operationPathComponents,

						CurrentPath:    c.OperationCurrentPath,
						SetCurrentPath: c.OperationSetCurrentPath,

						SelectedPath: c.operationSelectedPath,
						SetSelectedPath: func(s string) {
							c.operationSelectedPath = s
						},

						OpenCreateDirectoryModal: func() {
							c.createDirectoryModalOpen = true
						},
					},
					app.If(
						len(c.OperationIndex) > 0,
						func() app.UI {
							return &FileGrid{
								Index: c.OperationIndex,

								SelectedPath: c.operationSelectedPath,
								SetSelectedPath: func(s string) {
									c.operationSelectedPath = s
								},

								CurrentPath:    c.OperationCurrentPath,
								SetCurrentPath: c.OperationSetCurrentPath,

								Standalone: true,
							}
						},
					).Else(
						func() app.UI {
							return &EmptyState{}
						},
					),
				},
				Footer: []app.UI{
					func() app.UI {
						// Prefer selected path, fall back to current path if not selected
						newDirectory := c.operationSelectedPath
						if newDirectory == "" {
							newDirectory = c.OperationCurrentPath
						}

						return app.Button().
							Class("pf-v6-c-button pf-m-primary").
							Type("Button").
							OnClick(func(ctx app.Context, e app.Event) {
								// Prefer selected path, fall back to current path if not selected
								// We have to do this here again to prevent a race condition for c.operationSelectedPath
								newDirectory := c.operationSelectedPath
								if newDirectory == "" {
									newDirectory = c.OperationCurrentPath
								}

								newPath := path.Join(newDirectory, path.Base(c.selectedPath))

								c.MovePath(c.selectedPath, newPath)

								c.movePathModalOpen = false
								c.OperationSetCurrentPath("/")
								c.operationSelectedPath = ""
							}).
							Text(`Move to "` + path.Base(newDirectory) + `"`)
					}(),
					app.Button().
						Class("pf-v6-c-button pf-m-link").
						Type("button").
						OnClick(func(ctx app.Context, e app.Event) {
							c.movePathModalOpen = false
							c.OperationSetCurrentPath("/")
							c.operationSelectedPath = ""
						}).
						Text("Cancel"),
				},
			},

			&Modal{
				Open: c.copyPathModalOpen,
				Close: func() {
					c.copyPathModalOpen = false
				},

				ID:     "copy-path-modal-title",
				Large:  true,
				Nested: c.Nested,

				Title: `Copy "` + path.Base(c.selectedPath) + `"`,
				Body: []app.UI{
					&PathPickerToolbar{
						Index:        c.OperationIndex,
						RefreshIndex: c.RefreshIndex,

						PathComponents: operationPathComponents,

						CurrentPath:    c.OperationCurrentPath,
						SetCurrentPath: c.OperationSetCurrentPath,

						SelectedPath: c.operationSelectedPath,
						SetSelectedPath: func(s string) {
							c.operationSelectedPath = s
						},

						OpenCreateDirectoryModal: func() {
							c.createDirectoryModalOpen = true
						},
					},
					app.If(
						len(c.OperationIndex) > 0,
						func() app.UI {
							return &FileGrid{
								Index: c.OperationIndex,

								SelectedPath: c.operationSelectedPath,
								SetSelectedPath: func(s string) {
									c.operationSelectedPath = s
								},

								CurrentPath:    c.OperationCurrentPath,
								SetCurrentPath: c.OperationSetCurrentPath,

								Standalone: true,
							}
						},
					).Else(
						func() app.UI {
							return &EmptyState{}
						},
					),
				},
				Footer: []app.UI{
					func() app.UI {
						// Prefer selected path, fall back to current path if not selected
						newDirectory := c.operationSelectedPath
						if newDirectory == "" {
							newDirectory = c.OperationCurrentPath
						}

						return app.Button().
							Class("pf-v6-c-button pf-m-primary").
							Type("Button").
							OnClick(func(ctx app.Context, e app.Event) {
								// Prefer selected path, fall back to current path if not selected
								// We have to do this here again to prevent a race condition for c.operationSelectedPath
								newDirectory := c.operationSelectedPath
								if newDirectory == "" {
									newDirectory = c.OperationCurrentPath
								}

								newPath := path.Join(newDirectory, path.Base(c.selectedPath))

								c.CopyPath(c.selectedPath, newPath)

								c.copyPathModalOpen = false
								c.OperationSetCurrentPath("/")
								c.operationSelectedPath = ""
							}).
							Text(`Copy to "` + path.Base(newDirectory) + `"`)
					}(),
					app.Button().
						Class("pf-v6-c-button pf-m-link").
						Type("button").
						OnClick(func(ctx app.Context, e app.Event) {
							c.copyPathModalOpen = false
							c.OperationSetCurrentPath("/")
							c.operationSelectedPath = ""
						}).
						Text("Cancel"),
				},
			},

			&Modal{
				Open: c.uploadModalOpen,
				Close: func() {
					c.uploadModalOpen = false
				},

				ID:     "upload-modal-title",
				Nested: c.Nested,

				Title: "Upload",
				Body: []app.UI{
					app.Form().
						Class("pf-v6-c-form").
						ID("upload").
						OnSubmit(func(ctx app.Context, e app.Event) {
							e.PreventDefault()

							reader := app.Window().JSValue().Get("FileReader").New()
							input := app.Window().GetElementByID("upload-file-input")
							fileName := input.Get("files").Get("0").Get("name").String()

							reader.Set("onload", app.FuncOf(func(this app.Value, args []app.Value) interface{} {
								go func() {
									rawFileContent := app.Window().Get("Uint8Array").New(args[0].Get("target").Get("result"))

									fileContent := make([]byte, rawFileContent.Get("length").Int())
									app.CopyBytesToGo(fileContent, rawFileContent)

									c.WriteToPath(filepath.Join(c.CurrentPath, fileName), fileContent)

									// Manually refresh, as `c.WriteToPath` runs in a seperate goroutine
									ctx.Async(func() {
										c.RefreshIndex()
									})
								}()

								return nil
							}))

							reader.Call("readAsArrayBuffer", input.Get("files").Get("0"))

							// Clear the input
							input.Set("value", app.Null())

							c.uploadModalOpen = false
						}).
						Body(
							&FormGroup{
								Label: app.Label().
									For("upload-file-input").
									Class("pf-v6-c-form__label").
									Body(
										app.
											Span().
											Class("pf-v6-c-form__label-text").
											Text("File to upload"),
									),
								Input: app.Input().
									Name("upload-file-input").
									ID("upload-file-input").
									Type("file").
									Required(true).
									Class("pf-v6-c-form-control"),
							},
						),
				},
				Footer: []app.UI{
					app.Button().
						Class("pf-v6-c-button pf-m-primary").
						Type("submit").
						Form("upload").
						Text("Upload"),
					app.Button().
						Class("pf-v6-c-button pf-m-link").
						Type("button").
						OnClick(func(ctx app.Context, e app.Event) {
							// Clear the input
							app.Window().GetElementByID("upload-file-input").Set("value", app.Null())

							c.uploadModalOpen = false
						}).
						Text("Cancel"),
				},
			},

			&Modal{
				Open: c.editModalOpen,
				Close: func() {
					if c.textEditorDirty() {
						c.discardEditsModalOpen = true
					} else {
						c.discardEdits()
					}
				},

				ID:     "edit-modal-title",
				Nested: c.Nested,
				Large:  true,

				Title: `Editing "` + path.Base(c.selectedPath) + `"`,
				Body: []app.UI{
					&TextEditor{
						Content: c.EditPathContents,
						SetContent: func(s string) {
							if c.cleanEditPathContents == "" {
								c.cleanEditPathContents = c.EditPathContents
							}

							c.SetEditPathContents(s)
						},

						Refresh: func() {
							if c.textEditorDirty() {
								// When discarding edits, refresh
								c.discardEditsModalTargetsRefresh = true

								c.discardEditsModalOpen = true
							} else {
								c.editPath()
							}
						},
						Save: func() {
							c.WriteToPath(c.selectedPath, []byte(c.EditPathContents))
						},

						Language:       selectedPathContentType,
						VariableHeight: true,
					},
				},
				Footer: []app.UI{
					app.Button().
						Class("pf-v6-c-button pf-m-primary").
						OnClick(func(ctx app.Context, e app.Event) {
							c.saveEdits()
						}).
						Text("Save and close"),
					app.Button().
						Class("pf-v6-c-button pf-m-link").
						Type("button").
						OnClick(func(ctx app.Context, e app.Event) {
							if c.textEditorDirty() {
								c.discardEditsModalOpen = true
							} else {
								c.discardEdits()
							}
						}).
						Text("Cancel"),
				},
			},

			&Modal{
				Open: c.discardEditsModalOpen,
				Close: func() {
					c.discardEditsModalOpen = false
				},

				ID:     "discard-edits-modal-title",
				Nested: c.Nested,

				Title: `Discard changes in "` + path.Base(c.selectedPath) + `"?`,
				Body: []app.UI{
					app.P().Text(`If you discard the changes, they will be permanently lost.`),
				},
				Footer: []app.UI{
					app.Button().
						Class("pf-v6-c-button pf-m-danger").
						Type("button").
						OnClick(func(ctx app.Context, e app.Event) {
							if c.discardEditsModalTargetsRefresh {
								c.editPath()
								c.discardEditsModalOpen = false
							} else {
								c.discardEdits()
							}
						}).
						Text(func() string {
							if c.discardEditsModalTargetsRefresh {
								return "Discard and refresh"
							}

							return "Discard"
						}()),
					app.Button().
						Class("pf-v6-c-button pf-m-link").
						Type("button").
						OnClick(func(ctx app.Context, e app.Event) {
							c.discardEditsModalOpen = false
						}).
						Text("Cancel"),
				},
			},
		)
}

func (c *FileExplorer) closeOverflowMenus() {
	c.overflowMenuOpen = false
}

func (c *FileExplorer) sharePath() {
	c.closeOverflowMenus()

	c.SharePath(c.selectedPath)

	c.sharePathModalOpen = true
}

func (c *FileExplorer) deleteFile() {
	c.closeOverflowMenus()

	c.deletionConfirmModalOpen = true
}

func (c *FileExplorer) moveTo() {
	c.closeOverflowMenus()

	// Preseed the file picker and name input with the current path
	c.OperationSetCurrentPath(path.Dir(c.selectedPath))

	// Check if the selected item is a file
	selectedItemIsFile := false
	for _, part := range c.Index {
		if part.Name() == path.Base(c.selectedPath) && !part.IsDir() {
			selectedItemIsFile = true
		}
	}

	// Don't select the item as a destination if it is a file
	if selectedItemIsFile {
		c.operationSelectedPath = ""
	} else {
		c.operationSelectedPath = c.selectedPath
	}

	// Close the overflow menu
	c.overflowMenuOpen = false

	// Open the modal
	c.movePathModalOpen = true
}

func (c *FileExplorer) copyTo() {
	c.closeOverflowMenus()

	// Preseed the file picker and name input with the current path
	c.OperationSetCurrentPath(path.Dir(c.selectedPath))

	// Check if the selected item is a file
	selectedItemIsFile := false
	for _, part := range c.Index {
		if part.Name() == path.Base(c.selectedPath) && !part.IsDir() {
			selectedItemIsFile = true
		}
	}

	// Don't select the item as a destination if it is a file
	if selectedItemIsFile {
		c.operationSelectedPath = ""
	} else {
		c.operationSelectedPath = c.selectedPath
	}

	// Close the overflow menu
	c.overflowMenuOpen = false

	// Open the modal
	c.copyPathModalOpen = true
}

func (c *FileExplorer) rename() {
	c.closeOverflowMenus()

	// Preseed the input with the current file name
	c.newFileName = path.Base(c.selectedPath)

	c.renamePathModalOpen = true
}

func (c *FileExplorer) createDirectory() {
	c.closeOverflowMenus()

	c.createDirectoryModalOpen = true
}

func (c *FileExplorer) createEmptyFile() {
	c.closeOverflowMenus()

	c.createEmptyFileModalOpen = true
}

func (c *FileExplorer) uploadFile() {
	c.closeOverflowMenus()

	c.uploadModalOpen = true
}

func (c *FileExplorer) refresh() {
	c.closeOverflowMenus()

	c.RefreshIndex()
}

func (c *FileExplorer) mountDirectory() {
	c.closeOverflowMenus()

	c.mountFolderModalOpen = true
}

func (c *FileExplorer) editPath() {
	c.closeOverflowMenus()

	// Fetch the contents of the path to be edited
	c.EditPath(c.selectedPath)

	// Track the contents so that the "dirty" state can be changed
	c.cleanEditPathContents = ""

	c.editModalOpen = true
}

func (c *FileExplorer) textEditorDirty() bool {
	return c.cleanEditPathContents != "" && c.EditPathContents != c.cleanEditPathContents
}

func (c *FileExplorer) discardEdits() {
	c.discardEditsModalOpen = false
	c.editModalOpen = false
	c.SetEditPathContents("")
	c.cleanEditPathContents = ""
	c.discardEditsModalTargetsRefresh = false
}

func (c *FileExplorer) saveEdits() {
	c.WriteToPath(c.selectedPath, []byte(c.EditPathContents))
	c.editModalOpen = false
	c.SetEditPathContents("")
	c.cleanEditPathContents = ""
	c.discardEditsModalTargetsRefresh = false
}
