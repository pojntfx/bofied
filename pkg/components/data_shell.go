package components

import (
	"os"
	"path/filepath"

	"github.com/maxence-charriere/go-app/v8/pkg/app"
	"github.com/pojntfx/liwasc/pkg/components"
)

type DataShell struct {
	app.Compo

	AuthorizedWebDAVURL string
	ConfigFile          string
	Index               []os.FileInfo
	CurrentDir          string

	SetConfigFile      func(string)
	ValidateConfigFile func()
	FormatConfigFile   func()
	SaveConfigFile     func()
	SetCurrentDir      func(string, app.Context)
	UploadFile         func(string, []byte)
	Refresh            func(app.Context)

	Error   error
	Recover func()
	Ignore  func()
}

func (c *DataShell) Render() app.UI {
	if c.Error != nil {
		return app.Div().
			Body(
				app.Code().Text(c.Error),
				app.Button().
					OnClick(func(ctx app.Context, e app.Event) {
						c.Recover()
					}).
					Text("Recover"),
				app.Button().
					OnClick(func(ctx app.Context, e app.Event) {
						c.Ignore()
					}).
					Text("Ignore"),
			)
	}

	return app.Div().Body(
		app.Section().
			Body(
				app.Input().
					ReadOnly(true).
					Value(
						c.AuthorizedWebDAVURL,
					),
			),
		app.Section().
			Body(
				&components.Controlled{
					Component: app.Textarea().
						OnInput(func(ctx app.Context, e app.Event) {
							c.SetConfigFile(ctx.JSSrc.Get("value").String())
						}).
						Text(
							c.ConfigFile,
						),
					Properties: map[string]interface{}{
						"value": c.ConfigFile,
					},
				},
				app.Br(),
				app.Button().
					OnClick(func(ctx app.Context, e app.Event) {
						c.ValidateConfigFile()
					}).
					Text("Validate"),
				app.Button().
					OnClick(func(ctx app.Context, e app.Event) {
						c.FormatConfigFile()
					}).
					Text("Format"),
				app.Button().
					OnClick(func(ctx app.Context, e app.Event) {
						c.SaveConfigFile()
					}).
					Text("Save"),
				app.If(
					c.CurrentDir != ".",
					app.Button().
						OnClick(func(ctx app.Context, e app.Event) {
							c.SetCurrentDir(filepath.Dir(c.CurrentDir), ctx)
						}).
						Text("Up"),
				),
				app.Button().
					OnClick(func(ctx app.Context, e app.Event) {
						c.Refresh(ctx)
					}).
					Text("Refresh"),
			),
		app.Section().
			Body(
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

								c.UploadFile(fileName, fileContent)

								c.Refresh(ctx)
							}()

							return nil
						}))

						reader.Call("readAsArrayBuffer", ctx.JSSrc.Get("files").Get("0"))
					}),
			),
		app.Section().
			Body(
				app.Ul().Body(
					app.Range(c.Index).Slice(func(i int) app.UI {
						handler := func(app.Context) {}
						if c.Index[i].IsDir() {
							handler = func(ctx app.Context) {
								c.SetCurrentDir(c.Index[i].Name(), ctx)
							}
						}

						return app.Li().
							OnClick(func(ctx app.Context, e app.Event) {
								handler(ctx)
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
	)
}
