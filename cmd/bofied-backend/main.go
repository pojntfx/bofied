package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/pojntfx/bofied/pkg/config"
	"github.com/pojntfx/bofied/pkg/constants"
	"github.com/pojntfx/bofied/pkg/eventing"
	"github.com/pojntfx/bofied/pkg/servers"
	"github.com/pojntfx/bofied/pkg/services"
	"github.com/pojntfx/liwasc/pkg/validators"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	configFileKey                 = "configFile"
	workingDirKey                 = "workingDir"
	advertisedIPKey               = "advertisedIP"
	dhcpListenAddressKey          = "dhcpListenAddress"
	proxyDHCPListenAddressKey     = "proxyDHCPListenAddress"
	tftpListenAddressKey          = "tftpListenAddress"
	webDAVAndHTTPListenAddressKey = "extendedHTTPListenAddress"
	oidcIssuerKey                 = "oidcIssuer"
	oidcClientIDKey               = "oidcClientID"
	eventsListenAddressKey        = "eventsListenAddress"
)

func main() {
	// Create command
	cmd := &cobra.Command{
		Use:   "bofied-backend",
		Short: "Network boot nodes in a network.",
		Long: `bofied is a network boot server. It provides everything you need to PXE boot a node, from a (proxy)DHCP server for PXE service to a TFTP and HTTP server to serve boot files.

For more information, please visit https://github.com/pojntfx/bofied.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Bind config file
			if !(viper.GetString(configFileKey) == "") {
				viper.SetConfigFile(viper.GetString(configFileKey))

				if err := viper.ReadInConfig(); err != nil {
					return err
				}
			}

			// Initialize the working directory
			if err := config.CreateConfigIfNotExists(filepath.Join(viper.GetString(workingDirKey), constants.BootConfigFileName)); err != nil {
				log.Fatal(err)
			}

			// Create eventing utilities
			eventsHandler := eventing.NewEventHandler()

			// Create auth utilities
			oidcValidator := validators.NewOIDCValidator(viper.GetString(oidcIssuerKey), viper.GetString(oidcClientIDKey))
			if err := oidcValidator.Open(); err != nil {
				log.Fatal(err)
			}
			contextValidator := validators.NewContextValidator(services.AuthorizationMetadataKey, oidcValidator)

			// Create services
			eventsService := services.NewEventsService(eventsHandler, contextValidator)

			// Create servers
			dhcpServer := servers.NewDHCPServer(
				viper.GetString(dhcpListenAddressKey),
				viper.GetString(advertisedIPKey),
				eventsHandler,
			)
			proxyDHCPServer := servers.NewProxyDHCPServer(
				viper.GetString(proxyDHCPListenAddressKey),
				viper.GetString(advertisedIPKey),
				filepath.Join(viper.GetString(workingDirKey), constants.BootConfigFileName),
				eventsHandler,
			)
			tftpServer := servers.NewTFTPServer(
				viper.GetString(workingDirKey),
				viper.GetString(tftpListenAddressKey),
				eventsHandler,
			)
			eventsServer, eventsServerHandler := servers.NewEventsServer(viper.GetString(eventsListenAddressKey), eventsService)
			extendedHTTPServer := servers.NewExtendedHTTPServer(viper.GetString(workingDirKey), viper.GetString(webDAVAndHTTPListenAddressKey), oidcValidator, eventsServerHandler)

			// Start servers
			log.Printf(
				"bofied backend listening on %v (DHCP), %v (proxyDHCP), %v (TFTP), %v (WebDAV on %v, HTTP on %v and gRPC-Web on %v) and %v (gRPC)\n",
				viper.GetString(dhcpListenAddressKey),
				viper.GetString(proxyDHCPListenAddressKey),
				viper.GetString(tftpListenAddressKey),
				viper.GetString(webDAVAndHTTPListenAddressKey),
				servers.WebDAVPrefix,
				servers.HTTPPrefix,
				servers.EventsPrefix,
				viper.GetString(eventsListenAddressKey),
			)

			go func() {
				log.Fatal(dhcpServer.ListenAndServe())
			}()

			go func() {
				log.Fatal(proxyDHCPServer.ListenAndServe())
			}()

			go func() {
				log.Fatal(tftpServer.ListenAndServe())
			}()

			go func() {
				log.Fatal(eventsServer.ListenAndServe())
			}()

			return extendedHTTPServer.ListenAndServe()
		},
	}

	// Get default working dir
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("could not get home directory", err)
	}
	workingDirDefault := filepath.Join(home, ".local", "share", "bofied", "var", "lib", "bofied")

	// Bind flags
	cmd.PersistentFlags().StringP(configFileKey, "c", "", "Config file to use")
	cmd.PersistentFlags().StringP(workingDirKey, "d", workingDirDefault, "Working directory")
	cmd.PersistentFlags().String(advertisedIPKey, "100.64.154.246", "IP to advertise for DHCP clients")

	cmd.PersistentFlags().String(dhcpListenAddressKey, ":67", "Listen address for DHCP server")
	cmd.PersistentFlags().String(proxyDHCPListenAddressKey, ":4011", "Listen address for proxyDHCP server")
	cmd.PersistentFlags().String(tftpListenAddressKey, ":"+constants.TFTPPort, "Listen address for TFTP server")
	cmd.PersistentFlags().String(webDAVAndHTTPListenAddressKey, ":15256", "Listen address for WebDAV, HTTP and gRPC-Web server")
	cmd.PersistentFlags().String(eventsListenAddressKey, ":15257", "Listen address for events gRPC server")

	cmd.PersistentFlags().StringP(oidcIssuerKey, "i", "https://pojntfx.eu.auth0.com/", "OIDC issuer")
	cmd.PersistentFlags().StringP(oidcClientIDKey, "t", "myoidcclientid", "OIDC client ID")

	// Bind env variables
	if err := viper.BindPFlags(cmd.PersistentFlags()); err != nil {
		log.Fatal(err)
	}
	viper.SetEnvPrefix("bofied_backend")
	viper.AutomaticEnv()

	// Run command
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
