# bofied

Network boot nodes in a network.

[![hydrun CI](https://github.com/pojntfx/bofied/actions/workflows/hydrun.yaml/badge.svg)](https://github.com/pojntfx/bofied/actions/workflows/hydrun.yaml)
[![Docker CI](https://github.com/pojntfx/bofied/actions/workflows/docker.yaml/badge.svg)](https://github.com/pojntfx/bofied/actions/workflows/docker.yaml)
[![Go Reference](https://pkg.go.dev/badge/github.com/pojntfx/bofied.svg)](https://pkg.go.dev/github.com/pojntfx/bofied)
[![Matrix](https://img.shields.io/matrix/bofied:matrix.org)](https://matrix.to/#/#bofied:matrix.org?via=matrix.org)

## Overview

bofied is a network boot server. It provides everything you need to PXE boot a node, from a (proxy)DHCP server for PXE service to a TFTP and HTTP server to serve boot files.

It enables you to ...

- **Boot nodes from the network**: Using (proxy)DHCP for PXE service, it can configure nodes which are set to network boot
- **Serve boot files**: The integrated TFTP and HTTP servers can provide the iPXE network bootloader, Linux distros or other boot files
- **Easily manage and script the network boot config**: By using the browser or WebDAV, boot files can be managed and the scriptable configuration can be edited
- **Monitor network boot**: By monitoring (proxy)DHCP and TFTP traffic, bofied can give an insight to the network boot process using a browser or the gRPC API
- **Remotely provision nodes**: Because bofied is based on open web technologies and supports OpenID Connect authentication, it can be securely exposed to the public internet and be used to manage network boot in a offsite location

## Installation

### Containerized

You can get the Docker container like so:

```shell
$ docker pull pojntfx/bofied-backend
```

### Natively

If you prefer a native installation, static binaries are also available on [GitHub releases](https://github.com/pojntfx/bofied/releases).

You can install them like so:

```shell
$ curl -L -o /tmp/bofied-backend https://github.com/pojntfx/bofied/releases/download/latest/bofied-backend.linux-$(uname -m)
$ sudo install /tmp/bofied-backend /usr/local/bin
$ sudo setcap cap_net_bind_service+ep /usr/local/bin/bofied-backend # This allows rootless execution
```

### About the Frontend

The frontend is also available on [GitHub releases](https://github.com/pojntfx/bofied/releases) in the form of a static `.tar.gz` archive; to deploy it, simply upload it to a CDN or copy it to a web server. For most users, this shouldn't be necessary though; thanks to [@maxence-charriere](https://github.com/maxence-charriere)'s [go-app package](https://go-app.dev/), bofied is a progressive web app. By simply visiting the [public deployment](https://pojntfx.github.io/bofied/) once, it will be available for offline use whenever you need it.

## Usage

### Setting up Authentication

bofied uses [OpenID Connect](https://en.wikipedia.org/wiki/OpenID_Connect) for authentication, which means you can use almost any authentication provider, both self-hosted and as a service, that you want to. We've created a short tutorial video which shows how to set up [Auth0](https://auth0.com/) for this purpose, but feel free to use something like [Ory](https://github.com/ory/hydra) if you prefer a self-hosted solution:

[<img src="https://img.youtube.com/vi/N3cocCOsrGw/0.jpg" width="512" alt="Setting up OpenID Connect for Internal Apps YouTube Video" title="Setting up OpenID Connect for Internal Apps YouTube Video">](https://www.youtube.com/watch?v=N3cocCOsrGw)

### Verifying Port Availability

First, verify that ports `67/udp`, `4011/udp`, `69/udp`, `15256/tcp` and `15257/tcp` aren't in use by another app:

```bash
$ ss -tlnp | grep -E -- ':(15256|15257)'
$ ss -ulnp | grep -E -- ':(67|4011|69)'
```

Neither of these two commands should return anything; if they do, kill the process that listens on the port.

### Getting the Advertised IP

bofied integrates a (proxy)DHCP server, which advertises the IP address of the integrated TFTP server. To do so, you'll have to find out the IP of the node which is running bofied; you can find it with `ip a`:

```bash
$ ip -4 a
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
2: enp0s13f0u1u3: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc fq_codel state UP group default qlen 1000
    inet 192.168.178.147/24 brd 192.168.178.255 scope global dynamic noprefixroute enp0s13f0u1u3
       valid_lft 862274sec preferred_lft 862274sec
```

In the following, we'll assume that `192.168.178.147` is the IP address of this node.

### Starting the Backend (Containerized)

Using Docker (or an alternative like Podman), you can now easily start & configure the backend; see the [Reference](#reference) for more configuration parameters:

```shell
$ docker run \
    --name bofied-backend \
    -d \
    --restart always \
    --cap-add NET_BIND_SERVICE \
    -p 67:67/udp \
    -p 4011:4011/udp \
    -p 69:69/udp \
    -p 15256:15256/tcp \
    -p 15257:15257/tcp \
    -v ${HOME}/.local/share/bofied:/root/.local/share/bofied:z \
    -e BOFIED_BACKEND_OIDCISSUER=https://pojntfx.eu.auth0.com/ \
    -e BOFIED_BACKEND_OIDCCLIENTID=myoidcclientid \
    -e BOFIED_BACKEND_ADVERTISEDIP=192.168.178.147 \
    pojntfx/bofied-backend
```

The logs are available like so:

```shell
$ docker logs bofied-backend
```

### Starting the Backend (Natively)

If you prefer a native setup, a non-containerized installation is also possible.

First, set up a config file at `~/.local/share/bofied/etc/bofied/bofied-backend-config.yaml`; see the [Reference](#reference) for more configuration parameters:

```shell
$ mkdir -p ~/.local/share/bofied/etc/bofied/
$ cat <<EOT >~/.local/share/bofied/etc/bofied/bofied-backend-config.yaml
oidcIssuer: https://pojntfx.eu.auth0.com/
oidcClientID: myoidcclientid
advertisedIP: 192.168.178.147
EOT
```

Now, create a systemd service for it:

```shell
$ mkdir -p ~/.config/systemd/user/
$ cat <<EOT >~/.config/systemd/user/bofied-backend.service
[Unit]
Description=bofied

[Service]
ExecStart=/usr/local/bin/bofied-backend -c \${HOME}/.local/share/bofied/etc/bofied/bofied-backend-config.yaml

[Install]
WantedBy=multi-user.target
EOT
```

Finally, reload systemd and enable the service:

```shell
$ systemctl --user daemon-reload
$ systemctl --user enable --now bofied-backend
```

You can get the logs like so:

```shell
$ journalctl --user -u bofied-backend
```

### Setting up the Firewall

You might also have to open up the ports on your firewall:

```shell
$ for port in 67/udp 4011/udp 69/udp 15256/tcp 15257/tcp; do sudo firewall-cmd --permanent --add-port=${port}; done
```

### Connecting the Frontend

Now that the backend is running, head over to [https://pojntfx.github.io/bofied/](https://pojntfx.github.io/bofied/). Alternatively, as described in [About the Frontend](#about-the-frontend), you can also choose to self-host. Once you're on the page, you should be presented with the following setup page:

![Setup page](./assets/setup.png)

You'll have to enter your own information here; the `Backend URL` is the URL on which the backend (`http://localhost:15256/` by default) runs, the `OIDC Issuer`, `Client ID` and `Redirect URL` are the same values that you've set the backend up with above.

Finally, click on `Login`, and if everything worked out fine you should be presented with the initial launch screen:

![Initial page](./assets/initial.png)

🚀 **That's it**! We hope you enjoy using bofied.

## Reference

🚧 This section is a work-in-progress! It will be added soon. 🚧

## License

bofied (c) 2021 Felicitas Pojtinger and contributors

SPDX-License-Identifier: AGPL-3.0
