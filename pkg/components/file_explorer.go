package components

import (
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/studio-b12/gowebdav"
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

	CreatePath func(string)
	DeletePath func(string)
	MovePath   func(string, string)
	CopyPath   func(string, string)

	EditPathContents string
	EditPath         func(string)

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

	selectedPath     string
	newDirectoryName string
	newFileName      string

	overflowMenuOpen bool

	mountFolderModalOpen     bool
	sharePathModalOpen       bool
	createDirectoryModalOpen bool
	deletionConfirmModalOpen bool
	renamePathModalOpen      bool
	movePathModalOpen        bool
	copyPathModalOpen        bool
	uploadModalOpen          bool
	editModalOpen            bool

	operationSelectedPath string

	shareExpandableSectionOpen bool
	mountExpandableSectionOpen bool

	mobileMenuExpanded bool
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
			ctype := candidate.(gowebdav.File).ContentType()

			if ctype == "application/json" || strings.HasPrefix(ctype, "text/") {
				selectedPathContentType = ctype
			}

			break
		}
	}

	return app.Div().
		Class("pf-u-h-100").
		Body(
			app.Div().
				Class("pf-c-card pf-u-h-100").
				Body(
					app.Div().
						Class("pf-c-card__title").
						Body(
							app.Div().
								Class("pf-c-toolbar pf-u-pt-0").
								Body(
									app.Div().
										Class("pf-c-toolbar__content").
										Body(
											app.Div().
												Class("pf-c-toolbar__content-section pf-x-m-gap-md").
												Body(
													app.Div().
														Class("pf-c-toolbar__item pf-m-overflow-menu").
														Body(
															app.Div().
																Class("pf-c-overflow-menu").
																Body(
																	app.Div().
																		Class("pf-c-overflow-menu__content").
																		Body(
																			app.Div().
																				Class("pf-c-overflow-menu__group pf-m-button-group").
																				Body(
																					app.Div().
																						Class("pf-c-overflow-menu__item").
																						Body(
																							&Breadcrumbs{
																								PathComponents: pathComponents,

																								CurrentPath:    c.CurrentPath,
																								SetCurrentPath: c.SetCurrentPath,

																								SelectedPath: c.selectedPath,
																								SetSelectedPath: func(s string) {
																									c.selectedPath = s
																								},
																							},
																						),
																				),
																		),
																),
														),
													app.Div().
														Class("pf-c-toolbar__item pf-m-pagination").
														Body(
															app.Div().
																Class("pf-c-pagination pf-m-compact").
																Body(
																	app.Div().
																		Class("pf-c-pagination pf-m-compact").
																		Body(
																			app.Div().
																				Class("pf-c-overflow-menu").
																				Body(
																					app.Div().
																						Class("pf-c-overflow-menu__content pf-u-display-flex pf-u-display-none-on-lg").
																						Body(
																							app.Div().
																								Class("pf-c-overflow-menu__group pf-m-button-group").
																								Body(
																									app.Div().
																										Class("pf-c-overflow-menu__item").
																										Body(
																											app.Div().
																												Class(func() string {
																													classes := "pf-c-dropdown"

																													if c.mobileMenuExpanded {
																														classes += " pf-m-expanded"
																													}

																													return classes
																												}()).
																												Body(
																													app.Button().
																														Class("pf-c-dropdown__toggle pf-m-plain").
																														ID("page-default-nav-example-dropdown-kebab-1-button").
																														Aria("expanded", c.mobileMenuExpanded).Type("button").
																														Aria("label", "Actions").
																														Body(
																															app.I().
																																Class("fas fa-ellipsis-v pf-u-display-none-on-lg").
																																Aria("hidden", true),
																															app.I().
																																Class("fas fa-question-circle pf-u-display-none pf-u-display-inline-block-on-lg").
																																Aria("hidden", true),
																														).OnClick(func(ctx app.Context, e app.Event) {
																														c.mobileMenuExpanded = !c.mobileMenuExpanded
																													}),
																													app.Ul().
																														Class("pf-c-dropdown__menu pf-m-align-right").
																														Aria("aria-labelledby", "page-default-nav-example-dropdown-kebab-1-button").
																														Hidden(!c.mobileMenuExpanded).
																														Body(
																															app.If(
																																c.selectedPath != "",
																																app.If(
																																	selectedPathContentType != "",
																																	app.Li().
																																		Body(
																																			app.Button().
																																				Class("pf-c-button pf-c-dropdown__menu-item").
																																				Type("button").
																																				OnClick(func(ctx app.Context, e app.Event) {
																																					c.editPath()
																																				}).
																																				Text("Edit file"),
																																		),
																																),
																																app.Li().
																																	Body(
																																		app.Button().
																																			Class("pf-c-button pf-c-dropdown__menu-item").
																																			Type("button").
																																			OnClick(func(ctx app.Context, e app.Event) {
																																				c.sharePath()
																																			}).
																																			Text("Share file"),
																																	),
																																app.Li().
																																	Body(
																																		app.Button().
																																			Class("pf-c-button pf-c-dropdown__menu-item").
																																			Type("button").
																																			OnClick(func(ctx app.Context, e app.Event) {
																																				c.deleteFile()
																																			}).
																																			Text("Delete file"),
																																	),
																																app.Li().
																																	Class("pf-c-divider").
																																	Aria("role", "separator"),
																																app.Li().
																																	Body(
																																		app.Button().
																																			Class("pf-c-dropdown__menu-item").
																																			Type("button").
																																			OnClick(func(ctx app.Context, e app.Event) {
																																				c.moveTo()
																																			}).
																																			Text("Move to ..."),
																																	),
																																app.Li().
																																	Body(
																																		app.Button().
																																			Class("pf-c-dropdown__menu-item").
																																			Type("button").
																																			OnClick(func(ctx app.Context, e app.Event) {
																																				c.copyTo()
																																			}).
																																			Text("Copy to ..."),
																																	),
																																app.Li().
																																	Class("pf-c-divider").
																																	Aria("role", "separator"),
																																app.Li().
																																	Body(
																																		app.Button().
																																			Class("pf-c-dropdown__menu-item").
																																			Type("button").
																																			OnClick(func(ctx app.Context, e app.Event) {
																																				c.rename()
																																			}).
																																			Text("Rename"),
																																	),
																																app.Li().
																																	Class("pf-c-divider").
																																	Aria("role", "separator"),
																															),
																															app.Li().
																																Body(
																																	app.Button().
																																		Class("pf-c-button pf-c-dropdown__menu-item").
																																		Type("button").
																																		OnClick(func(ctx app.Context, e app.Event) {
																																			c.createDirectory()
																																		}).
																																		Text("Create directory"),
																																),
																															app.Li().
																																Body(
																																	app.Button().
																																		Class("pf-c-button pf-c-dropdown__menu-item").
																																		Type("button").
																																		OnClick(func(ctx app.Context, e app.Event) {
																																			c.uploadFile()
																																		}).
																																		Text("Upload file"),
																																),
																															app.Li().
																																Class("pf-c-divider").
																																Aria("role", "separator"),
																															app.Li().
																																Body(
																																	app.Button().
																																		Class("pf-c-button pf-c-dropdown__menu-item").
																																		Type("button").
																																		OnClick(func(ctx app.Context, e app.Event) {
																																			c.refresh()
																																		}).
																																		Text("Refresh"),
																																),
																															app.Li().
																																Body(
																																	app.Button().
																																		Class("pf-c-button pf-c-dropdown__menu-item").
																																		Type("button").
																																		OnClick(func(ctx app.Context, e app.Event) {
																																			c.mountDirectory()
																																		}).
																																		Text("Mount directory"),
																																),
																														),
																												),
																										),
																								),
																						),
																					app.Div().
																						Class("pf-c-overflow-menu__content pf-u-display-none pf-u-display-flex-on-lg").
																						Body(
																							app.Div().
																								Class("pf-c-overflow-menu__group pf-m-button-group").
																								Body(
																									app.If(
																										c.selectedPath != "",
																										app.If(
																											selectedPathContentType != "",
																											app.Div().Class("pf-c-overflow-menu__item").
																												Body(
																													app.Button().
																														Type("button").
																														Aria("label", "Edit file").
																														Title("Edit file").
																														Class("pf-c-button pf-m-plain").
																														OnClick(func(ctx app.Context, e app.Event) {
																															c.editPath()
																														}).
																														Body(
																															app.I().
																																Class("fas fa-edit").
																																Aria("hidden", true),
																														),
																												),
																										),
																										app.Div().Class("pf-c-overflow-menu__item").
																											Body(
																												app.Button().
																													Type("button").
																													Aria("label", "Share file").
																													Title("Share file").
																													Class("pf-c-button pf-m-plain").
																													OnClick(func(ctx app.Context, e app.Event) {
																														c.sharePath()
																													}).
																													Body(
																														app.I().
																															Class("fas fa-share-alt").
																															Aria("hidden", true),
																													),
																											),
																										app.Div().
																											Class("pf-c-overflow-menu__item").
																											Body(
																												app.Button().
																													Type("button").
																													Aria("label", "Delete file").
																													Title("Delete file").
																													Class("pf-c-button pf-m-plain").
																													OnClick(func(ctx app.Context, e app.Event) {
																														c.deleteFile()
																													}).
																													Body(
																														app.I().
																															Class("fas fa-trash").
																															Aria("hidden", true),
																													),
																											),
																										app.Div().
																											Class("pf-c-overflow-menu__control").
																											Body(
																												app.Div().
																													Class(func() string {
																														classes := "pf-c-dropdown"
																														if c.overflowMenuOpen {
																															classes += " pf-m-expanded"
																														}

																														return classes
																													}()).
																													Body(
																														app.Button().
																															Class("pf-c-dropdown__toggle pf-m-plain").
																															ID("toolbar-overflow-menu-button").
																															Aria("expanded", c.overflowMenuOpen).
																															Type("button").
																															Aria("label", "Toggle overflow menu").
																															OnClick(func(ctx app.Context, e app.Event) {
																																c.overflowMenuOpen = !c.overflowMenuOpen
																															}).
																															Body(
																																app.I().
																																	Class("fas fa-ellipsis-v").
																																	Aria("hidden", true),
																															),
																														app.Ul().
																															Class("pf-c-dropdown__menu").
																															Aria("labelledby", "toolbar-overflow-menu-button").
																															Hidden(!c.overflowMenuOpen).
																															Body(
																																app.Li().
																																	Body(
																																		app.Button().
																																			Class("pf-c-dropdown__menu-item").
																																			Type("button").
																																			OnClick(func(ctx app.Context, e app.Event) {
																																				c.moveTo()
																																			}).
																																			Text("Move to ..."),
																																	),
																																app.Li().
																																	Body(
																																		app.Button().
																																			Class("pf-c-dropdown__menu-item").
																																			Type("button").
																																			OnClick(func(ctx app.Context, e app.Event) {
																																				c.copyTo()
																																			}).
																																			Text("Copy to ..."),
																																	),
																																app.Li().
																																	Class("pf-c-divider").
																																	Aria("role", "separator"),
																																app.Li().
																																	Body(
																																		app.Button().
																																			Class("pf-c-dropdown__menu-item").
																																			Type("button").
																																			OnClick(func(ctx app.Context, e app.Event) {
																																				c.rename()
																																			}).
																																			Text("Rename"),
																																	),
																															),
																													),
																											),
																										app.Div().
																											Class("pf-c-divider pf-m-vertical pf-m-inset-md pf-u-mr-sm").
																											Aria("role", "separator"),
																									),
																									app.Div().
																										Class("pf-c-overflow-menu__item").
																										Body(
																											app.Button().
																												Type("button").
																												Aria("label", "Create directory").
																												Title("Create directory").
																												Class("pf-c-button pf-m-plain").
																												OnClick(func(ctx app.Context, e app.Event) {
																													c.createDirectory()
																												}).
																												Body(
																													app.I().
																														Class("fas fa-folder-plus").
																														Aria("hidden", true),
																												),
																										),
																									app.Div().
																										Class("pf-c-overflow-menu__item").
																										Body(
																											app.Button().
																												Type("button").
																												Aria("label", "Upload file").
																												Title("Upload file").
																												Class("pf-c-button pf-m-plain").
																												OnClick(func(ctx app.Context, e app.Event) {
																													c.uploadFile()
																												}).
																												Body(
																													app.I().
																														Class("fas fa-cloud-upload-alt").
																														Aria("hidden", true),
																												),
																										),
																								),
																						),
																					app.Div().
																						Class("pf-c-divider pf-m-vertical pf-m-inset-md pf-u-display-none pf-u-display-flex-on-lg").
																						Aria("role", "separator"),
																					app.Div().
																						Class("pf-c-overflow-menu__group pf-m-button-group pf-u-display-none pf-u-display-flex-on-lg").
																						Body(
																							app.Div().
																								Class("pf-c-overflow-menu__item").
																								Body(
																									app.Button().
																										Type("button").
																										Aria("label", "Refresh").
																										Title("Refresh").
																										Class("pf-c-button pf-m-plain").
																										OnClick(func(ctx app.Context, e app.Event) {
																											c.refresh()
																										}).
																										Body(
																											app.I().
																												Class("fas fas fa-sync").
																												Aria("hidden", true),
																										),
																								),
																							app.Div().
																								Class("pf-c-overflow-menu__item").
																								Body(
																									app.Button().
																										Class("pf-c-button pf-m-control").
																										Type("button").
																										OnClick(func(ctx app.Context, e app.Event) {
																											c.mountDirectory()
																										}).
																										Body(
																											app.Span().
																												Class("pf-c-button__icon pf-m-start").
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
												),
										),
								),
						),
					app.Div().
						Class("pf-c-card__body").
						Body(
							app.If(
								len(c.Index) > 0,
								&FileGrid{
									Index: c.Index,

									SelectedPath: c.selectedPath,
									SetSelectedPath: func(s string) {
										c.selectedPath = s
									},

									CurrentPath:    c.CurrentPath,
									SetCurrentPath: c.SetCurrentPath,
								},
							).Else(
								&EmptyState{
									Action: app.Button().
										Class("pf-c-button pf-m-primary").
										Type("button").
										OnClick(func(ctx app.Context, e app.Event) {
											c.uploadModalOpen = true
										}).
										Body(
											app.Span().
												Class("pf-c-button__icon pf-m-start").
												Body(
													app.I().
														Class("fas fa-cloud-upload-alt").
														Aria("hidden", true),
												),
											app.Text("Upload File"),
										),
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
						Class("pf-c-content").
						Body(
							app.P().
								Text(`You can mount this folder as a WebDAV share using your system's file explorer. To do so, please use the following credentials:`),
							app.Form().
								Class("pf-c-form").
								OnSubmit(func(ctx app.Context, e app.Event) {
									e.PreventDefault()
								}).
								Body(
									app.Div().
										Class("pf-c-form__group").
										Body(
											app.Div().
												Class("pf-c-form__group-label").
												Body(
													app.Label().
														Class("pf-c-form__label").
														For("webdav-address").
														Body(
															app.Span().
																Class("pf-c-form__label-text").
																Text("Address"),
														),
												),
											app.Div().
												Class("pf-c-form__group-control").
												Body(
													&CopyableInput{
														Component: &Controlled{
															Component: app.Input().
																Class("pf-c-form-control").
																ReadOnly(true).
																Type("text").
																Value(c.WebDAVAddress.String()).
																Aria("label", "WebDAV server address").
																Name("webdav-address").
																ID("webdav-address"),
															Properties: map[string]interface{}{
																"value": c.WebDAVAddress.String(),
															},
														},
														ID: "webdav-address",
													},
												),
										),
									app.Div().
										Class("pf-c-form__group").
										Body(
											app.Div().
												Class("pf-c-form__group-label").
												Body(
													app.Label().
														Class("pf-c-form__label").
														For("webdav-username").
														Body(
															app.Span().
																Class("pf-c-form__label-text").
																Text("Username"),
														),
												),
											app.Div().
												Class("pf-c-form__group-control").
												Body(
													&CopyableInput{
														Component: &Controlled{
															Component: app.Input().
																Class("pf-c-form-control").
																ReadOnly(true).
																Type("text").
																Value(c.WebDAVUsername).
																Aria("label", "WebDAV username").
																Name("webdav-username").
																ID("webdav-username"),
															Properties: map[string]interface{}{
																"value": c.WebDAVUsername,
															},
														},
														ID: "webdav-username",
													},
												),
										),
									app.Div().
										Class("pf-c-form__group").
										Body(
											app.Div().
												Class("pf-c-form__group-label").
												Body(
													app.Label().
														Class("pf-c-form__label").
														For("webdav-password").
														Body(
															app.Span().
																Class("pf-c-form__label-text").
																Text("Password (One-Time Token)"),
														),
												),
											app.Div().
												Class("pf-c-form__group-control").
												Body(
													&CopyableInput{
														Component: &Controlled{
															Component: app.Input().
																Class("pf-c-form-control").
																ReadOnly(true).
																Type("text").
																Value(c.WebDAVPassword).
																Aria("label", "WebDAV password").
																Name("webdav-password").
																ID("webdav-password"),
															Properties: map[string]interface{}{
																"value": c.WebDAVPassword,
															},
														},
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
												Class("pf-c-form pf-m-horizontal pf-u-mb-md").
												Body(
													&FormGroup{
														NoTopPadding: true,
														Label: app.Label().
															For("use-advertised-ip-for-webdav").
															Class("pf-c-form__label").
															Body(
																app.
																	Span().
																	Class("pf-c-form__label-text").
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
														NoTopPadding: true,
														Label: app.Label().
															For("use-davs").
															Class("pf-c-form__label").
															Body(
																app.
																	Span().
																	Class("pf-c-form__label-text").
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
						Class("pf-c-button pf-m-primary").
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
						Class("pf-c-content").
						Body(
							app.P().
								Text(`You should be able to access this path using the following addresses. Please note that, especially for the TFTP address, access is probably only possible from the LAN on which bofied is running:`),
							app.Form().
								Class("pf-c-form").
								OnSubmit(func(ctx app.Context, e app.Event) {
									e.PreventDefault()
								}).
								Body(
									app.Div().
										Class("pf-c-form__group").
										Body(
											app.Div().
												Class("pf-c-form__group-label").
												Body(
													app.Label().
														Class("pf-c-form__label").
														For("http-address").
														Body(
															app.Span().
																Class("pf-c-form__label-text").
																Text("HTTP Address"),
														),
												),
											app.Div().
												Class("pf-c-form__group-control").
												Body(
													&CopyableInput{
														Component: &Controlled{
															Component: app.Input().
																Class("pf-c-form-control").
																ReadOnly(true).
																Type("text").
																Value(c.HTTPShareLink.String()).
																Aria("label", "HTTP address").
																Name("http-address").
																ID("http-address"),
															Properties: map[string]interface{}{
																"value": c.HTTPShareLink.String(),
															},
														},
														ID: "http-address",
													},
												),
										),
									app.Div().
										Class("pf-c-form__group").
										Body(
											app.Div().
												Class("pf-c-form__group-label").
												Body(
													app.Label().
														Class("pf-c-form__label").
														For("tftp-address").
														Body(
															app.Span().
																Class("pf-c-form__label-text").
																Text("TFTP Address"),
														),
												),
											app.Div().
												Class("pf-c-form__group-control").
												Body(
													&CopyableInput{
														Component: &Controlled{
															Component: app.Input().
																Class("pf-c-form-control").
																ReadOnly(true).
																Type("text").
																Value(c.TFTPShareLink.String()).
																Aria("label", "TFTP address").
																Name("tftp-address").
																ID("tftp-address"),
															Properties: map[string]interface{}{
																"value": c.TFTPShareLink.String(),
															},
														},
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
												Class("pf-c-form pf-m-horizontal pf-u-mb-md").
												Body(
													&FormGroup{
														NoTopPadding: true,
														Label: app.Label().
															For("use-advertised-ip").
															Class("pf-c-form__label").
															Body(
																app.
																	Span().
																	Class("pf-c-form__label-text").
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
														NoTopPadding: true,
														Label: app.Label().
															For("use-https").
															Class("pf-c-form__label").
															Body(
																app.
																	Span().
																	Class("pf-c-form__label-text").
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
						Class("pf-c-button pf-m-primary").
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
						Class("pf-c-form").
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
									Class("pf-c-form__label").
									Body(
										app.
											Span().
											Class("pf-c-form__label-text").
											Text("Directory name"),
									),
								Input: &Controlled{
									Component: &Autofocused{
										Component: app.Input().
											Name("directory-name-input").
											ID("directory-name-input").
											Type("text").
											Required(true).
											Class("pf-c-form-control").
											OnInput(func(ctx app.Context, e app.Event) {
												c.newDirectoryName = ctx.JSSrc().Get("value").String()
											}),
									},
									Properties: map[string]interface{}{
										"value": c.newDirectoryName,
									},
								},
								Required: true,
							},
						),
				},
				Footer: []app.UI{
					app.Button().
						Class("pf-c-button pf-m-primary").
						Type("submit").
						Form("create-directory").
						Text("Create"),
					app.Button().
						Class("pf-c-button pf-m-link").
						Type("button").
						OnClick(func(ctx app.Context, e app.Event) {
							c.newDirectoryName = ""
							c.createDirectoryModalOpen = false
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
						Class("pf-c-button pf-m-danger").
						Type("button").
						OnClick(func(ctx app.Context, e app.Event) {
							c.DeletePath(c.selectedPath)

							c.selectedPath = ""
							c.deletionConfirmModalOpen = false
						}).
						Text("Delete"),
					app.Button().
						Class("pf-c-button pf-m-link").
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
						Class("pf-c-form").
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
									Class("pf-c-form__label").
									Body(
										app.
											Span().
											Class("pf-c-form__label-text").
											Text("New name"),
									),
								Input: &Controlled{
									Component: &Autofocused{
										Component: app.Input().
											Name("path-rename-input").
											ID("path-rename-input").
											Type("text").
											Required(true).
											Class("pf-c-form-control").
											OnInput(func(ctx app.Context, e app.Event) {
												c.newFileName = ctx.JSSrc().Get("value").String()
											}),
									},
									Properties: map[string]interface{}{
										"value": c.newFileName,
									},
								},
								Required: true,
							},
						),
				},
				Footer: []app.UI{
					app.Button().
						Class("pf-c-button pf-m-primary").
						Type("submit").
						Form("rename-path").
						Text("Rename"),
					app.Button().
						Class("pf-c-button pf-m-link").
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
						&FileGrid{
							Index: c.OperationIndex,

							SelectedPath: c.operationSelectedPath,
							SetSelectedPath: func(s string) {
								c.operationSelectedPath = s
							},

							CurrentPath:    c.OperationCurrentPath,
							SetCurrentPath: c.OperationSetCurrentPath,

							Standalone: true,
						},
					).Else(
						&EmptyState{},
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
							Class("pf-c-button pf-m-primary").
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
						Class("pf-c-button pf-m-link").
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
						&FileGrid{
							Index: c.OperationIndex,

							SelectedPath: c.operationSelectedPath,
							SetSelectedPath: func(s string) {
								c.operationSelectedPath = s
							},

							CurrentPath:    c.OperationCurrentPath,
							SetCurrentPath: c.OperationSetCurrentPath,

							Standalone: true,
						},
					).Else(
						&EmptyState{},
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
							Class("pf-c-button pf-m-primary").
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
						Class("pf-c-button pf-m-link").
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
						Class("pf-c-form").
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
									ctx.Emit(func() {
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
									Class("pf-c-form__label").
									Body(
										app.
											Span().
											Class("pf-c-form__label-text").
											Text("File to upload"),
									),
								Input: app.Input().
									Name("upload-file-input").
									ID("upload-file-input").
									Type("file").
									Required(true).
									Class("pf-c-form-control"),
							},
						),
				},
				Footer: []app.UI{
					app.Button().
						Class("pf-c-button pf-m-primary").
						Type("submit").
						Form("upload").
						Text("Upload"),
					app.Button().
						Class("pf-c-button pf-m-link").
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
					c.editModalOpen = false
				},

				ID:     "edit-modal-title",
				Nested: c.Nested,
				Large:  true,

				Title: `Editing "` + path.Base(c.selectedPath) + `"`,
				Body: []app.UI{
					&TextEditor{
						Content: c.EditPathContents,
						// SetContent: c.SetEditPathContents,

						Refresh: c.editPath,
						// Save:    c.SaveEditPathContents,

						Language: selectedPathContentType,
					},
				},
			},
		)
}

func (c *FileExplorer) sharePath() {
	// Close the overflow menus
	c.mobileMenuExpanded = false
	c.overflowMenuOpen = false

	c.SharePath(c.selectedPath)

	c.sharePathModalOpen = true
}

func (c *FileExplorer) deleteFile() {
	// Close the overflow menus
	c.mobileMenuExpanded = false
	c.overflowMenuOpen = false

	c.deletionConfirmModalOpen = true
}

func (c *FileExplorer) moveTo() {
	// Close the overflow menus
	c.mobileMenuExpanded = false
	c.overflowMenuOpen = false

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
	// Close the overflow menus
	c.mobileMenuExpanded = false
	c.overflowMenuOpen = false

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
	// Close the overflow menus
	c.mobileMenuExpanded = false
	c.overflowMenuOpen = false

	// Preseed the input with the current file name
	c.newFileName = path.Base(c.selectedPath)

	c.renamePathModalOpen = true
}

func (c *FileExplorer) createDirectory() {
	// Close the overflow menus
	c.mobileMenuExpanded = false
	c.overflowMenuOpen = false

	c.createDirectoryModalOpen = true
}

func (c *FileExplorer) uploadFile() {
	// Close the overflow menus
	c.mobileMenuExpanded = false
	c.overflowMenuOpen = false

	c.uploadModalOpen = true
}

func (c *FileExplorer) refresh() {
	// Close the overflow menus
	c.mobileMenuExpanded = false
	c.overflowMenuOpen = false

	c.RefreshIndex()
}

func (c *FileExplorer) mountDirectory() {
	// Close the overflow menus
	c.mobileMenuExpanded = false
	c.overflowMenuOpen = false

	c.mountFolderModalOpen = true
}

func (c *FileExplorer) editPath() {
	// Close the overflow menus
	c.mobileMenuExpanded = false
	c.overflowMenuOpen = false

	c.EditPath(c.selectedPath)
	c.editModalOpen = true
}
