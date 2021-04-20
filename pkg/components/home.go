package components

import (
	"github.com/maxence-charriere/go-app/v8/pkg/app"
	"github.com/pojntfx/bofied/pkg/authorization"
	"github.com/pojntfx/bofied/pkg/constants"
	"github.com/pojntfx/bofied/pkg/providers"
	"github.com/pojntfx/liwasc/pkg/components"
	metaproviders "github.com/pojntfx/liwasc/pkg/providers"
	"github.com/studio-b12/gowebdav"
)

type Home struct {
	app.Compo
}

func (c *Home) Render() app.UI {
	return &metaproviders.ConfigurationProvider{
		StoragePrefix:       "bofied.configuration",
		StateQueryParameter: "state",
		CodeQueryParameter:  "code",
		Children: func(cpcp metaproviders.SetupProviderChildrenProps) app.UI {
			// This div is required so that there are no authorization loops
			return app.Div().
				TabIndex(-1).
				Class("pf-x-ws-router").
				Body(
					app.If(cpcp.Ready,
						// Identity provider
						&metaproviders.IdentityProvider{
							Issuer:        cpcp.OIDCIssuer,
							ClientID:      cpcp.OIDCClientID,
							RedirectURL:   cpcp.OIDCRedirectURL,
							HomeURL:       "/",
							Scopes:        []string{"profile", "email"},
							StoragePrefix: "bofied.identity",
							Children: func(ipcp metaproviders.IdentityProviderChildrenProps) app.UI {
								// Configuration shell
								if ipcp.Error != nil {
									return &components.SetupShell{
										LogoSrc:          "/web/logo.svg",
										Title:            "Log in to bofied",
										ShortDescription: "Network boot nodes in a network.",
										LongDescription:  `bofied is a network boot server. It provides everything you need to PXE boot a node, from a (proxy)DHCP server for PXE service to a TFTP and HTTP server to serve boot files.`,
										HelpLink:         "https://github.com/pojntfx/liwasc#Usage",
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

								// Authorized WebDAV Client
								webDAVClient := gowebdav.NewClient(cpcp.BackendURL, constants.OIDCOverBasicAuthUsername, ipcp.IDToken)
								header, value := authorization.GetOIDCOverBasicAuthHeader(constants.OIDCOverBasicAuthUsername, ipcp.IDToken)
								webDAVClient.SetHeader(header, value)

								// Data provider
								return &providers.DataProvider{
									BackendURL:   cpcp.BackendURL,
									IDToken:      ipcp.IDToken,
									WebDAVClient: webDAVClient,
									Children: func(dpcp providers.DataProviderChildrenProps) app.UI {
										// Data shell
										return &DataShell{
											// Config file editor
											ConfigFile:    dpcp.ConfigFile,
											SetConfigFile: dpcp.SetConfigFile,

											FormatConfigFile:  dpcp.FormatConfigFile,
											RefreshConfigFile: dpcp.RefreshConfigFile,
											SaveConfigFile:    dpcp.SaveConfigFile,

											ConfigFileError:        dpcp.ConfigFileError,
											RecoverConfigFileError: dpcp.RecoverConfigFileError,
											IgnoreConfigFileError:  dpcp.IgnoreConfigFileError,

											// File explorer
											CurrentPath:    dpcp.CurrentPath,
											SetCurrentPath: dpcp.SetCurrentPath,

											Index:        dpcp.Index,
											RefreshIndex: dpcp.RefreshIndex,
											UploadFile:   dpcp.UploadFile,

											ShareLink: dpcp.ShareLink,
											SharePath: dpcp.SharePath,

											CreateDirectory: dpcp.CreateDirectory,
											DeletePath:      dpcp.DeletePath,
											MovePath:        dpcp.MovePath,
											CopyPath:        dpcp.CopyPath,

											AuthorizedWebDAVURL: dpcp.AuthorizedWebDAVURL,

											FileExplorerError:        dpcp.FileExplorerError,
											RecoverFileExplorerError: dpcp.RecoverFileExplorerError,
											IgnoreFileExplorerError:  dpcp.IgnoreFileExplorerError,
										}
									},
								}
							},
						},
					).Else(
						// Configuration shell
						&components.SetupShell{
							LogoSrc:          "/web/logo.svg",
							Title:            "Log in to bofied",
							ShortDescription: "Network boot nodes in a network.",
							LongDescription:  `bofied is a network boot server. It provides everything you need to PXE boot a node, from a (proxy)DHCP server for PXE service to a TFTP and HTTP server to serve boot files.`,
							HelpLink:         "https://github.com/pojntfx/liwasc#Usage",
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
						},
					),
				)
		},
	}
}
