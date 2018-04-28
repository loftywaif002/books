---
Title: Easy generation of XML struct definition
Id: 334
---

Writing struct definitions for XML parsing can be tedious.

You can use [chidley](https://github.com/gnewton/chidley) to automatically generate struct definitions from sample XML file.

Install the tool with `go get -u github.com/gnewton/chidley`.

Run: `chidley sample.xml`.

This will print Go struct definitions to standard out.

List all options with `chidley`.

Options that I find useful:

* `-X` : sort generated structs by order in XML file (default is alphabetically)
* `-a ""` : prefix for attribute names (default is `Attr`)
* `-e ""` : prefix for struct names (default is `C` i.e. by default struct name for XML element `Api` would `CApi`; this makes it just `Api`)
* `-t` : by default all fields are strings. This flag tries to infer type from values

All options:

```
chidley <flags> xmlFileName|url
xmlFileName can be .gz or .bz2: uncompressed transparently
Usage of chidley:
  -A string
        The tag name attribute to use for the max length Go annotations
  -B    Add database metadata to created Go structs
  -C    Structs have underscores instead of CamelCase; how chidley used to produce output; includes name spaces (see -n)
  -D string
        Base directory for generated Java code (root of maven project) (default "java")
  -F    Assume complete representative XML and collapse tags with only a single string and no attributes
  -G    Only write generated Go structs to stdout (default true)
  -I    If XML decoding error encountered, continue
  -J    Generated Java code for Java/JAXB
  -K    Do not change the case of the first letter of the XML tag names (default true)
  -L    Ignore lower case XML tags
  -M string
        Set name of CDATA string field (default "Text")
  -N string
        The tag name to use for the max length Go annotations
  -P string
        Java package name (rightmost in full package name
  -S string
        The tag name separator to use for the max length Go annotations (default ":")
  -T string
        Field template for the struct field definition. Can include annotations. Default is for XML and JSON (default "{{.GoName}} {{.GoTypeArrayOrPointer}}{{.GoType}} `xml:\"{{if notEmpty .XMLNameSpace}}{{.XMLNameSpace}} {{end}}{{.XMLName}},omitempty\" json:\"{{.XMLName}},omitempty\"`")
  -W    Generate Go code to convert XML to JSON or XML (latter useful for validation) and write it to stdout
  -X    Sort output of structs in Go code by order encounered in source XML (default is alphabetical order)
  -Z int
        The padding on the max length tag attribute
  -a string
        Prefix to attribute names (default "Attr")
  -c    Read XML from standard input
  -d    Debug; prints out much information
  -e string
        Prefix to struct (element) names; must start with a capital (default "C")
  -h string
        List of XML tags to ignore; comma separated
  -k string
        App name for Java code (appended to ca.gnewton.chidley Java package name)) (default "jaxb")
  -m    Validate the field template. Useful to make sure the template defined with -T is valid
  -n    Use the XML namespace prefix as prefix to JSON name
  -p    Pretty-print json in generated code (if applicable)
  -r    Progress: every 50000 input tags (elements)
  -t    Use type info obtained from XML (int, bool, etc); default is to assume everything is a string; better chance at working if XMl sample is not complete
  -u    Filename interpreted as an URL
  -x    Add XMLName (Space, Local) for each XML element, to JSON
```
