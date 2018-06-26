---
Title: Access documentation offline
Id: 8
SOId: 24998
---

To browse full documentation locally, run:

```sh
godoc -http=localhost:9000
```

Then open [localhost:9000](http://localhost:9000) in the browser. It'll show the same information as [https://golang.org/doc/](https://golang.org/doc/).

You can run guided tour locally with:

```sh
go get golang.org/x/tour/gotour
go tool tour
```

This is the equivalent of [https://tour.golang.org/](https://tour.golang.org/).

You can use `godoc` for quick reference. For example, to see documentation for fmt.Print:

```sh
godoc cmd/fmt Print
```

General help is also available from the command-line:

```sh
go help [command]
```

For example, to see documentation for `go build` use `go help build`.
