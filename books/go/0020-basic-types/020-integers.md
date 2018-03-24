---
Title: Integers
Search: ints
Id: 14
SOId: rd6000kf
---

Go has fixed-size signed and unsigned integers:

* `int8`, `uint8`
* `byte` is an alias for `uint8`
* `int16`, `uint16`
* `int32`, `uint32`
* `rune` is an alias for `int32`
* `int64`, `uint64`

It also has architecture-dependent integers:

* `int` is `int32` on 32-bit processors and `int64` on 64-bit processors
* `uint` is `uint32` on 32-bit processors and `uint64` on 64-bit processors

[Zero value](29) of an integer is 0.
