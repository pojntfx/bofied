package components

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/pojntfx/bofied/pkg/providers"
)

type Home struct {
	app.Compo
}

func (c *Home) Render() app.UI {
	return &providers.SetupProvider{
		StoragePrefix:       "bofied.configuration",
		StateQueryParameter: "state",
		CodeQueryParameter:  "code",
		Children: func(cpcp providers.SetupProviderChildrenProps) app.UI {
			// This div is required so that there are no authorization loops
			return app.Div().
				Class("pf-x-ws-router").
				Body(
					app.If(cpcp.Ready,
						func() app.UI {
							// Identity provider
							return &providers.IdentityProvider{
								Issuer:        cpcp.OIDCIssuer,
								ClientID:      cpcp.OIDCClientID,
								RedirectURL:   cpcp.OIDCRedirectURL,
								HomeURL:       "/",
								Scopes:        []string{"profile", "email"},
								StoragePrefix: "bofied.identity",
								Children: func(ipcp providers.IdentityProviderChildrenProps) app.UI {
									// Configuration shell
									if ipcp.Error != nil {
										return &SetupShell{
											LogoSrc:          "/web/logo.svg",
											Title:            "Log in to bofied",
											ShortDescription: "Modern network boot server.",
											LongDescription:  `bofied is a network boot server. It provides everything you need to PXE boot a node, from a (proxy)DHCP server for PXE service to a TFTP and HTTP server to serve boot files.`,
											HelpLink:         "https://github.com/pojntfx/bofied#Usage",
											Links: map[string]string{
												"License":       "https://github.com/pojntfx/bofied/blob/main/LICENSE",
												"Source Code":   "https://github.com/pojntfx/bofied",
												"Documentation": "https://github.com/pojntfx/bofied#Usage",
											},

											BackendURL:      cpcp.BackendURL,
											OIDCIssuer:      cpcp.OIDCIssuer,
											OIDCClientID:    cpcp.OIDCClientID,
											OIDCRedirectURL: cpcp.OIDCRedirectURL,

											SetBackendURL:      cpcp.SetBackendURL,
											SetOIDCIssuer:      cpcp.SetOIDCIssuer,
											SetOIDCClientID:    cpcp.SetOIDCClientID,
											SetOIDCRedirectURL: cpcp.SetOIDCRedirectURL,
											ApplyConfig:        cpcp.ApplyConfig,

											Error: ipcp.Error,
										}
									}

									// Configuration placeholder
									if ipcp.IDToken == "" || ipcp.UserInfo.Email == "" {
										return app.P().Text("Authorizing ...")
									}

									// Data provider
									return &providers.DataProvider{
										BackendURL: cpcp.BackendURL,
										IDToken:    ipcp.IDToken,
										Children: func(dpcp providers.DataProviderChildrenProps) app.UI {
											// Data shell
											return &DataShell{
												// Config file editor
												ConfigFile:    dpcp.ConfigFile,
												SetConfigFile: dpcp.SetConfigFile,

												FormatConfigFile:  dpcp.FormatConfigFile,
												RefreshConfigFile: dpcp.RefreshConfigFile,
												SaveConfigFile:    dpcp.SaveConfigFile,

												ConfigFileError:       dpcp.ConfigFileError,
												IgnoreConfigFileError: dpcp.IgnoreConfigFileError,

												// File explorer
												CurrentPath:    dpcp.CurrentPath,
												SetCurrentPath: dpcp.SetCurrentPath,

												Index:        dpcp.Index,
												RefreshIndex: dpcp.RefreshIndex,
												WriteToPath:  dpcp.WriteToPath,

												HTTPShareLink: dpcp.HTTPShareLink,
												TFTPShareLink: dpcp.TFTPShareLink,
												SharePath:     dpcp.SharePath,

												CreatePath:      dpcp.CreatePath,
												CreateEmptyFile: dpcp.CreateEmptyFile,
												DeletePath:      dpcp.DeletePath,
												MovePath:        dpcp.MovePath,
												CopyPath:        dpcp.CopyPath,

												EditPathContents:    dpcp.EditPathContents,
												SetEditPathContents: dpcp.SetEditPathContents,
												EditPath:            dpcp.EditPath,

												WebDAVAddress:  dpcp.WebDAVAddress,
												WebDAVUsername: dpcp.WebDAVUsername,
												WebDAVPassword: dpcp.WebDAVPassword,

												OperationIndex: dpcp.OperationIndex,

												OperationCurrentPath:    dpcp.OperationCurrentPath,
												OperationSetCurrentPath: dpcp.OperationSetCurrentPath,

												FileExplorerError:        dpcp.FileExplorerError,
												RecoverFileExplorerError: dpcp.RecoverFileExplorerError,
												IgnoreFileExplorerError:  dpcp.IgnoreFileExplorerError,

												Events: dpcp.Events,

												EventsError:        dpcp.EventsError,
												RecoverEventsError: dpcp.RecoverEventsError,
												IgnoreEventsError:  dpcp.IgnoreEventsError,

												UserInfo: ipcp.UserInfo,
												Logout:   ipcp.Logout,

												// Metadata
												UseAdvertisedIP:    dpcp.UseAdvertisedIP,
												SetUseAdvertisedIP: dpcp.SetUseAdvertisedIP,

												UseAdvertisedIPForWebDAV:    dpcp.UseAdvertisedIPForWebDAV,
												SetUseAdvertisedIPForWebDAV: dpcp.SetUseAdvertisedIPForWebDAV,

												SetUseHTTPS: dpcp.SetUseHTTPS,
												SetUseDavs:  dpcp.SetUseDavs,
											}
										},
									}
								},
							}
						},
					).Else(
						func() app.UI {
							// Configuration shell
							return &SetupShell{
								LogoSrc:          "/web/logo.svg",
								Title:            "Log in to bofied",
								ShortDescription: "Modern network boot server.",
								LongDescription:  `bofied is a network boot server. It provides everything you need to PXE boot a node, from a (proxy)DHCP server for PXE service to a TFTP and HTTP server to serve boot files.`,
								HelpLink:         "https://github.com/pojntfx/bofied#Usage",
								Links: map[string]string{
									"License":       "https://github.com/pojntfx/bofied/blob/main/LICENSE",
									"Source Code":   "https://github.com/pojntfx/bofied",
									"Documentation": "https://github.com/pojntfx/bofied#Usage",
								},

								BackendURL:      cpcp.BackendURL,
								OIDCIssuer:      cpcp.OIDCIssuer,
								OIDCClientID:    cpcp.OIDCClientID,
								OIDCRedirectURL: cpcp.OIDCRedirectURL,

								SetBackendURL:      cpcp.SetBackendURL,
								SetOIDCIssuer:      cpcp.SetOIDCIssuer,
								SetOIDCClientID:    cpcp.SetOIDCClientID,
								SetOIDCRedirectURL: cpcp.SetOIDCRedirectURL,
								ApplyConfig:        cpcp.ApplyConfig,

								Error: cpcp.Error,
							}
						},
					),
				)
		},
	}
}
