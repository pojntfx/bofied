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
														app.Raw(`<div class="pf-c-pagination pf-m-compact">
        <div class="pf-c-pagination pf-m-compact pf-m-compact">
            <div class="pf-c-overflow-menu">
                <div class="pf-c-overflow-menu__content">
                    <div class="pf-c-overflow-menu__group pf-m-button-group">
                        <div class="pf-c-overflow-menu__item">
                            <button type="button" aria-label="Refresh" title="Refresh" class="pf-c-button pf-m-plain">
                                <i class="fas fas fa-sync" aria-hidden="true"></i>
                            </button>
                        </div>
                        <div class="pf-c-divider pf-m-vertical pf-m-inset-md pf-u-mr-sm" role="separator">
                        </div>
                        <div class="pf-c-overflow-menu__item">
                            <button type="button" aria-label="Create Directory" title="Create Directory"
                                class="pf-c-button pf-m-plain">
                                <i class="fas fa-folder-plus" aria-hidden="true"></i>
                            </button>
                        </div>
                        <div class="pf-c-overflow-menu__item">
                            <button type="button" aria-label="Upload file" title="Upload file"
                                class="pf-c-button pf-m-plain">
                                <i class="fas fa-cloud-upload-alt" aria-hidden="true"></i>
                            </button>
                        </div>
                        <div class="pf-c-divider pf-m-vertical pf-m-inset-md pf-u-mr-sm" role="separator">
                        </div>
                        <div class="pf-c-overflow-menu__item">
                            <button type="button" aria-label="Share file" title="Share file"
                                class="pf-c-button pf-m-plain">
                                <i class="fas fa-share-alt" aria-hidden="true"></i>
                            </button>
                        </div>
                        <div class="pf-c-overflow-menu__item">
                            <button type="button" aria-label="Delete file" title="Delete file"
                                class="pf-c-button pf-m-plain">
                                <i class="fas fa-trash" aria-hidden="true"></i>
                            </button>
                        </div>
                    </div>
                </div>
                <div class="pf-c-overflow-menu__control">
                    <div class="pf-c-dropdown">
                        <button class="pf-c-button pf-c-dropdown__toggle pf-m-plain" type="button"
                            aria-label="Overflow menu">
                            <i class="fas fa-ellipsis-v" aria-hidden="true"></i>
                        </button>
                        <ul class="pf-c-dropdown__menu"
                            aria-labelledby="toolbar-attribute-value-search-filter-desktop-example-overflow-menu-dropdown-toggle"
                            hidden>
                            <li>
                                <button class="pf-c-dropdown__menu-item">
                                    Tertiary
                                </button>
                            </li>
                        </ul>
                    </div>
                </div>
                <div class="pf-c-divider pf-m-vertical pf-m-inset-md pf-u-mr-lg" role="separator"></div>
                <div class="pf-c-overflow-menu__item">
                    <div class="pf-c-clipboard-copy">
                        <div class="pf-c-clipboard-copy__group">
                            <input class="pf-c-form-control" readonly type="text"
                                value="dav://user:eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IlFtSXdJNVprRjBJTmh0b0dwVWFmcyJ9.eyJuaWNrbmFtZSI6ImZlbGl4IiwibmFtZSI6ImZlbGl4QHBvanRpbmdlci5jb20iLCJwaWN0dXJlIjoiaHR0cHM6Ly9zLmdyYXZhdGFyLmNvbS9hdmF0YXIvZGI4NTZkZjMzZmE0ZjRiY2U0NDE4MTlmNjA0YzkwZDU_cz00ODAmcj1wZyZkPWh0dHBzJTNBJTJGJTJGY2RuLmF1dGgwLmNvbSUyRmF2YXRhcnMlMkZmZS5wbmciLCJ1cGRhdGVkX2F0IjoiMjAyMS0wNC0yMVQxNToxNzoxMy4zNDZaIiwiZW1haWwiOiJmZWxpeEBwb2p0aW5nZXIuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImlzcyI6Imh0dHBzOi8vcG9qbnRmeC5ldS5hdXRoMC5jb20vIiwic3ViIjoiYXV0aDB8NjA3YzJkNmY0OWZkODQwMDczODRmODZlIiwiYXVkIjoiMEdRWEt5dTlwUEtBeW1heFg4N29WUEkzNXJpMGt4Zk0iLCJpYXQiOjE2MTkwMTkwOTEsImV4cCI6MTYxOTA1NTA5MX0.LnSDaOEA1Do8DjYSK73GOuYzoD9gFTF7xvnnGPVUvXwJHLPCbHeiLLsL-ZMl4g80ErtYdmiIn1qDV7VEijepjGPfN-MoYlCy8Lml2EqMdy3ODxCd4CUj6Rx3ggsVXLxpZh6wutrgFLGNUeaiWFC2MAxjnItRVtAdwXHzvL4mIjOLfAeZcuighvhwfeGX7PHfUH1HHCoWvpjZVBN_wKC4A-vQyos4CDGGL5nvw2b86ND6QtpAIrKGXLFHqCqAjHfU1dKrWGvBMB15bju68RVguPr1NSQyPUgGWfNlDMf1hnSdyb5CVw0P3wD2R56jAcnQhi0RfDlq2t2pWmwDn-Emww@localhost:15256/private"
                                aria-label="Copyable input" />
                            <button class="pf-c-button pf-m-control" type="button" aria-label="Copy to clipboard"
                                aria-labelledby="basic-readonly-copy-button basic-readonly-text-input">
                                <i class="fas fa-copy" aria-hidden="true"></i>
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>`),
													),
												),
										),
								),
						),
					app.Div().
						Class("pf-c-card__body").
						Body(
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
