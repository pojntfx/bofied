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
																													app.A().
																														Class("pf-c-breadcrumb__link").
																														OnClick(func(ctx app.Context, e app.Event) {
																															c.selectedPath = ""

																															c.SetCurrentPath("/")
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
																															app.A().
																																Class(classes).
																																OnClick(func(ctx app.Context, e app.Event) {
																																	c.selectedPath = ""

																																	c.SetCurrentPath(link)
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
													app.Div().Class("pf-c-toolbar__item pf-m-pagination").Body(
														app.Div().Class("pf-c-pagination pf-m-compact").Body(
															app.Div().Class("pf-c-pagination pf-m-compact pf-m-compact").Body(
																app.Div().Class("pf-c-overflow-menu").Body(
																	app.Div().Class("pf-c-overflow-menu__content").Body(
																		app.Div().Class("pf-c-overflow-menu__group pf-m-button-group").Body(
																			app.Div().Class("pf-c-overflow-menu__item").Body(
																				app.Button().Type("button").Aria("label", "Refresh").Title("Refresh").Class("pf-c-button pf-m-plain").Body(
																					app.I().Class("fas fas fa-sync").Aria("hidden", true),
																				),
																			),
																			app.Div().Class("pf-c-divider pf-m-vertical pf-m-inset-md pf-u-mr-sm").Aria("role", "separator"),
																			app.Div().Class("pf-c-overflow-menu__item").Body(
																				app.Button().Type("button").Aria("label", "Create directory").Title("Create directory").Class("pf-c-button pf-m-plain").Body(
																					app.I().Class("fas fa-folder-plus").Aria("hidden", true),
																				),
																			),
																			app.Div().Class("pf-c-overflow-menu__item").Body(
																				app.Button().Type("button").Aria("label", "Upload file").Title("Upload file").Class("pf-c-button pf-m-plain").Body(
																					app.I().Class("fas fa-cloud-upload-alt").Aria("hidden", true),
																				),
																			),
																			app.Div().Class("pf-c-divider pf-m-vertical pf-m-inset-md pf-u-mr-sm").Aria("role", "separator"),
																			app.Div().Class("pf-c-overflow-menu__item").Body(
																				app.Button().Type("button").Aria("label", "Share file").Title("Share file").Class("pf-c-button pf-m-plain").Body(
																					app.I().Class("fas fa-share-alt").Aria("hidden", true),
																				),
																			),
																			app.Div().Class("pf-c-overflow-menu__item").Body(
																				app.Button().Type("button").Aria("label", "Delete file").Title("Delete file").Class("pf-c-button pf-m-plain").Body(
																					app.I().Class("fas fa-trash").Aria("hidden", true),
																				),
																			),
																		),
																	),
																	app.Div().Class("pf-c-overflow-menu__control").Body(
																		app.Raw(`<div class="pf-c-dropdown pf-m-expanded">
  <button class="pf-c-dropdown__toggle pf-m-plain" id="dropdown-kebab-expanded-button" aria-expanded="true" type="button" aria-label="Actions">
    <i class="fas fa-ellipsis-v" aria-hidden="true"></i>
  </button>
  <ul class="pf-c-dropdown__menu" aria-labelledby="dropdown-kebab-expanded-button">
    <li>
      <button class="pf-c-dropdown__menu-item" href="#">Move to ...</button>
    </li>
    <li>
      <button class="pf-c-dropdown__menu-item" type="button">Copy to ...</button>
    </li>
    <li class="pf-c-divider" role="separator"></li>
    <li>
      <button class="pf-c-dropdown__menu-item" href="#">Rename</button>
    </li>
  </ul>
</div>`),
																	),
																	app.Div().Class("pf-c-divider pf-m-vertical pf-m-inset-md pf-u-mr-lg").Aria("role", "separator"),
																	app.Div().Class("pf-c-overflow-menu__item").Body(
																		app.Div().Class("pf-c-clipboard-copy").Body(
																			app.Div().Class("pf-c-clipboard-copy__group").Body(
																				&Controlled{
																					Component: app.Input().Class("pf-c-form-control").ReadOnly(true).Type("text").Value(c.AuthorizedWebDAVURL).Aria("label", "Authorized WebDAV URL").ID("authorized-webdav-url-input"),
																					Properties: map[string]interface{}{
																						"value": c.AuthorizedWebDAVURL,
																					},
																				},
																				app.Button().Class("pf-c-button pf-m-control").Type("button").Aria("label", "Copy to clipboard").Aria("labelledby", "authorized-webdav-url-input").Body(
																					app.I().Class("fas fa-copy").Aria("hidden", true),
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

														c.selectedPath = ""

														c.SetCurrentPath(filepath.Join(c.CurrentPath, c.Index[i].Name()))
													}
												}).
												Body(
													app.Div().Class(
														func() string {
															classes := "pf-c-card pf-m-plain pf-m-hoverable pf-m-selectable"
															if c.selectedPath == filepath.Join(c.CurrentPath, c.Index[i].Name()) {
																classes += " pf-m-selected"
															}

															return classes
														}()).
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
								app.Div().Class("pf-c-empty-state").Body(
									app.Div().Class("pf-c-empty-state__content").Body(
										app.I().Class("fas fa-folder-open pf-c-empty-state__icon").Aria("hidden", true),
										app.H1().Class("pf-c-title pf-m-lg").Text("No files or directories here yet"),
										app.Div().Class("pf-c-empty-state__body").Text("You can upload a file or create a directory to make it available for nodes."),
										app.Button().Class("pf-c-button pf-m-primary").Type("button").Body(
											app.Span().Class("pf-c-button__icon pf-m-start").Body(
												app.I().Class("fas fa-cloud-upload-alt").Aria("hidden", true),
											),
											app.Text("Upload File"),
										),
									),
								),
							),
						),
				),
			app.Div().Body(
				app.Div().
					Body(
						// Path navigation
						app.H2().
							Text("Files"),
						// Current path
						app.Div().
							Body(
								app.Code().
									Text(c.CurrentPath),
							),
						// Set current path
						app.Div().Body(
							app.Input().
								Type("text").
								OnInput(func(ctx app.Context, e app.Event) {
									c.newCurrentPath = ctx.JSSrc.Get("value").String()
								}),
							app.Button().
								OnClick(func(ctx app.Context, e app.Event) {
									c.SetCurrentPath(c.newCurrentPath)
								}).
								Text("Navigate"),
						),
					),
				app.Div().
					Body(
						// Refresh
						app.Button().
							OnClick(func(ctx app.Context, e app.Event) {
								c.RefreshIndex()
							}).
							Text("Refresh"),
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
							// Delete
							app.Button().
								OnClick(func(ctx app.Context, e app.Event) {
									c.DeletePath(c.selectedPath)
								}).
								Text("Delete"),
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
				app.Div().
					Body(
						&Controlled{
							Component: app.Input().
								ReadOnly(true).
								Value(c.AuthorizedWebDAVURL),
							Properties: map[string]interface{}{
								"value": c.AuthorizedWebDAVURL,
							},
						},
					),
			),
			app.Div().
				Body(
					app.Ul().
						Body(
							app.Range(c.Index).Slice(func(i int) app.UI {
								return app.Li().
									Style("cursor", "pointer").
									Style("background", func() string {
										if c.selectedPath == filepath.Join(c.CurrentPath, c.Index[i].Name()) {
											return "lightgrey"
										}

										return "inherit"
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

											c.selectedPath = ""

											c.SetCurrentPath(filepath.Join(c.CurrentPath, c.Index[i].Name()))
										}
									}).
									Body(
										app.Text(c.Index[i].Name()),
										app.If(
											c.Index[i].IsDir(),
											app.Text("/"),
										),
									)
							}),
						),
				),
			app.If(
				c.Error != nil,
				app.Div().
					Body(
						app.H3().
							Text("Error"),
						app.Code().
							Text(c.Error),
						app.Button().
							OnClick(func(ctx app.Context, e app.Event) {
								c.Ignore()
							}).
							Text("Ignore"),
						app.Button().
							OnClick(func(ctx app.Context, e app.Event) {
								c.Recover()
							}).
							Text("Recover"),
					),
			),
		)
}
