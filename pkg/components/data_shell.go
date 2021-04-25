package components

import (
	"os"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/pojntfx/bofied/pkg/providers"
)

type DataShell struct {
	app.Compo

	// Config file editor
	ConfigFile    string
	SetConfigFile func(string)

	FormatConfigFile  func()
	RefreshConfigFile func()
	SaveConfigFile    func()

	ConfigFileError       error
	IgnoreConfigFileError func()

	// File explorer
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

	FileExplorerError        error
	RecoverFileExplorerError func()
	IgnoreFileExplorerError  func()

	Events []providers.Event

	EventsError        error
	RecoverEventsError func(app.Context)
	IgnoreEventsError  func()
}

func (c *DataShell) Render() app.UI {
	// Gather notifications
	notifications := []Notification{}
	for _, event := range c.Events {
		notifications = append(notifications, Notification{
			CreatedAt: event.CreatedAt.String(),
			Message:   event.Message,
		})
	}

	return app.Div().Body(
		app.Section().
			Body(
				&ConfigFileEditor{
					ConfigFile:    c.ConfigFile,
					SetConfigFile: c.SetConfigFile,

					FormatConfigFile:  c.FormatConfigFile,
					RefreshConfigFile: c.RefreshConfigFile,
					SaveConfigFile:    c.SaveConfigFile,

					Error:  c.ConfigFileError,
					Ignore: c.IgnoreConfigFileError,
				},
			),
		app.Section().
			Body(
				&FileExplorer{
					CurrentPath:    c.CurrentPath,
					SetCurrentPath: c.SetCurrentPath,

					Index:        c.Index,
					RefreshIndex: c.RefreshIndex,
					WriteToPath:  c.WriteToPath,

					HTTPShareLink: c.HTTPShareLink,
					TFTPShareLink: c.TFTPShareLink,
					SharePath:     c.SharePath,

					CreatePath: c.CreatePath,
					DeletePath: c.DeletePath,
					MovePath:   c.MovePath,
					CopyPath:   c.CopyPath,

					WebDAVAddress:  c.WebDAVAddress,
					WebDAVUsername: c.WebDAVUsername,
					WebDAVPassword: c.WebDAVPassword,

					OperationIndex: c.OperationIndex,

					OperationCurrentPath:    c.OperationCurrentPath,
					OperationSetCurrentPath: c.OperationSetCurrentPath,

					Error:   c.FileExplorerError,
					Recover: c.RecoverFileExplorerError,
					Ignore:  c.IgnoreFileExplorerError,
				},
			),
		app.Section().Body(
			&NotificationDrawer{
				Notifications: notifications,
			},
		),
	)
}
