package components

import (
	"os"

	"github.com/maxence-charriere/go-app/v8/pkg/app"
)

type DataShell struct {
	app.Compo

	// Config file editor
	ConfigFile    string
	SetConfigFile func(string)

	FormatConfigFile  func()
	RefreshConfigFile func()
	SaveConfigFile    func()

	ConfigFileError        error
	RecoverConfigFileError func()
	IgnoreConfigFileError  func()

	// File explorer
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

	AuthorizedWebDAVURL string

	FileExplorerError        error
	RecoverFileExplorerError func()
	IgnoreFileExplorerError  func()
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
		app.Section().
			Body(
				&FileExplorer{
					CurrentPath:    c.CurrentPath,
					SetCurrentPath: c.SetCurrentPath,

					Index:        c.Index,
					RefreshIndex: c.RefreshIndex,
					UploadFile:   c.UploadFile,

					ShareLink: c.ShareLink,
					SharePath: c.SharePath,

					CreateDirectory: c.CreateDirectory,
					DeletePath:      c.DeletePath,
					MovePath:        c.MovePath,
					CopyPath:        c.CopyPath,

					AuthorizedWebDAVURL: c.AuthorizedWebDAVURL,

					Error:   c.FileExplorerError,
					Recover: c.RecoverFileExplorerError,
					Ignore:  c.IgnoreFileExplorerError,
				},
			),
	)
}
