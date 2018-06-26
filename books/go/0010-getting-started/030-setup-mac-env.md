---
Title: Mac OS setup
Id: 4
SOId: rd600058
---

After installing the compiler you need to configure [`GOPATH`](10) environment variable.

Since Go 1.8, the `GOPATH` environment variable has the default value of `$HOME/go`, so you can skip setting it.

Create the go directory with `mkdir $HOME/go`.

Add the following to your `~/.bash_profile` file:

```sh
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```

Load the changes with `source ~/.bash_profile` or launch a new terminal for the changes to take effect.

**Explanation**

File `~/.bash_profile` is executed at startup by the default `bash` shell.

By adding commands there we make them available inside every shell session.

Adding `$GOPATH/bin` to `PATH` is a matter of convenience. When you install Go programs with `go get ...`, you will be able to invoke them without typing the full path. For example, you would run `gotest1` rather than `$HOME/go/bin/gotest1`.

**More configuration**

To quickly cd to a directory with Go source code I add bash alias to `~/.bash_profile`:

```sh
alias cdgo="cd $GOPATH/src/github.com/kjk"
```

You need to change `github.com/kjk` for your GitHub account.

You should understand the [effect of GOPATH](10).
