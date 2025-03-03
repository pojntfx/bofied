<img alt="Project icon" style="vertical-align: middle;" src="./docs/icon.svg" width="128" height="128" align="left">

# bofied

Modern network boot server.

<br/>

https://user-images.githubusercontent.com/28832235/117546130-42745500-b029-11eb-804d-134ff5049ccc.mp4

[![hydrun CI](https://github.com/pojntfx/bofied/actions/workflows/hydrun.yaml/badge.svg)](https://github.com/pojntfx/bofied/actions/workflows/hydrun.yaml)
[![Docker CI](https://github.com/pojntfx/bofied/actions/workflows/docker.yaml/badge.svg)](https://github.com/pojntfx/bofied/actions/workflows/docker.yaml)
![Go Version](https://img.shields.io/badge/go%20version-%3E=1.22-61CFDD.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/pojntfx/bofied.svg)](https://pkg.go.dev/github.com/pojntfx/bofied)
[![Matrix](https://img.shields.io/matrix/bofied:matrix.org)](https://matrix.to/#/#bofied:matrix.org?via=matrix.org)

## Overview

bofied is a network boot server. It provides everything you need to PXE boot a node, from a (proxy)DHCP server for PXE service to a TFTP and HTTP server to serve boot files.

It enables you to:

- **Boot nodes from the network**: Using (proxy)DHCP for PXE service, it can configure nodes which are set to network boot
- **Serve boot files**: The integrated TFTP and HTTP servers can provide the iPXE network bootloader, Linux distros or other boot files
- **Easily manage and script the network boot config**: By using the browser or WebDAV, boot files can be managed and the scriptable configuration can be edited
- **Monitor network boot**: By monitoring it's (proxy)DHCP and TFTP traffic, bofied can give an insight to the network boot process using a browser or the gRPC API
- **Remotely provision nodes**: Because bofied is based on open web technologies and supports OpenID Connect authentication, it can be securely exposed to the public internet and be used to manage network boot in a offsite location

## Installation

### Containerized

You can get the OCI image like so:

```console
$ podman pull ghcr.io/pojntfx/bofied-backend
```

### Natively

Static binaries are available on [GitHub releases](https://github.com/pojntfx/bofied/releases).

On Linux, you can install them like so:

```console
$ curl -L -o /tmp/bofied-backend "https://github.com/pojntfx/bofied/releases/latest/download/bofied-backend.linux-$(uname -m)"
$ sudo install /tmp/bofied-backend /usr/local/bin
$ sudo setcap cap_net_admin+ep /usr/local/bin/bofied-backend # This allows rootless execution
```

On macOS, you can use the following:

```console
$ curl -L -o /tmp/bofied-backend "https://github.com/pojntfx/bofied/releases/latest/download/bofied-backend.darwin-$(uname -m)"
$ sudo install /tmp/bofied-backend /usr/local/bin
```

On Windows, the following should work (using PowerShell as administrator):

```console
PS> Invoke-WebRequest https://github.com/pojntfx/bofied-backend/releases/latest/download/bofied-backend.windows-x86_64.exe -OutFile \Windows\System32\bofied-backend.exe
```

You can find binaries for more operating systems and architectures on [GitHub releases](https://github.com/pojntfx/bofied/releases).

### Frontend

The frontend is publicly accessible and can be opened in any browser. It's a [PWA](https://en.wikipedia.org/wiki/Progressive_web_app), meaning it will continue to work offline after your first visit:

<a href="https://pojntfx.github.io/bofied/"><img src="https://github.com/pojntfx/webnetesctl/raw/main/img/launch.png" alt="PWA badge" width="200"/></a>

To self-host, see [contributing](#contributing). A `.tar.gz` archive of the frontend, ready for deployment on any static web server, is available on [GitHub releases](https://github.com/pojntfx/bofied/releases).

## Tutorial

> TL;DR: Setup OIDC authentication, configure the backend, start the backend, open up a port on the firewall, and connect using the frontend

### 1. Setting up Authentication

bofied uses [OpenID Connect](https://en.wikipedia.org/wiki/OpenID_Connect) for authentication, which means you can use almost any authentication provider, both self-hosted and as a service, that you want to. We've created a short tutorial video which shows how to set up [Auth0](https://auth0.com/) for this purpose, but feel free to use something like [Ory](https://github.com/ory/hydra) if you prefer a self-hosted solution:

[<img src="https://img.youtube.com/vi/N3cocCOsrGw/0.jpg" width="256" alt="Setting up OpenID Connect for Internal Apps YouTube Video" title="Setting up OpenID Connect for Internal Apps YouTube Video">](https://www.youtube.com/watch?v=N3cocCOsrGw)

### 2. Verifying Port Availability

First, verify that ports `67/udp`, `4011/udp`, `69/udp`, `15256/tcp` and `15257/tcp` aren't in use by another app:

```bash
$ ss -tlnp | grep -E -- ':(15256|15257)'
$ ss -ulnp | grep -E -- ':(67|4011|69)'
```

Neither of these two commands should return anything; if they do, kill the process that listens on the port.

### 3. Getting the Advertised IP

bofied integrates a (proxy)DHCP server, which advertises the IP address of the integrated TFTP server. To do so, you'll have to find out the IP of the node which is running bofied; you can find it with `ip a`:

```bash
$ ip -4 a
# ...
2: enp0s13f0u1u3: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc fq_codel state UP group default qlen 1000
    inet 192.168.178.147/24 brd 192.168.178.255 scope global dynamic noprefixroute enp0s13f0u1u3
       valid_lft 862274sec preferred_lft 862274sec
# ...
```

In the following, we'll assume that `192.168.178.147` is the IP address of this node.

### 4. Starting the Backend

Next up is starting the backend. It is available both in containerized and static binary form; see the [Reference](#reference) for more configuration parameters:

<details>
  <summary>Expand containerized instructions</summary>

Run the following:

```console
$ sudo podman run \
    --name bofied-backend \
    -d \
    --restart always \
    --net host \
    --cap-add NET_BIND_SERVICE \
    -v ${HOME}/.local/share/bofied:/root/.local/share/bofied:z \
    -e BOFIED_BACKEND_OIDCISSUER=https://pojntfx.eu.auth0.com/ \
    -e BOFIED_BACKEND_OIDCCLIENTID=myoidcclientid \
    -e BOFIED_BACKEND_ADVERTISEDIP=192.168.178.147 \
    pojntfx/bofied-backend
$ sudo podman generate systemd --new bofied-backend | sudo tee /lib/systemd/system/bofied-backend.service

$ sudo systemctl daemon-reload

$ sudo systemctl enable --now bofied-backend
```

The logs are available like so:

```console
$ journalctl -u bofied-backend
```

</details>

<details>
  <summary>Expand native instructions</summary>

First, set up a config file at `~/.local/share/bofied/etc/bofied/bofied-backend-config.yaml`; see the [Reference](#reference) for more configuration parameters:

```console
$ sudo su
# mkdir -p ~/.local/share/bofied/etc/bofied/
# cat <<EOT >~/.local/share/bofied/etc/bofied/bofied-backend-config.yaml
oidcIssuer: https://pojntfx.eu.auth0.com/
oidcClientID: myoidcclientid
advertisedIP: 192.168.178.147
EOT
```

Now, create a systemd service for it:

```console
$ sudo tee /etc/systemd/system/bofied-backend.service<<'EOT'
[Unit]
Description=bofied Backend

[Service]
ExecStart=/usr/local/bin/bofied-backend -c \${HOME}/.local/share/bofied/etc/bofied/bofied-backend-config.yaml

[Install]
WantedBy=multi-user.target
EOT
```

Finally, reload systemd and enable the service:

```console
$ sudo systemctl daemon-reload
$ sudo systemctl enable --now bofied-backend
```

You can get the logs like so:

```console
$ journalctl -u bofied-backend
```

</details>

### 5. Setting up the Firewall

You might also have to open up the ports on your firewall:

```console
$ for port in 67/udp 4011/udp 69/udp 15256/tcp 15257/tcp; do sudo firewall-cmd --permanent --add-port=${port}; done
```

### 6. Connecting the Frontend

Now that the backend is running, head over to [https://pojntfx.github.io/bofied/](https://pojntfx.github.io/bofied/):

[<img src="https://github.com/alphahorizonio/webnetesctl/raw/main/img/launch.png" width="240">](https://pojntfx.github.io/bofied/)

Alternatively, as described in [Frontend](#frontend), you can also choose to self-host. Once you're on the page, you should be presented with the following setup page:

![Setup page](./docs/screenshot-setup.png)

You'll have to enter your own information here; the `Backend URL` is the URL on which the backend (`http://localhost:15256/` by default) runs, the `OIDC Issuer`, `Client ID` and `Redirect URL` are the same values that you've set the backend up with above.

Finally, click on `Login`, and if everything worked out fine you should be presented with the initial launch screen:

![Initial page](./docs/screenshot-initial.png)

🚀 **That's it**! We hope you enjoy using bofied.

## Screenshots

Click on an image to see a larger version.

<a display="inline" href="./docs/screenshot-syntax-validation.png?raw=true">
<img src="./docs/screenshot-syntax-validation.png" width="45%" alt="Screenshot of syntax validation" title="Screenshot of syntax validation">
</a>

<a display="inline" href="./docs/screenshot-monitoring.png?raw=true">
<img src="./docs/screenshot-monitoring.png" width="45%" alt="Screenshot of monitoring" title="Screenshot of monitoring">
</a>

<a display="inline" href="./docs/screenshot-file-operations-2.png?raw=true">
<img src="./docs/screenshot-file-operations-2.png" width="45%" alt="Screenshot 2 of file operations" title="Screenshot 2 of file operations">
</a>

<a display="inline" href="./docs/screenshot-file-operations-3.png?raw=true">
<img src="./docs/screenshot-file-operations-3.png" width="45%" alt="Screenshot 3 of file operations" title="Screenshot 3 of file operations">
</a>

<a display="inline" href="./docs/screenshot-sharing.png?raw=true">
<img src="./docs/screenshot-sharing.png" width="45%" alt="Screenshot of sharing" title="Screenshot of sharing">
</a>

<a display="inline" href="./docs/screenshot-text-editor.png?raw=true">
<img src="./docs/screenshot-text-editor.png" width="45%" alt="Screenshot of the text editor" title="Screenshot of the text editor">
</a>

<a display="inline" href="./docs/screenshot-about-modal.png?raw=true">
<img src="./docs/screenshot-about-modal.png" width="45%" alt="Screenshot of about modal" title="Screenshot of about modal">
</a>

## Reference

### Command Line Arguments

```console
$ bofied-backend --help
bofied is a network boot server. It provides everything you need to PXE boot a node, from a (proxy)DHCP server for PXE service to a TFTP and HTTP server to serve boot files.

For more information, please visit https://github.com/pojntfx/bofied.

Usage:
  bofied-backend [flags]

Flags:
      --advertisedIP string                IP to advertise for DHCP clients (default "100.64.154.246")
  -c, --configFile string                  Config file to use
      --dhcpListenAddress string           Listen address for DHCP server (default ":67")
      --extendedHTTPListenAddress string   Listen address for WebDAV, HTTP and gRPC-Web server (default ":15256")
      --grpcListenAddress string           Listen address for gRPC server (default ":15257")
  -h, --help                               help for bofied-backend
  -t, --oidcClientID string                OIDC client ID (default "myoidcclientid")
  -i, --oidcIssuer string                  OIDC issuer (default "https://pojntfx.eu.auth0.com/")
      --proxyDHCPListenAddress string      Listen address for proxyDHCP server (default ":4011")
  -p, --pureConfig Configuration           Prevent usage of stdlib in configuration file, even if enabled in Configuration function
  -s, --skipStarterDownload                Don't initialize by downloading the starter on the first run
      --starterURL string                  Download URL to a starter .tar.gz archive; the default chainloads https://netboot.xyz/ (default "https://github.com/pojntfx/ipxe-binaries/releases/latest/download/ipxe.tar.gz")
      --tftpListenAddress string           Listen address for TFTP server (default ":69")
  -d, --workingDir string                  Working directory (default "/home/pojntfx/.local/share/bofied/var/lib/bofied")
```

### Environment Variables

All command line arguments described above can also be set using environment variables; for example, to set `--advertisedIP` to `192.168.178.147` with an environment variable, use `BOFIED_BACKEND_ADVERTISEDIP=192.168.178.147`.

### Configuration File

Just like with the environment variables, bofied can also be configured using a configuration file; see [examples/bofied-backend-config.yaml](./examples/bofied-backend-config.yaml) for an example configuration file.

### Config Script

The config script is separate from the config file and is used to dynamically decide which file to send to which node based on it's IP address, MAC address and processor architecture. It can be set & validated using either the frontend or WebDAV. The default config script, which is fetched from [pojntfx/ipxe-binaries](https://github.com/pojntfx/ipxe-binaries), returns a matching executable based on the architecture:

```go
package config

func Filename(
	ip string,
	macAddress string,
	arch string,
	archID int,
) string {
	switch arch {
	case "x86 BIOS":
		return "ipxe-i386.kpxe"
	case "x86 UEFI":
		return "ipxe-i386.efi"
	case "x64 UEFI":
		return "ipxe-x86_64.efi"
	case "ARM 32-bit UEFI":
		return "ipxe-arm32.efi"
	case "ARM 64-bit UEFI":
		return "ipxe-arm64.efi"
	default:
		return "ipxe-i386.kpxe"
	}
}

func Configure() map[string]string {
	return map[string]string{
		"useStdlib": "false",
	}
}
```

The script is just a small [Go](https://go.dev/) program which exports two functions: **`Filename`** and **`Configure`**. **`Configure`** is called to configure the interpreter; for example, if you want to use the standard library, i.e. to log information with `log.Println` or to make a HTTP request with `http.Get`, you can set `"useStdlib": "true",`. **`Filename`** is called with the IP address, MAC address and architecture (as a string and as an ID), and should return the name of the file to send to the booting node. The following architecture values are available (see [IANA Processor Architecture Types](https://www.iana.org/assignments/dhcpv6-parameters/dhcpv6-parameters.xhtml#processor-architecture)):

<details>
  <summary>Expand available architecture values</summary>

| `archID` Parameter | `arch` Parameter                     |
| ------------------ | ------------------------------------ |
| 0x00               | x86 BIOS                             |
| 0x01               | NEC/PC98 (DEPRECATED)                |
| 0x02               | Itanium                              |
| 0x03               | DEC Alpha (DEPRECATED)               |
| 0x04               | Arc x86 (DEPRECATED)                 |
| 0x05               | Intel Lean Client (DEPRECATED)       |
| 0x06               | x86 UEFI                             |
| 0x07               | x64 UEFI                             |
| 0x08               | EFI Xscale (DEPRECATED)              |
| 0x09               | EBC                                  |
| 0x0a               | ARM 32-bit UEFI                      |
| 0x0b               | ARM 64-bit UEFI                      |
| 0x0c               | PowerPC Open Firmware                |
| 0x0d               | PowerPC ePAPR                        |
| 0x0e               | POWER OPAL v3                        |
| 0x0f               | x86 uefi boot from http              |
| 0x10               | x64 uefi boot from http              |
| 0x11               | ebc boot from http                   |
| 0x12               | arm uefi 32 boot from http           |
| 0x13               | arm uefi 64 boot from http           |
| 0x14               | pc/at bios boot from http            |
| 0x15               | arm 32 uboot                         |
| 0x16               | arm 64 uboot                         |
| 0x17               | arm uboot 32 boot from http          |
| 0x18               | arm uboot 64 boot from http          |
| 0x19               | RISC-V 32-bit UEFI                   |
| 0x1a               | RISC-V 32-bit UEFI boot from http    |
| 0x1b               | RISC-V 64-bit UEFI                   |
| 0x1c               | RISC-V 64-bit UEFI boot from http    |
| 0x1d               | RISC-V 128-bit UEFI                  |
| 0x1e               | RISC-V 128-bit UEFI boot from http   |
| 0x1f               | s390 Basic                           |
| 0x20               | s390 Extended                        |
| 0x21               | MIPS 32-bit UEFI                     |
| 0x22               | MIPS 64-bit UEFI                     |
| 0x23               | Sunway 32-bit UEFI                   |
| 0x24               | Sunway 64-bit UEFI                   |
| 0x25               | LoongArch 32-bit UEFI                |
| 0x26               | LoongArch 32-bit UEFI boot from http |
| 0x27               | LoongArch 64-bit UEFI                |
| 0x28               | LoongArch 64-bit UEFI boot from http |
| 0x29               | arm rpiboot                          |
| 0x2a-0xff          | Unassigned                           |

</details>

When bofied is first started, it automatically downloads [pojntfx/ipxe-binaries](https://github.com/pojntfx/ipxe-binaries) to the boot file directory, so without configuring anything you can already network boot many Linux distros and other operating systems thanks to [netboot.xyz](https://netboot.xyz/). This behavior can of course also be disabled, in which case only a minimal config file will be created; see [Reference](#reference).

### WebDAV

In addition to using the frontend to manage boot files, you can also mount them using [WebDAV](https://en.wikipedia.org/wiki/WebDAV). You can the required credentials by using the `Mount directory` button in the frontend:

![Mount directory modal](./docs/screenshot-mount-directory.png)

Using a file manager like [Files](https://en.wikipedia.org/wiki/GNOME_Files), you can now mount the folder:

![GNOME Files WebDAV mounting](./docs/screenshot-gnome-files-webdav-mounting.png)

When transfering large files, using WebDAV directly is the recommended method.

![GNOME Files WebDAV listing](./docs/screenshot-gnome-files-webdav-listing.png)

### gRPC API

bofied exposes a streaming gRPC and gRPC-Web API for monitoring network boot, which is also in use internally in the frontend. You can find the relevant `.proto` files in [api/proto/v1](./api/proto/v1); send the OpenID Connect token with the `X-Bofied-Authorization` metadata key.

## Acknowledgements

- This project would not have been possible were it not for [@maxence-charriere](https://github.com/maxence-charriere)'s [go-app package](https://go-app.dev/); if you enjoy using bofied, please donate to him!
- The open source [PatternFly design system](https://www.patternfly.org/v4/) provides a professional design and reduced the need for custom CSS to a minimium (less than 30 SLOC!).
- [pin/tftp](https://github.com/pin/tftp) provides the TFTP functionality for bofied.
- [studio-b12/gowebdav](https://github.com/studio-b12/gowebdav) provides the WebDAV client for the bofied frontend.
- The [yaegi Go interpreter](https://github.com/traefik/yaegi) is used to securely evaluate the config script.

## Contributing

To contribute, please use the [GitHub flow](https://guides.github.com/introduction/flow/) and follow our [Code of Conduct](./CODE_OF_CONDUCT.md).

To build and start a development version of bofied locally, run the following:

```console
$ git clone https://github.com/pojntfx/bofied.git
$ cd bofied
$ make depend
$ BOFIED_BACKEND_OIDCISSUER=https://pojntfx.eu.auth0.com/ BOFIED_BACKEND_OIDCCLIENTID=myoidcclientid BOFIED_BACKEND_ADVERTISEDIP=192.168.178.147 make -j dev
```

The backend should now be started and the frontend be available on [http://localhost:15225/](http://localhost:15225/).

Have any questions or need help? Chat with us [on Matrix](https://matrix.to/#/#bofied:matrix.org?via=matrix.org)!

## Related Projects

If you want to have persistent inventory of services and nodes on your network or turn the nodes in it on remotely, check out [liwasc](https://github.com/pojntfx/liwasc)!

## Troubleshooting

- If you run bofied over Wifi or advertise your Wifi adapter's IP, old PXE clients might hang at PXE errors such as `TFTP cannot read from connection` or `TFTP read timeout` due to very slow transfer speeds or other problems with complex network topologies. To fix this, disconnect other network adapters - i.e. if you're running bofied on a laptop with both a Wifi and an ethernet card and you advertise the Wifi card's IP, disconnect the ethernet cable. bofied's log and [Wireshark](https://www.wireshark.org/) might also give you more insights in such situations. For the best reliablity, run bofied on a wired connection.

## License

bofied (c) 2024 Felicitas Pojtinger and contributors

SPDX-License-Identifier: AGPL-3.0-or-later
