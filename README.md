Kratgo
======

[![Build Status](https://travis-ci.org/savsgio/kratgo.svg?branch=master)](https://travis-ci.org/savsgio/kratgo)
[![Go Report Card](https://goreportcard.com/badge/github.com/savsgio/kratgo)](https://goreportcard.com/report/github.com/savsgio/kratgo)
[![GitHub release](https://img.shields.io/github/release/savsgio/kratgo.svg)](https://github.com/savsgio/kratgo/releases)
[![Docker](https://img.shields.io/docker/automated/savsgio/kratgo.svg?colorB=blue&style=flat)](https://hub.docker.com/r/savsgio/kratgo)
<!-- [![GoDoc](https://godoc.org/github.com/savsgio/kratgo?status.svg)](https://godoc.org/github.com/savsgio/kratgo) -->

Simple, lightweight and ultra-fast HTTP Cache to speed up your websites.


### Requirements

- [Go](https://golang.org/dl/) >= 1.12.X
- make
- git


## Features:

- Cache proxy.
- Load balancing beetwen backends.
- Cache invalidation via API (Admin).
- Configuration to non-cache certain requests.
- Configuration to set or unset headers on especific requests.

## General

To known if request pass across Kratgo Cache in backend servers, check the request header `X-Kratgo-Cache` with value `true`.


## Install

Clone the repository:

```bash
git clone https://github.com/savsgio/kratgo.git && cd kratgo
```

and execute:

```bash
make
make install
```

The binary file will install in `/usr/local/bin/kratgo` and configuration file in `/etc/kratgo/kratgo.conf.yml`


## Cache invalidation (Admin)

The cache invalidation is available via API. The API's address is configured in ***admin*** section of the configuration file.

This API only accepts ***POST*** requests with ***json***, under the path `/invalidate/`.

Ex: `http://localhost:6082/invalidate/`

The complete json body must be as following example:

```json
{
	"host": "www.example.com",
	"path": "/es/",
	"header": {
		"key": "Content-Type",
		"value": "text/plain; charset=utf-8"
	}
}
```

**IMPORTANT: All fields are optional, but at least you must specify one.**

All invalidations will process by workers in Kratgo. You can configure the maximum available workers in the configuration.

The workers are activated only when necessary.


## Docker

The docker image is available in Docker Hub: [savsgio/kratgo](https://hub.docker.com/r/savsgio/kratgo)

Get a basic configuration from [here](https://github.com/savsgio/kratgo/blob/master/config/kratgo.conf.yml) and customize it.

Run with:

```bash
docker run --rm --name kratgo -it -v <VOLUME WITH CONFIG> -p 6081:6081 -p 6082:6082 savsgio/kratgo -config <CONFIG FILE PATH IN THE VOLUME>
```

## Developers

Copy configuration file `./config/kratgo.conf.yml` to `./config/kratgo-dev.conf.yml`, and customize it.

Run with:

```bash
make run
```

Contributing
------------

**Feel free to contribute it or fork me...** :wink:
