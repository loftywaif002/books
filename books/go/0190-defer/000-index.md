---
Title: Defer
Id: 128
SOId: 2795
---

In a complicated function, it's easy to forgot to release a resource (e.g. to close a file handle or to unlock a mutex).

In C++ you would use RAII to ensure a resource is always released.

In Go you use `defer` statement.

```go
func foo() {
  f, err := os.Open("myfile.txt")
  if err != nil {
    return
  }
  defer f.Close()

  // ... lots of code
}
```

In the above example, `defer f.Close()` ensures that `f.Close()` will be called before we exit `foo`, even if a [panic](131) happens.

Placing `defer f.Close()` right after `os.Open()` makes it easy to audit the code to verify `Close` is always called.

This is especially useful for large functions with multiple exit points.

If deferred code is more complicated, you can use a function literal:

```go
func foo() {
  mutex1.Lock()
  mutex2.Lock()

  defer func() {
    mutex2.Unlock()
    mutex1.Unlock()
  }()

  // ... more code
}
```

You can use multiple `defer` statements. They'll be called in reverse order.
