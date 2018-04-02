---
Title: Strings
Id: 40
SOId: 9666
---

Strings in Go are immutable sequences of bytes.

Unlike languages like Python or Java, they are not internally represented as Unicode. Consequently, when reading strings from files or network connections, there is no conversion step from bytes to internal representation. When writing strings to files, there is no conversion to a code page.

Go strings don't assume any particular code page. They are just bytes.

Go source code files are always UTF-8 so strings defined in source code are also valid utf8 strings.

Additionally, functions in the standard library that involve converting characters to upper-case or lower-case etc. assume that raw bytes represent UTF-8-encoded Unicode strings and perform transformations using Unicode rules.

Basic string usage:

@file index.go output sha1:643856e224099ead3e2fc62a2e3c19ee17f62374 goplayground:RgP9gvhvXqF

Important standard library packages for working on strings:

* [strings](https://golang.org/pkg/strings/) implements string searching, splitting, case conversions
* [bytes](https://golang.org/pkg/bytes/) has the same functionality as `strings` package but operates on `[]byte` byte slices
* [strconv](https://golang.org/pkg/strconv/) for conversion between strings and integer and float numbers
* [unicode/utf8](https://golang.org/pkg/unicode/utf8/) decodes from UTF-8-encoded strings and encodes to UTF-8-encoded strings
* [regexp](https://golang.org/pkg/regexp/) implements regular expressions
* [text/scanner](https://golang.org/pkg/text/scanner/) for scanning and tokenizing UTF-8-encoded text
* [text/template](https://golang.org/pkg/text/template/) for generating larger strings from templates
* [html/template](https://golang.org/pkg/html/template/) has all the functionality of `text/template` but understands the structure of HTML for generation of HTML that is safe from code injection attacks
