package components

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type FileExplorer struct {
	app.Compo

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

	Error   error
	Recover func()
	Ignore  func()

	selectedPath     string
	newDirectoryName string
	pathToCopyTo     string
	newFileName      string

	overflowMenuOpen bool

	mountFolderModalOpen     bool
	sharePathModalOpen       bool
	createDirectoryModalOpen bool
	deletionConfirmModalOpen bool
	renamePathModalOpen      bool
	movePathModalOpen        bool

	operationSelectedPath string
}

func (c *FileExplorer) Render() app.UI {
	rawPathParts := strings.Split(c.CurrentPath, string(os.PathSeparator))
	pathParts := []string{}
	for _, pathPart := range rawPathParts {
		// Ignore empty paths
		if pathPart != "" {
			pathParts = append(pathParts, pathPart)
		}
	}

	return app.Div().
		Body(
			app.Div().
				Class("pf-c-card").
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
												Class("pf-c-toolbar__content-section").
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
																							app.Nav().
																								Class("pf-c-breadcrumb").
																								Aria("label", "Current path").
																								Body(
																									app.Ol().
																										Class("pf-c-breadcrumb__list").
																										Body(
																											app.Li().
																												Class("pf-c-breadcrumb__item").
																												Body(
																													app.Span().
																														Class("pf-c-breadcrumb__item-divider").
																														Body(
																															app.I().
																																Class("fas fa-angle-right").
																																Aria("hidden", true),
																														),
																													app.Button().
																														Type("button").
																														Class("pf-c-breadcrumb__link").
																														TabIndex(0).
																														OnClick(func(ctx app.Context, e app.Event) {
																															c.SetCurrentPath("/")

																															c.selectedPath = ""
																														}).
																														Text("Files"),
																												),
																											app.Range(pathParts).Slice(func(i int) app.UI {
																												link := path.Join(append([]string{"/"}, pathParts[:i+1]...)...)

																												// The last path part shouldn't be marked as a link
																												classes := "pf-c-breadcrumb__link"
																												if i == len(pathParts)-1 {
																													classes += " pf-m-current"
																												}

																												return app.Li().
																													Class("pf-c-breadcrumb__item").
																													Body(
																														app.Span().
																															Class("pf-c-breadcrumb__item-divider").
																															Body(
																																app.I().
																																	Class("fas fa-angle-right").
																																	Aria("hidden", true),
																															),
																														app.If(
																															// The last path part shouldn't be an action
																															i == len(pathParts)-1,
																															app.A().
																																Class(classes).
																																Text(pathParts[i]),
																														).Else(
																															app.Button().
																																Type("button").
																																Class(classes).
																																OnClick(func(ctx app.Context, e app.Event) {
																																	c.SetCurrentPath(link)

																																	c.selectedPath = ""
																																}).
																																Text(pathParts[i]),
																														),
																													)
																											}),
																										),
																								),
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
																		Class("pf-c-pagination pf-m-compact pf-m-compact").
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
																									app.If(
																										c.selectedPath != "",
																										app.Div().Class("pf-c-overflow-menu__item").
																											Body(
																												app.Button().
																													Type("button").
																													Aria("label", "Share file").
																													Title("Share file").
																													Class("pf-c-button pf-m-plain").
																													OnClick(func(ctx app.Context, e app.Event) {
																														c.SharePath(c.selectedPath)

																														c.sharePathModalOpen = true
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
																														c.deletionConfirmModalOpen = true
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
																																			OnClick(func(ctx app.Context, e app.Event) {
																																				// Preseed the file picker and name input with the current path
																																				c.OperationSetCurrentPath(path.Dir(c.selectedPath))
																																				c.operationSelectedPath = c.selectedPath

																																				// Close the overflow menu
																																				c.overflowMenuOpen = false

																																				// Open the modal
																																				c.movePathModalOpen = true
																																			}).
																																			Text("Move to ..."),
																																	),
																																app.Li().
																																	Body(
																																		app.Button().
																																			Class("pf-c-dropdown__menu-item").
																																			Text("Copy to ..."),
																																	),
																																app.Li().
																																	Class("pf-c-divider").
																																	Aria("role", "separator"),
																																app.Li().
																																	Body(
																																		app.Button().
																																			Class("pf-c-dropdown__menu-item").
																																			OnClick(func(ctx app.Context, e app.Event) {
																																				// Preseed the input with the current file name
																																				c.newFileName = path.Base(c.selectedPath)

																																				// Close the overflow menu
																																				c.overflowMenuOpen = false

																																				c.renamePathModalOpen = true
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
																													c.createDirectoryModalOpen = true
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
																												Body(
																													app.I().
																														Class("fas fa-cloud-upload-alt").
																														Aria("hidden", true),
																												),
																										),
																								),
																						),
																					app.Div().
																						Class("pf-c-divider pf-m-vertical pf-m-inset-md").
																						Aria("role", "separator"),
																					app.Div().
																						Class("pf-c-overflow-menu__group pf-m-button-group").
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
																											c.RefreshIndex()
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
																											c.mountFolderModalOpen = true
																										}).
																										Body(
																											app.Span().
																												Class("pf-c-button__icon pf-m-start").
																												Body(
																													app.I().
																														Class("fas fa-hdd").
																														Aria("hidden", true),
																												),
																											app.Text("Mount Directory"),
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
								app.Div().
									Class("pf-c-empty-state").
									Body(
										app.Div().
											Class("pf-c-empty-state__content").
											Body(
												app.I().
													Class("fas fa-folder-open pf-c-empty-state__icon").
													Aria("hidden", true),
												app.H1().
													Class("pf-c-title pf-m-lg").
													Text("No files or directories here yet"),
												app.Div().
													Class("pf-c-empty-state__body").
													Text("You can upload a file or create a directory to make it available for nodes."),
												app.Button().
													Class("pf-c-button pf-m-primary").
													Type("button").
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
											),
									),
							),
						),
				),
			app.Div().Body(
				// Upload file
				app.Input().
					Type("file").
					OnChange(func(ctx app.Context, e app.Event) {
						reader := app.Window().JSValue().Get("FileReader").New()
						fileName := ctx.JSSrc.Get("files").Get("0").Get("name").String()

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

						reader.Call("readAsArrayBuffer", ctx.JSSrc.Get("files").Get("0"))
					}),
				app.If(
					c.selectedPath != "",
					// Copy
					app.Div().Body(
						&Controlled{
							Component: app.Input().
								Type("text").
								Value(c.pathToCopyTo).
								OnInput(func(ctx app.Context, e app.Event) {
									c.pathToCopyTo = ctx.JSSrc.Get("value").String()
								}),
							Properties: map[string]interface{}{
								"value": c.pathToCopyTo,
							},
						},
						app.Button().
							OnClick(func(ctx app.Context, e app.Event) {
								c.CopyPath(c.selectedPath, filepath.Join(c.CurrentPath, c.pathToCopyTo))

								c.pathToCopyTo = ""
							}).
							Text("Copy"),
					),
				),
			),
			&Modal{
				Open: c.mountFolderModalOpen,
				Close: func() {
					c.mountFolderModalOpen = false

					// This manual update is required as the event is fired from `app.Window`
					c.Update()
				},

				ID: "mount-folder-modal-title",

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
																Value(c.WebDAVAddress).
																Aria("label", "WebDAV server address").
																Name("webdav-address").
																ID("webdav-address"),
															Properties: map[string]interface{}{
																"value": c.WebDAVAddress,
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

					// This manual update is required as the event is fired from `app.Window`
					c.Update()
				},

				ID: "share-path-modal-title",

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
																Value(c.HTTPShareLink).
																Aria("label", "HTTP address").
																Name("http-address").
																ID("http-address"),
															Properties: map[string]interface{}{
																"value": c.HTTPShareLink,
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
																Value(c.TFTPShareLink).
																Aria("label", "TFTP address").
																Name("tftp-address").
																ID("tftp-address"),
															Properties: map[string]interface{}{
																"value": c.TFTPShareLink,
															},
														},
														ID: "tftp-address",
													},
												),
										),
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

					// This manual update is required as the event is fired from `app.Window`
					c.Update()
				},

				ID: "create-directory-modal-title",

				Title: "Create Directory",
				Body: []app.UI{
					app.Form().
						Class("pf-c-form").
						ID("create-directory").
						OnSubmit(func(ctx app.Context, e app.Event) {
							e.PreventDefault()

							c.CreatePath(filepath.Join(c.CurrentPath, c.newDirectoryName))

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
												c.newDirectoryName = ctx.JSSrc.Get("value").String()
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

					// This manual update is required as the event is fired from `app.Window`
					c.Update()
				},

				ID: "deletion-confirm-modal-title",

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

					// This manual update is required as the event is fired from `app.Window`
					c.Update()
				},

				ID: "rename-path-modal-title",

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
												c.newFileName = ctx.JSSrc.Get("value").String()
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

					// This manual update is required as the event is fired from `app.Window`
					c.Update()
				},

				ID:    "move-path-modal-title",
				Large: true,

				Title: `Move "` + path.Base(c.selectedPath) + `"`,
				Body: []app.UI{
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
		)
}
