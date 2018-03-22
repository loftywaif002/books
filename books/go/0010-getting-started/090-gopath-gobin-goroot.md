---
Title: GOPATH, GOROOT, GOBIN
Id: 10
SOId: 14406
---

## `GOPATH`

It's important to understand the effect of the `GOPATH` environment variable.

If you come from other programming languages, you're probably used to placing source code anywhere in the file system.

Go tools expect a certain layout of the source code.

`GOPATH` is the root of the workspace and contains the following folders:

* `src` — location of source files: `.go`, `.c`, `.g`, `.s`
* `pkg` — location of compiled packages (`.a` files)
* `bin` — location of executables built by Go

Like the system `PATH` environment variable, Go path is a `:` (`;` on Windows) delimited list of directories where Go will look for packages. The `go get` tool will also download packages to the first directory in this list.

Since Go 1.8, the `GOPATH` environment variable will have a default value if it is unset. It defaults to `$HOME/go` on Unix/Linux and `%USERPROFILE%/go` on Windows.

Some tools assume that `GOPATH` only consists of a single directory.

## `GOBIN`

The bin directory where `go install` and `go get` will place binaries after building `main` packages. Generally this is set to somewhere on the system `PATH` so that installed binaries can be run and discovered easily.

## `GOROOT`

This is the location of your Go installation. It is used to find the standard libraries. It is very rare to have to set this variable as Go embeds the build path into the toolchain. Setting `GOROOT` is needed if the installation directory differs from the build directory (or the value set when building).

See [go env](328) for a full list of environment variables.
