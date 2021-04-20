package components

import (
	"os"
	"path/filepath"

	"github.com/maxence-charriere/go-app/v8/pkg/app"
)

type FileExplorer struct {
	app.Compo

	CurrentPath    string
	SetCurrentPath func(string)

	Index        []os.FileInfo
	RefreshIndex func()
	UploadFile   func(string, []byte)

	ShareLink string
	SharePath func(string)

	CreateDirectory func(string)
	DeletePath      func(string)
	MovePath        func(string, string)
	CopyPath        func(string, string)
	RenamePath      func(string, string)

	AuthorizedWebDAVURL string

	FileExplorerError        error
	RecoverFileExplorerError func()
	IgnoreFileExplorerError  func()

	newCurrentPath string
}

func (c *FileExplorer) Render() app.UI {
	return app.Div().
		Body(
			app.H2().
				Text("Files"),
			app.Div().
				Body(
					app.Div().Body(
						app.Code().
							Text(c.CurrentPath),
					),
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
					app.Ul().
						Body(
							app.Range(c.Index).Slice(func(i int) app.UI {
								return app.Li().
									Body(
										app.If(
											c.Index[i].IsDir(),
											app.Button().
												OnClick(func(ctx app.Context, e app.Event) {
													c.SetCurrentPath(filepath.Join(c.CurrentPath, c.Index[i].Name()))
												}).
												Text(c.Index[i].Name()+"/"),
										).Else(
											app.Text(c.Index[i].Name()),
										),
									)
							}),
						),
				),
		)
}
