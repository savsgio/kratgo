Kratgo
======

[![Build Status](https://travis-ci.org/savsgio/kratgo.svg?branch=master)](https://travis-ci.org/savsgio/kratgo)
[![Go Report Card](https://goreportcard.com/badge/github.com/savsgio/kratgo)](https://goreportcard.com/report/github.com/savsgio/kratgo)
<!-- [![Coverage Status](https://coveralls.io/repos/github/savsgio/kratgo/badge.svg?branch=master)](https://coveralls.io/github/savsgio/kratgo?branch=master) -->
<!-- [![GoDoc](https://godoc.org/github.com/savsgio/kratgo?status.svg)](https://godoc.org/github.com/savsgio/kratgo) -->
<!-- [![GitHub release](https://img.shields.io/github/release/savsgio/kratgo.svg)](https://github.com/savsgio/kratgo/releases) -->

Simple, lightweight and ultra-fast HTTP Cache written in Go for accelerate your websites speed.

## Requirements

- [Go](https://golang.org/dl/) >= 1.11.X
- make
- git

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

## Developers

Copy configuration file `./config/kratgo.conf.yml` to `./config/kratgo-dev.conf.yml`, and customize it.

Run with:

```bash
make run
```

Contributing
------------

**Feel free to contribute it or fork me...** :wink:
