package providers

import (
	"context"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
	api "github.com/pojntfx/bofied/pkg/api/proto/v1"
	"github.com/pojntfx/bofied/pkg/authorization"
	"github.com/pojntfx/bofied/pkg/constants"
	"github.com/pojntfx/bofied/pkg/servers"
	"github.com/pojntfx/bofied/pkg/services"
	"github.com/pojntfx/bofied/pkg/validators"
	"github.com/pojntfx/go-app-grpc-chat-frontend-web/pkg/websocketproxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/studio-b12/gowebdav"
)

type Event struct {
	CreatedAt time.Time
	Message   string
}

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

	HTTPShareLink url.URL
	TFTPShareLink url.URL
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

	Events []Event

	EventsError        error
	RecoverEventsError func(app.Context)
	IgnoreEventsError  func()

	// Metadata
	UseAdvertisedIP    bool
	SetUseAdvertisedIP func(bool)

	SetUseSecureProtocol func(bool)
}

type DataProvider struct {
	app.Compo

	BackendURL string
	IDToken    string
	Children   func(dpcp DataProviderChildrenProps) app.UI

	webDAVClient         *gowebdav.Client
	authenticatedContext context.Context
	eventsService        api.EventsServiceClient
	metadataService      api.MetadataServiceClient

	configFile    string
	configFileErr error

	currentPath     string
	index           []os.FileInfo
	httpShareLink   url.URL
	tftpShareLink   url.URL
	fileExplorerErr error

	operationCurrentPath string
	operationIndex       []os.FileInfo

	events    []Event
	eventsErr error

	advertisedIP string

	useAdvertisedIP   bool
	useSecureProtocol bool
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

		HTTPShareLink: func() url.URL {
			u := c.httpShareLink

			if c.useAdvertisedIP {
				u.Host = c.advertisedIP
			}

			if c.useSecureProtocol {
				u.Scheme = "https"
			}

			return u
		}(),
		TFTPShareLink: func() url.URL {
			u := c.tftpShareLink

			if c.useAdvertisedIP {
				u.Host = c.advertisedIP
			}

			return u
		}(),
		SharePath: c.sharePath,

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

		Events: c.events,

		EventsError:        c.eventsErr,
		RecoverEventsError: c.recoverEventsError,
		IgnoreEventsError:  c.ignoreEventsError,

		// Metadata
		UseAdvertisedIP:    c.useAdvertisedIP,
		SetUseAdvertisedIP: c.setUseAdvertisedIP,

		SetUseSecureProtocol: c.setUseSecureProtocol,
	})
}

