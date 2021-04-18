package providers

import (
	"github.com/maxence-charriere/go-app/v8/pkg/app"

	"github.com/studio-b12/gowebdav"
)

type DataProviderChildrenProps struct {
	GetIDToken func() string
}

type DataProvider struct {
	app.Compo

	IDToken      string
	WebDAVClient *gowebdav.Client
	Children     func(dpcp DataProviderChildrenProps) app.UI
}

func (c *DataProvider) Render() app.UI {
	return c.Children(DataProviderChildrenProps{
		GetIDToken: func() string {
			return c.IDToken
		},
	})
}
