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

	ShareLink string
	SharePath func(string)

	CreatePath func(string)
	DeletePath func(string)
	MovePath   func(string, string)
	CopyPath   func(string, string)

	AuthorizedWebDAVURL string

	Error   error
	Recover func()
	Ignore  func()

	newCurrentPath   string
	selectedPath     string
	newDirectoryName string
	pathToMoveTo     string
	pathToCopyTo     string
	newFileName      string

	overflowMenuOpen bool
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
																														c.DeletePath(c.selectedPath)

																														c.selectedPath = ""
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
																										Body(
																											app.Span().
																												Class("pf-c-button__icon pf-m-start").
																												Body(
																													app.I().
																														Class("fas fa-hdd").
																														Aria("hidden", true),
																												),
																											app.Text("Mount Folder"),
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
								app.Div().
									Class("pf-l-grid pf-m-gutter pf-m-all-4-col-on-sm pf-m-all-4-col-on-md pf-m-all-3-col-on-lg pf-m-all-2-col-on-xl").
									Body(
										app.Range(c.Index).Slice(func(i int) app.UI {
											return app.Div().
												Class("pf-l-grid__item pf-u-text-align-center").
												Body(
													app.Div().
														Class(
															func() string {
																classes := "pf-c-card pf-m-plain pf-m-hoverable pf-m-selectable"
																if c.selectedPath == filepath.Join(c.CurrentPath, c.Index[i].Name()) {
																	classes += " pf-m-selected"
																}

																return classes
															}()).
														OnClick(func(ctx app.Context, e app.Event) {
															newSelectedPath := filepath.Join(c.CurrentPath, c.Index[i].Name())
															if c.selectedPath == newSelectedPath {
																newSelectedPath = ""
															}

															c.selectedPath = newSelectedPath
														}).
														OnDblClick(func(ctx app.Context, e app.Event) {
															if c.Index[i].IsDir() {
																e.PreventDefault()

																c.SetCurrentPath(filepath.Join(c.CurrentPath, c.Index[i].Name()))

																c.selectedPath = ""
															}
														}).
														Aria("role", "button").
														TabIndex(0).
														Body(
															app.Div().
																Class("pf-c-card__body").
																Body(
																	app.I().
																		Class(func() string {
																			classes := "fas pf-u-font-size-3xl"
																			if c.Index[i].IsDir() {
																				classes += " fa-folder"
																			} else {
																				classes += " fa-file-alt"
																			}

																			return classes
																		}()).
																		Aria("hidden", true),
																),
															app.Div().
																Class("pf-c-card__footer").
																Body(
																	app.Text(c.Index[i].Name()),
																),
														),
												)
										}),
									),
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
				// Create directory
				app.Div().Body(
					&Controlled{
						Component: app.Input().
							Type("text").
							Value(c.newDirectoryName).
							OnInput(func(ctx app.Context, e app.Event) {
								c.newDirectoryName = ctx.JSSrc.Get("value").String()
							}),
						Properties: map[string]interface{}{
							"value": c.newDirectoryName,
						},
					},
					app.Button().
						OnClick(func(ctx app.Context, e app.Event) {
							c.CreatePath(filepath.Join(c.CurrentPath, c.newDirectoryName))

							c.newDirectoryName = ""
						}).
						Text("Create Directory"),
				),
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
					// Share
					app.Div().
						Body(
							app.Button().
								OnClick(func(ctx app.Context, e app.Event) {
									c.SharePath(c.selectedPath)
								}).
								Text("Share"),
							app.If(
								c.ShareLink != "",
								app.Div().
									Body(
										app.A().
											Target("_blank").
											Href(c.ShareLink).
											Text(c.ShareLink),
									),
							),
						),
					// Move
					app.Div().Body(
						&Controlled{
							Component: app.Input().
								Type("text").
								Value(c.pathToMoveTo).
								OnInput(func(ctx app.Context, e app.Event) {
									c.pathToMoveTo = ctx.JSSrc.Get("value").String()
								}),
							Properties: map[string]interface{}{
								"value": c.pathToMoveTo,
							},
						},
						app.Button().
							OnClick(func(ctx app.Context, e app.Event) {
								c.MovePath(c.selectedPath, filepath.Join(c.CurrentPath, c.pathToMoveTo))

								c.pathToMoveTo = ""
							}).
							Text("Move"),
					),
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
					// Rename
					app.Div().Body(
						&Controlled{
							Component: app.Input().
								Type("text").
								Value(c.newFileName).
								OnInput(func(ctx app.Context, e app.Event) {
									c.newFileName = ctx.JSSrc.Get("value").String()
								}),
							Properties: map[string]interface{}{
								"value": c.newFileName,
							},
						},
						app.Button().
							OnClick(func(ctx app.Context, e app.Event) {
								c.MovePath(c.selectedPath, filepath.Join(c.CurrentPath, c.newFileName))

								c.newFileName = ""
							}).
							Text("Rename"),
					),
				),
			),
		)
}
