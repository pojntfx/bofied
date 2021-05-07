# bofied

Network boot nodes in a network.

[![hydrun CI](https://github.com/pojntfx/bofied/actions/workflows/hydrun.yaml/badge.svg)](https://github.com/pojntfx/bofied/actions/workflows/hydrun.yaml)
[![Docker CI](https://github.com/pojntfx/bofied/actions/workflows/docker.yaml/badge.svg)](https://github.com/pojntfx/bofied/actions/workflows/docker.yaml)
[![Go Reference](https://pkg.go.dev/badge/github.com/pojntfx/bofied.svg)](https://pkg.go.dev/github.com/pojntfx/bofied)
[![Matrix](https://img.shields.io/matrix/bofied:matrix.org)](https://matrix.to/#/#bofied:matrix.org?via=matrix.org)

## Overview

🚧 This project is a work-in-progress! Instructions will be added as soon as it is usable. 🚧

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

## License

bofied (c) 2021 Felicitas Pojtinger and contributors

SPDX-License-Identifier: AGPL-3.0
