---
Title: Parse command line arguments
Search: flags
Id: 171
SOId: 4023
---

Package [`flag`](https://golang.org/pkg/flag/) in standard library is for parsing cmd-line arguments:

@file index.go output noplayground

Output above is a result of calling `go run $file -echo echo-arg additional arg`.

## Defining arguments

Let's say your program has an integer `-retries` option.

You register such option with `flag` package using:

```go
var flgRetries int
defaultRetries := 0
usage := "retries specifies number of times to retry the operation"
flag.IntVar(&flgRetries, "retries", defaultRetries, usage)
```

There are functions for other common types:

* `flag.BoolVar`
* `flag.DurationVar`
* `flag.Float64Var`
* `flag.IntVar`, `flag.UIntVar`, `flag.Int64Var`, `flag.UInt64Var`
* `flag.StringVar`

If you register int argument with name `retries`, the way to provide it on cmd-line is `-retries ${value}` or `-retries=${value}`.

POSIX variant `--retries` or Windows variant `/retries` are **not** recognized.

For boolean values you can say: `-help` (implicitly true), `-help=true` or `-help=false`.

`-help false` is not a valid form for boolean variables.

## Parsing and accessing remaining arguments

After parsing arguments, call `flag.Parse()`.

Parsing fails if:

* unknown flag was given on command-line
* a flag didn't parse based on its type (e.g. it was registered as int but the value was not a valid number)

In case of failure, help text describing flags is shown and program exits with error code 2.

You can explicitly print help text using `flag.Usage()`. This is often triggered by `-help` flag.

Help text is based on usage text provided in `flag.IntVar` and others.

Command-line arguments that don't start with `-` are untouched and can be accessed with `flag.Args()`.

## Limitations

Features missing from `flag` package:

* no support for POSIX style `--name`, only `-name` is supported
* no support for short alternatives e.g. `-n` being synonym with `--name`
* no suport for Windows style `/name`

If you need those features, your options are:

* access [raw cmd-line arguments](172)
* use a [third party library](173)
