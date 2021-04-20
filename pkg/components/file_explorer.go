package components

import (
	"os"
	"path/filepath"

	"github.com/maxence-charriere/go-app/v8/pkg/app"
	"github.com/pojntfx/liwasc/pkg/components"
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
	return app.Div().
		Body(
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

									c.Update()
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
							&components.Controlled{
								Component: app.Input().
									Type("text").
									Value(c.newDirectoryName).
									OnInput(func(ctx app.Context, e app.Event) {
										c.newDirectoryName = ctx.JSSrc.Get("value").String()

										c.Update()
									}),
								Properties: map[string]interface{}{
									"value": c.newDirectoryName,
								},
							},
							app.Button().
								OnClick(func(ctx app.Context, e app.Event) {
									c.CreatePath(filepath.Join(c.CurrentPath, c.newDirectoryName))

									c.newDirectoryName = ""

									c.Update()
								}).
								Text("Create Directory"),
						),
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

										c.WriteToPath(fileName, fileContent)

										c.RefreshIndex()
									}()

									return nil
								}))

								reader.Call("readAsArrayBuffer", ctx.JSSrc.Get("files").Get("0"))
							}),
						app.If(
							c.selectedPath != "",
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
							app.Button().
								OnClick(func(ctx app.Context, e app.Event) {
									c.DeletePath(c.selectedPath)
								}).
								Text("Delete"),
							app.Div().Body(
								app.Input().
									Type("text").
									OnInput(func(ctx app.Context, e app.Event) {
										c.pathToMoveTo = ctx.JSSrc.Get("value").String()

										c.Update()
									}),
								app.Button().
									OnClick(func(ctx app.Context, e app.Event) {
										c.MovePath(c.selectedPath, c.pathToMoveTo)

										c.pathToMoveTo = ""

										c.Update()

										c.RefreshIndex()
									}).
									Text("Move"),
							),
							app.Div().Body(
								app.Input().
									Type("text").
									OnInput(func(ctx app.Context, e app.Event) {
										c.pathToCopyTo = ctx.JSSrc.Get("value").String()

										c.Update()
									}),
								app.Button().
									OnClick(func(ctx app.Context, e app.Event) {
										c.CopyPath(c.selectedPath, c.pathToCopyTo)

										c.pathToCopyTo = ""

										c.Update()

										c.RefreshIndex()
									}).
									Text("Copy"),
							),
							app.Div().Body(
								app.Input().
									Type("text").
									OnInput(func(ctx app.Context, e app.Event) {
										c.newFileName = ctx.JSSrc.Get("value").String()

										c.Update()
									}),
								app.Button().
									OnClick(func(ctx app.Context, e app.Event) {
										c.MovePath(c.selectedPath, filepath.Join(c.CurrentPath, c.newFileName))

										c.newFileName = ""

										c.Update()

										c.RefreshIndex()
									}).
									Text("Rename"),
							),
						),
					),
				app.Div().
					Body(
						&components.Controlled{
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

										c.Update()
									}).
									OnDblClick(func(ctx app.Context, e app.Event) {
										if c.Index[i].IsDir() {
											e.PreventDefault()

											c.selectedPath = ""

											c.Update()

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
