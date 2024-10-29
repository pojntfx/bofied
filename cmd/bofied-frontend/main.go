package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/kataras/compress"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/pojntfx/bofied/pkg/components"
)

func main() {
	// Client-side code
	{
		// Define the routes
		app.Route("/", func() app.Composer {
			return &components.Home{}
		})

		// Start the app
		app.RunWhenOnBrowser()
	}

	// Server-/build-side code
	{
		// Parse the flags
		build := flag.Bool("build", false, "Create static build")
		out := flag.String("out", "out/bofied-frontend", "Out directory for static build")
		path := flag.String("path", "", "Base path for static build")
		serve := flag.Bool("serve", false, "Build and serve the frontend")
		laddr := flag.String("laddr", "localhost:15255", "Address to serve the frontend on")

		flag.Parse()

		// Define the handler
		h := &app.Handler{
			Author:          "Felicitas Pojtinger",
			BackgroundColor: "#151515",
			Description:     "Modern network boot server.",
			Icon: app.Icon{
				Default: "/web/icon.png",
			},
			Keywords: []string{
				"pxe-boot",
				"ipxe",
				"netboot",
				"network-boot",
				"http-server",
				"dhcp-server",
				"pxe",
				"webdav-server",
				"tftp-server",
				"proxy-dhcp",
			},
			LoadingLabel: "Modern network boot server.",
			Name:         "bofied",
			RawHeaders: []string{
				`<meta property="og:url" content="https://pojntfx.github.io/bofied/">`,
				`<meta property="og:title" content="bofied">`,
				`<meta property="og:description" content="Modern network boot server.">`,
				`<meta property="og:image" content="https://pojntfx.github.io/bofied/web/icon.png">`,
			},
			Styles: []string{
				`https://unpkg.com/@patternfly/patternfly@4.102.2/patternfly.css`,
				`https://unpkg.com/@patternfly/patternfly@4.102.2/patternfly-addons.css`,
				`/web/index.css`,
			},
			ThemeColor: "#151515",
			Title:      "bofied",
		}

		// Create static build if specified
		if *build {
			// Deploy under a path
			if *path != "" {
				h.Resources = app.GitHubPages(*path)
			}

			if err := app.GenerateStaticWebsite(*out, h); err != nil {
				log.Fatalf("could not build: %v\n", err)
			}
		}

		// Serve if specified
		if *serve {
			log.Printf("bofied frontend listening on %v\n", *laddr)

			if err := http.ListenAndServe(*laddr, compress.Handler(h)); err != nil {
				log.Fatalf("could not open bofied frontend: %v\n", err)
			}
		}
	}
}
