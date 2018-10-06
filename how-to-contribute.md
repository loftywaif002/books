## How to contribute


### Making a change

All books are managed on [Notion.so](https://www.notion.so/kjkpublic/All-Books-9d463535c38f45c592756211e56977cf).

To make a small change to an article, you can click on `Suggest an edit` link at top of each page. This opens the file on Notion.so. If you're logged in to Notion, you can add a comment.

### Suggesting a change

For general ideas on how to improve the books, the project etc., use [issue tracker](https://github.com/essentialbooks/books/issues)

### Making more changes

Source code examples for articles are hosted on GitHub.

Toolchain for building the books (i.e. converting markdown sources and source files into HTML) is written in go, so you'll need to [install Go](https://golang.org/dl/).

For cross-platform portability, helper scripts are written in PowerShell, so if you're not on Windows, you'll have to [install it too](https://github.com/PowerShell/PowerShell). Or you can write equivalent bash script. They are trivial.

You'll also need an editor. I use [Visual Studio Code](https://code.visualstudio.com/) with [Code Runner](https://marketplace.visualstudio.com/items?itemName=formulahendry.code-runner) and [Terminal Here](https://marketplace.visualstudio.com/items?itemName=Tyriar.vscode-terminal-here) extensions.

A crucial tool is `./s/preview.ps1`.

It rebuilds all HTML, starts a web server for local preview of changes.

### What to improve?

Each book can be improved by adding more articles.

Areas that are likely to always need improvement:
* examples for common tasks (e.g. parsing json/xml/cvs/markdown files, accessing databases etc.)
* examples for popular third-party libraries
