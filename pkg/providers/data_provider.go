package providers

import (
	"net/url"
	"os"
	"path/filepath"

	"github.com/maxence-charriere/go-app/v8/pkg/app"
	"github.com/pojntfx/bofied/pkg/constants"
	"github.com/pojntfx/bofied/pkg/validators"

	"github.com/studio-b12/gowebdav"
)

type DataProviderChildrenProps struct {
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

type DataProvider struct {
	app.Compo

	BackendURL   string
	IDToken      string
	WebDAVClient *gowebdav.Client
	Children     func(dpcp DataProviderChildrenProps) app.UI

	configFile string
	index      []os.FileInfo
	currentDir string

	err           error
	configFileErr error
}

func (c *DataProvider) Render() app.UI {
	return c.Children(DataProviderChildrenProps{
		AuthorizedWebDAVURL: c.getAuthorizedWebDAVURL(),
		Index:               c.index,
		CurrentDir:          c.currentDir,

		SetCurrentDir: c.setCurrentDir,
		UploadFile:    c.uploadFile,
		Refresh:       c.refresh,

		Error:   c.err,
		Recover: c.recover,
		Ignore:  c.ignore,

		// Config file editor
		ConfigFile:    c.configFile,
		SetConfigFile: c.setConfigFile,

		FormatConfigFile:  c.formatConfigFile,
		RefreshConfigFile: c.refreshConfigFile,
		SaveConfigFile:    c.saveConfigFile,

		ConfigFileError:        c.configFileErr,
		RecoverConfigFileError: c.recoverConfigFileError,
		IgnoreConfigFileError:  c.ignoreConfigFileError,
	})
}

func (c *DataProvider) OnMount(ctx app.Context) {
	c.refresh(ctx)
}

func (c *DataProvider) refresh(ctx app.Context) {
	c.configFile = c.getConfigFile()

	// On initial render, render the working directory
	if c.currentDir == "" {
		c.currentDir = "."
	}
	c.setCurrentDir(c.currentDir, ctx)

	c.Update()
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

func (c *DataProvider) setCurrentDir(dir string, ctx app.Context) {
	rawDirs, err := c.WebDAVClient.ReadDir(dir)
	if err != nil {
		c.panic(err)

		return
	}

	filteredDirs := []os.FileInfo{}
	for _, dir := range rawDirs {
		if dir.Name() != constants.BootConfigFileName {
			filteredDirs = append(filteredDirs, dir)
		}
	}

	ctx.Dispatch(func() {
		c.currentDir = dir
		c.index = filteredDirs

		c.Update()
	})
}

func (c *DataProvider) validateConfigFile() {
	if err := validators.CheckGoSyntax(c.configFile); err != nil {
		c.panic(err)

		return
	}
}

func (c *DataProvider) uploadFile(name string, content []byte) {
	if err := c.WebDAVClient.Write(filepath.Join(c.currentDir, name), content, os.ModePerm); err != nil {
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

// Config file editor
func (c *DataProvider) setConfigFile(s string) {
	// Clear the error (if it's still faulty it will be set again below)
	c.configFileErr = nil

	// Set the new config file
	c.configFile = s

	// Check the syntax and show errors if they exist
	if err := validators.CheckGoSyntax(c.configFile); err != nil {
		c.panicConfigFileError(err)
	}

	c.Update()
}

func (c *DataProvider) formatConfigFile() {
	formattedConfigFile, err := validators.FormatGoSrc(c.configFile)
	if err != nil {
		c.panicConfigFileError(err)

		return
	}

	c.configFile = formattedConfigFile

	c.Update()
}

func (c *DataProvider) refreshConfigFile() {
	c.configFile = c.getConfigFile()

	c.Update()
}

func (c *DataProvider) saveConfigFile() {
	if err := c.WebDAVClient.Write(constants.BootConfigFileName, []byte(c.configFile), os.ModePerm); err != nil {
		c.panicConfigFileError(err)

		return
	}
}

func (c *DataProvider) recoverConfigFileError() {
	c.ignoreConfigFileError()
}

func (c *DataProvider) ignoreConfigFileError() {
	// Only clear the error
	c.configFileErr = nil

	c.Update()
}

func (c *DataProvider) panicConfigFileError(err error) {
	// Set the error
	c.configFileErr = err

	c.Update()
}
