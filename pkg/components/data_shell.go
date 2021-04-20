package components

import (
	"os"

	"github.com/maxence-charriere/go-app/v8/pkg/app"
)

type DataShell struct {
	app.Compo

	AuthorizedWebDAVURL string
	Index               []os.FileInfo
	CurrentDir          string

	SetCurrentDir func(string, app.Context)
	UploadFile    func(string, []byte)
	Refresh       func(app.Context)

	Error   error
	Recover func()
	Ignore  func()

	// Config file editor
	ConfigFile    string
	SetConfigFile func(string)

	FormatConfigFile  func()
	RefreshConfigFile func()
	SaveConfigFile    func()

	ConfigFileError        error
	RecoverConfigFileError func()
	IgnoreConfigFileError  func()
}

func (c *DataShell) Render() app.UI {
	return app.Div().Body(
		app.Section().
			Body(
				&ConfigFileEditor{
					ConfigFile:    c.ConfigFile,
					SetConfigFile: c.SetConfigFile,

					FormatConfigFile:  c.FormatConfigFile,
					RefreshConfigFile: c.RefreshConfigFile,
					SaveConfigFile:    c.SaveConfigFile,

					Error:   c.ConfigFileError,
					Recover: c.RecoverConfigFileError,
					Ignore:  c.IgnoreConfigFileError,
				},
			),
		// app.Section().
		// 	Body(
		// 		app.Input().
		// 			ReadOnly(true).
		// 			Value(
		// 				c.AuthorizedWebDAVURL,
		// 			),
		// 	),
		// app.Section().
		// 	Body(
		// 		app.Input().
		// 			Type("file").
		// 			OnChange(func(ctx app.Context, e app.Event) {
		// 				reader := app.Window().JSValue().Get("FileReader").New()
		// 				fileName := ctx.JSSrc.Get("files").Get("0").Get("name").String()

		// 				reader.Set("onload", app.FuncOf(func(this app.Value, args []app.Value) interface{} {
		// 					go func() {
		// 						rawFileContent := app.Window().Get("Uint8Array").New(args[0].Get("target").Get("result"))

		// 						fileContent := make([]byte, rawFileContent.Get("length").Int())
		// 						app.CopyBytesToGo(fileContent, rawFileContent)

		// 						c.UploadFile(fileName, fileContent)

		// 						c.Refresh(ctx)
		// 					}()

		// 					return nil
		// 				}))

		// 				reader.Call("readAsArrayBuffer", ctx.JSSrc.Get("files").Get("0"))
		// 			}),
		// 	),
		// app.Section().
		// 	Body(
		// 		app.Ul().Body(
		// 			app.Range(c.Index).Slice(func(i int) app.UI {
		// 				handler := func(app.Context) {}
		// 				if c.Index[i].IsDir() {
		// 					handler = func(ctx app.Context) {
		// 						c.SetCurrentDir(c.Index[i].Name(), ctx)
		// 					}
		// 				}

		// 				return app.Li().
		// 					OnClick(func(ctx app.Context, e app.Event) {
		// 						handler(ctx)
		// 					}).
		// 					Body(
		// 						app.Text(c.Index[i].Name()),
		// 						app.If(
		// 							c.Index[i].IsDir(),
		// 							app.Text("/"),
		// 						),
		// 					)
		// 			}),
		// 		),
		// 	),
	)
}
