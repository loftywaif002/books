---
Title: Defer gotchas
Search: Defer pitfalls
Id: 130
SOId: rd6000ub
---

When using `defer` keep the following in mind.

## Deferred functions are called at function end

Deferred statements have a function scope, not a block scope.

In other words: deferred calls are executed when exiting a function not when executing block created with `if` or `for` statements.

@file defer_gotcha.go output sha1:56723587ae20f269117036508b8fd5905b3b58e9 goplayground:6ivihpJbkZb

You might expect that deferred statement to be executed when we exit `if` branch but it's executed as the last thing in a function.

## Deferred function arguments

@file defer_gotcha2.go output sha1:fb7c1105477451d2e56e43c432c0779fe6635108 goplayground:YnB5lleWxNa

You might have expected that this code will print `0` and `1`, because those are the values of `i` when we evaluate `defer`.

This is how the code looks like after compiler rewrites the loop to implement `defer`:

```go
var i int
for i = 0; i < 2; i++ {
}

fmt.Printf("%d\n", i)
fmt.Printf("%d\n", i)
```

Now it's clear that when we call deferred `fmt.Printf`, `i` is `2`.

We can fix this by using a [closure](118) to capture the variable:

@file defer_gotcha3.go output sha1:debd2f22d4a3be8fbf97c4d3dbf049cbae4a1cfb goplayground:zTCIuDzpXS9

A closure is more expensive as it requires allocating an object to collect all the variables captured by the closure. In this case that's the price of correctness.
