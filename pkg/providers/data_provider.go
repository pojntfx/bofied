package providers

import (
	"net/url"

	"github.com/maxence-charriere/go-app/v8/pkg/app"
	"github.com/pojntfx/bofied/pkg/constants"

	"github.com/studio-b12/gowebdav"
)

type DataProviderChildrenProps struct {
	GetAuthorizedWebDAVURL func() string

	Error   error
	Recover func()
	Ignore  func()
}

type DataProvider struct {
	app.Compo

	BackendURL   string
	IDToken      string
	WebDAVClient *gowebdav.Client
	Children     func(dpcp DataProviderChildrenProps) app.UI

	err error
}

func (c *DataProvider) Render() app.UI {
	return c.Children(DataProviderChildrenProps{
		GetAuthorizedWebDAVURL: func() string {
			// Parse URL
			u, err := url.Parse(c.BackendURL)
			if err != nil {
				c.panic(err)
			}

			// Make it a WebDAV URL
			if u.Scheme == "https" {
				u.Scheme = "davs"
			} else {
				u.Scheme = "dav"
			}

			// Add basic auth
			u.User = url.UserPassword(constants.OIDCOverBasicAuthUsername, c.IDToken)

			return u.String()
		},

		Error:   c.err,
		Recover: c.recover,
		Ignore:  c.ignore,
	})
}

func (c *DataProvider) recover() {
	// Recover ignore for now, as updating will re-evaluate potentially fault expressions
	c.ignore()
}

func (c *DataProvider) ignore() {
	// Only clear the error
	c.err = nil

	c.Update()
}

func (c *DataProvider) panic(err error) {
	// Set the error
	c.err = err

	c.Update()
}
