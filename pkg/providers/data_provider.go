package providers

import (
	"net/url"
	"os"

	"github.com/maxence-charriere/go-app/v8/pkg/app"
	"github.com/pojntfx/bofied/pkg/constants"
	"github.com/pojntfx/bofied/pkg/validators"

	"github.com/studio-b12/gowebdav"
)

type DataProviderChildrenProps struct {
	AuthorizedWebDAVURL string
	ConfigFile          string

	SetConfigFile      func(string)
	ValidateConfigFile func()
	SaveConfigFile     func()

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

	configFile string

	err error
}

func (c *DataProvider) Render() app.UI {
	return c.Children(DataProviderChildrenProps{
		AuthorizedWebDAVURL: c.getAuthorizedWebDAVURL(),
		ConfigFile:          c.configFile,

		SetConfigFile:      c.setConfigFile,
		ValidateConfigFile: c.validateConfigFile,
		SaveConfigFile:     c.saveConfigFile,

		Error:   c.err,
		Recover: c.recover,
		Ignore:  c.ignore,
	})
}

func (c *DataProvider) OnMount(ctx app.Context) {
	ctx.Dispatch(func() {
		c.configFile = c.getConfigFile()

		c.Update()
	})
}

func (c *DataProvider) getAuthorizedWebDAVURL() string {
	// Parse URL
	u, err := url.Parse(c.BackendURL)
	if err != nil {
		c.panic(err)

		return ""
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
}

func (c *DataProvider) getConfigFile() string {
	content, err := c.WebDAVClient.Read(constants.BootConfigFileName)
	if err != nil {
		c.panic(err)

		return ""
	}

	return string(content)
}

func (c *DataProvider) setConfigFile(s string) {
	c.configFile = s

	c.Update()
}

func (c *DataProvider) validateConfigFile() {
	if err := validators.CheckGoSyntax(c.configFile); err != nil {
		c.panic(err)

		return
	}
}

func (c *DataProvider) saveConfigFile() {
	if err := validators.CheckGoSyntax(c.configFile); err != nil {
		c.panic(err)

		return
	}

	if err := c.WebDAVClient.Write(constants.BootConfigFileName, []byte(c.configFile), os.ModePerm); err != nil {
		c.panic(err)

		return
	}
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
