---
Title: Find string in another string
Search: searching for strings
Id: 41
SOId: 30331
---

## Find the position of a string within another string

Using [`strings.Index`](https://golang.org/pkg/strings/#Index):

@file finding_strings.go output sha1:a47b4a300ddca07351a2b42f7e66bd358f175ca5 goplayground:1c8V7OP6HJh

## Find the position of string within another string starting from the end

Above we searched from the beginning of the string. We can also search from the end using [`strings.LastIndex`](https://golang.org/pkg/strings/#LastIndex):

@file finding_strings2.go output sha1:9fa45203c4ae55075c6b726b2692c813240971e9 goplayground:iyxoEzHw4p7

## Find all occurrences of a substring

Above, we only found the first occurrence of the substring. Here's how to find all occurrences:

@file finding_strings3.go output sha1:cd32f1eaba4a825c9d36f6fe786dc760a3caa48e goplayground:C9VJxywxROJ

<!-- TODO: example using a regular expression -->

## Check if a string contains another string

Using [`strings.Contains`](https://golang.org/pkg/strings/#Contains):

@file finding_strings4.go output sha1:6615c90c15775f6f402a176ae0d6209f494cea77 goplayground:BSLr-qaYPvJ

## Check if a string starts with another string

Using [`strings.HasPrefix`](https://golang.org/pkg/strings/#HasPrefix):

@file finding_strings5.go output sha1:babb86fb2fbae4369e20ccc5239230d606280d66 goplayground:WgVWyWBQ_gu

## Check if a string ends with another string

Using [`strings.HasSuffix`](https://golang.org/pkg/strings/#HasSuffix):

@file finding_strings6.go output sha1:d183ea2089b304e519b217fbe51f5e4af1d8f6be goplayground:bZWDWmkNklX
