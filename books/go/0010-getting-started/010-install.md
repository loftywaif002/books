---
Title: Install Go toolchain
Id: 2
SOId: 20381
---

There are 2 production quality compilers:

* gc, the official compiler
* [gccgo](https://golang.org/doc/install/gccgo). Also created by the Go project, but not used as often as gc.

This article describes installing gc compiler.

## Install on Windows

Run `.msi` installer downloaded from [https://golang.org/dl/](https://golang.org/dl/).

Read [setting up Windows environment](3-windows-setup).

## Install on Mac OS

### Using official binaries

Download `.pkg` installer from [https://golang.org/dl/](https://golang.org/dl/) and run it.

### Using Homebrew

* if you don't have Homebrew installed, install it following the [instructions](https://brew.sh/)
* `brew install go`

Read [setting up Mac environment](4-mac-os-setup).

## Install on Ubuntu

### Using Ubuntu provided package

```sh
$ sudo apt-get update
$ sudo apt-get install go
```

Note that packages provided by Ubuntu are often outdated. A new version of Go is released every 6 months but Ubuntu distribution moves at a slower pace.

For that reason we recommend installing the binary packages.

### Using binary packages

These instructions work on pretty much every Linux distribution:

```sh
$ sudo apt-get update
$ sudo apt-get install -y build-essential git curl wget
$ wget https://storage.googleapis.com/golang/go<versions>.gz
```

You can find the version list [here](https://golang.org/doc/install).

```sh
# To install go1.9.3 use
$ wget https://storage.googleapis.com/golang/go1.9.3.linux-amd64.tar.gz

# Untar the file
$ sudo tar -C /usr/local -xzf go1.9.3.linux-amd64.tar.gz
$ sudo chown -R $USER:$USER /usr/local/go
$ rm go1.9.3.linux-amd64.tar.gz
```

Read [setting up Unix environment](5).

## Install on other OSes

Follow [official instructions](https://golang.org/doc/install).
