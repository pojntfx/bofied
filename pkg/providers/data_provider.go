package providers

import (
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/pojntfx/bofied/pkg/constants"
	"github.com/pojntfx/bofied/pkg/servers"
	"github.com/pojntfx/bofied/pkg/validators"

	"github.com/studio-b12/gowebdav"
)

type DataProviderChildrenProps struct {
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
}

type DataProvider struct {
	app.Compo

	BackendURL   string
	IDToken      string
	WebDAVClient *gowebdav.Client
	Children     func(dpcp DataProviderChildrenProps) app.UI

	configFile    string
	configFileErr error

	currentPath     string
	index           []os.FileInfo
	httpShareLink   string
	tftpShareLink   string
	fileExplorerErr error

	operationCurrentPath string
	operationIndex       []os.FileInfo
}

func (c *DataProvider) Render() app.UI {
	address, username, password := c.getWebDAVCredentials()

	return c.Children(DataProviderChildrenProps{
		// Config file editor
		ConfigFile:    c.configFile,
		SetConfigFile: c.setConfigFile,

		FormatConfigFile:  c.formatConfigFile,
		RefreshConfigFile: c.refreshConfigFile,
		SaveConfigFile:    c.saveConfigFile,

		ConfigFileError:       c.configFileErr,
		IgnoreConfigFileError: c.ignoreConfigFileError,

		// File explorer
		CurrentPath: c.currentPath,
		SetCurrentPath: func(s string) {
			c.setCurrentPath(s, false)
		},

		Index:        c.index,
		RefreshIndex: c.refreshIndex,
		WriteToPath:  c.writeToPath,

		HTTPShareLink: c.httpShareLink,
		TFTPShareLink: c.tftpShareLink,
		SharePath:     c.sharePath,

		CreatePath: c.createPath,
		DeletePath: c.deletePath,
		MovePath:   c.movePath,
		CopyPath:   c.copyPath,

		WebDAVAddress:  address,
		WebDAVUsername: username,
		WebDAVPassword: password,

		OperationIndex: c.operationIndex,

		OperationCurrentPath: c.operationCurrentPath,
		OperationSetCurrentPath: func(s string) {
			c.setCurrentPath(s, true)
		},

		FileExplorerError:        c.fileExplorerErr,
		RecoverFileExplorerError: c.recoverFileExplorerError,
		IgnoreFileExplorerError:  c.ignoreFileExplorerError,
	})
}

func (c *DataProvider) OnMount(ctx app.Context) {
	c.refreshConfigFile()
	c.refreshIndex()
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
}

func (c *DataProvider) formatConfigFile() {
	formattedConfigFile, err := validators.FormatGoSrc(c.configFile)
	if err != nil {
		c.panicConfigFileError(err)

		return
	}

	c.configFile = formattedConfigFile
}

func (c *DataProvider) refreshConfigFile() {
	content, err := c.WebDAVClient.Read(constants.BootConfigFileName)
	if err != nil {
		c.panicConfigFileError(err)
	}

	c.setConfigFile(string(content))
}

func (c *DataProvider) saveConfigFile() {
	if err := validators.CheckGoSyntax(c.configFile); err != nil {
		c.panicConfigFileError(err)

		return
	}

	if err := c.WebDAVClient.Write(constants.BootConfigFileName, []byte(c.configFile), os.ModePerm); err != nil {
		c.panicConfigFileError(err)

		return
	}
}

func (c *DataProvider) ignoreConfigFileError() {
	// Only clear the error
	c.configFileErr = nil
}

func (c *DataProvider) panicConfigFileError(err error) {
	// Set the error
	c.configFileErr = err
}

// File explorer
func (c *DataProvider) setCurrentPath(path string, operationPath bool) {
	rawDirs, err := c.WebDAVClient.ReadDir(path)
	if err != nil {
		c.panicFileExplorerError(err)

		return
	}

	filteredDirs := []os.FileInfo{}
	for _, dir := range rawDirs {
		if dir.Name() != constants.BootConfigFileName {
			filteredDirs = append(filteredDirs, dir)
		}
	}

	if operationPath {
		c.operationCurrentPath = path
		c.operationIndex = filteredDirs
	} else {
		c.currentPath = path
		c.index = filteredDirs
	}
}

func (c *DataProvider) refreshIndex() {
	// On initial render, render root
	if c.currentPath == "" {
		c.currentPath = "/"
	}
	if c.operationCurrentPath == "" {
		c.operationCurrentPath = "/"
	}

	c.setCurrentPath(c.currentPath, false)
	c.setCurrentPath(c.operationCurrentPath, true)
}

func (c *DataProvider) writeToPath(path string, content []byte) {
	if err := c.WebDAVClient.Write(path, content, os.ModePerm); err != nil {
		c.panicFileExplorerError(err)

		return
	}
}

func (c *DataProvider) sharePath(path string) {
	// Parse URL
	u, err := url.Parse(c.BackendURL)
	if err != nil {
		c.panicFileExplorerError(err)
	}

	// Replace `private` prefix with `public` prefix
	pathParts := filepath.SplitList(u.Path)
	if len(pathParts) > 0 {
		pathParts = pathParts[1:]
	}
	u.Path = filepath.Join(filepath.Join(append([]string{servers.HTTPPrefix}, pathParts...)...), path)

	// Set HTTP share link
	c.httpShareLink = u.String()

	u.Scheme = "tftp"
	u.Host = u.Hostname() + ":" + constants.TFTPPort

	// Set TFTP share link
	c.tftpShareLink = u.String()
}

func (c *DataProvider) createPath(path string) {
	if err := c.WebDAVClient.MkdirAll(path, os.ModePerm); err != nil {
		c.panicFileExplorerError(err)

		return
	}

	c.refreshIndex()
}

func (c *DataProvider) deletePath(path string) {
	if err := c.WebDAVClient.RemoveAll(path); err != nil {
		c.panicFileExplorerError(err)

		return
	}

	c.refreshIndex()
}

func (c *DataProvider) movePath(src string, dst string) {
	if err := c.WebDAVClient.Rename(src, dst, true); err != nil {
		c.panicFileExplorerError(err)

		return
	}

	c.refreshIndex()
}

func (c *DataProvider) copyPath(src string, dst string) {
	if err := c.WebDAVClient.Copy(src, dst, true); err != nil {
		c.panicFileExplorerError(err)

		return
	}

	c.refreshIndex()
}

func (c *DataProvider) getWebDAVCredentials() (address string, username string, password string) {
	// Parse URL
	u, err := url.Parse(c.BackendURL)
	if err != nil {
		c.panicFileExplorerError(err)

		return "", "", ""
	}

	// Make it a WebDAV URL
	if u.Scheme == "https" {
		u.Scheme = "davs"
	} else {
		u.Scheme = "dav"
	}

	// Add current folder
	u.Path = path.Join(u.Path, c.currentPath)

	return u.String(), constants.OIDCOverBasicAuthUsername, c.IDToken
}

func (c *DataProvider) recoverFileExplorerError() {
	c.ignoreFileExplorerError()
}

func (c *DataProvider) ignoreFileExplorerError() {
	// Only clear the error
	c.fileExplorerErr = nil
}

func (c *DataProvider) panicFileExplorerError(err error) {
	// Set the error
	c.fileExplorerErr = err
}