func (c *DataProvider) OnMount(ctx app.Context) {
	// Initialize events
	c.events = []Event{}

	// Create WebDAV client
	webDAVClient := gowebdav.NewClient(path.Join(c.BackendURL, servers.WebDAVPrefix), constants.OIDCOverBasicAuthUsername, c.IDToken)
	header, value := authorization.GetOIDCOverBasicAuthHeader(constants.OIDCOverBasicAuthUsername, c.IDToken)
	webDAVClient.SetHeader(header, value)
	c.webDAVClient = webDAVClient

	// Parse URL for gRPC client
	u, err := url.Parse(c.BackendURL)
	if err != nil {
		c.panicEventsError(err)

		return
	}

	// Make it a WebSocket URL
	if u.Scheme == "https" {
		u.Scheme = "wss"
	} else {
		u.Scheme = "ws"
	}

	// Create gRPC client
	conn, err := grpc.Dial(path.Join(u.String(), servers.GRPCPrefix), grpc.WithContextDialer(websocketproxy.NewWebSocketProxyClient(time.Minute).Dialer), grpc.WithInsecure())
	if err != nil {
		c.panicEventsError(err)

		return
	}
	c.eventsService = api.NewEventsServiceClient(conn)
	c.metadataService = api.NewMetadataServiceClient(conn)
	c.authenticatedContext = metadata.AppendToOutgoingContext(context.Background(), services.AuthorizationMetadataKey, c.IDToken)

	// Refresh/subscribe to data
	c.refreshConfigFile()
	c.refreshIndex()
	go c.subscribeToEvents(ctx)
	go c.getMetadata(ctx)
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
	content, err := c.webDAVClient.Read(constants.BootConfigFileName)
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

	if err := c.webDAVClient.Write(constants.BootConfigFileName, []byte(c.configFile), os.ModePerm); err != nil {
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
	rawDirs, err := c.webDAVClient.ReadDir(path)
	if err != nil {
		c.panicFileExplorerError(err)

		return
	}

	filteredDirs := []os.FileInfo{}
	for _, dir := range rawDirs {
		if dir.Name() != constants.BootConfigFileName {
			// Hide directories for operations
			if operationPath && !dir.IsDir() {
				continue
			}

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
	if err := c.webDAVClient.Write(path, content, os.ModePerm); err != nil {
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

	// Sync with the share options
	if u.Host == c.advertisedIP {
		c.useAdvertisedIP = true
	}

	// Replace `private` prefix with `public` prefix
	u.Path = filepath.Join(filepath.Join(append([]string{servers.HTTPPrefix}, filepath.SplitList(u.Path)...)...), path)

	// Set HTTP share link
	c.httpShareLink = *u

	u.Scheme = "tftp"
	u.Host = u.Hostname() + ":" + constants.TFTPPort

	// Set TFTP share link
	c.tftpShareLink = *u
}

func (c *DataProvider) createPath(path string) {
	if err := c.webDAVClient.MkdirAll(path, os.ModePerm); err != nil {
		c.panicFileExplorerError(err)

		return
	}

	c.refreshIndex()
}

func (c *DataProvider) deletePath(path string) {
	if err := c.webDAVClient.RemoveAll(path); err != nil {
		c.panicFileExplorerError(err)

		return
	}

	c.refreshIndex()
}

func (c *DataProvider) movePath(src string, dst string) {
	if err := c.webDAVClient.Rename(src, dst, true); err != nil {
		c.panicFileExplorerError(err)

		return
	}

	c.refreshIndex()
}

func (c *DataProvider) copyPath(src string, dst string) {
	if err := c.webDAVClient.Copy(src, dst, true); err != nil {
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

	// Add the prefix
	u.Path = path.Join(u.Path, servers.WebDAVPrefix)

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

func (c *DataProvider) subscribeToEvents(ctx app.Context) {
	// Get stream from service
	events, err := c.eventsService.SubscribeToEvents(c.authenticatedContext, &emptypb.Empty{})
	if err != nil {
		// We have to use `Context.Emit` here as this runs from a separate Goroutine
		ctx.Emit(func() {
			c.panicEventsError(err)
		})

		return
	}

	// Process stream
	for {
		// Receive event from stream
		event, err := events.Recv()
		if err != nil {
			// We have to use `Context.Emit` here as this runs from a separate Goroutine
			ctx.Emit(func() {
				c.panicEventsError(err)
			})

			return
		}

		// Parse the event's date
		eventCreatedAt, err := time.Parse(time.RFC3339, event.GetCreatedAt())
		if err != nil {
			// We have to use `Context.Emit` here as this runs from a separate Goroutine
			ctx.Emit(func() {
				c.panicEventsError(err)
			})

			return
		}

		// Add the event (we have to use `Context.Emit` here as this runs from a separate Goroutine)
		ctx.Emit(func() {
			c.events = append(c.events, Event{
				CreatedAt: eventCreatedAt,
				Message:   event.GetMessage(),
			})
		})
	}
}

func (c *DataProvider) recoverEventsError(ctx app.Context) {
	// Clear the error
	c.eventsErr = nil

	// Resubscribe
	c.OnMount(ctx)
}

func (c *DataProvider) ignoreEventsError() {
	// Only clear the error
	c.eventsErr = nil
}

func (c *DataProvider) panicEventsError(err error) {
	// Set the error
	c.eventsErr = err
}

// Metadata
func (c *DataProvider) getMetadata(ctx app.Context) {
	metadata, err := c.metadataService.GetMetadata(c.authenticatedContext, &emptypb.Empty{})
	if err != nil {
		// We have to use `Context.Emit` here as this runs from a separate Goroutine
		ctx.Emit(func() {
			c.panicEventsError(err)
		})

		return
	}

	// We have to use `Context.Emit` here as this runs from a separate Goroutine
	ctx.Emit(func() {
		c.advertisedIP = metadata.GetAdvertisedIP()
	})
}

func (c *DataProvider) setUseAdvertisedIP(b bool) {
	c.useAdvertisedIP = b
}

func (c *DataProvider) setUseSecureProtocol(b bool) {
	c.useSecureProtocol = b
}
