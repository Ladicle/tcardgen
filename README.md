# Twitter Card Image Generator

Generate Twitter card image for your blog posts.

![sample](./example/blog-post.png)

## Installation

`go get github.com/Ladicle/tcardgen`

## Usage

```
$ tcardgen -h
Generate twitter card image from the Hugo post.

Usage:
  tcardgen [-f <FONTDIR>] [-o <OUTDIR>] [-t <TEMPLATE>] <FILE>...

Flags:
  -f, --font string       Set a font directory. (default "font")
  -h, --help              help for tcardgen
  -o, --out string        Set an output directory. (default "out")
  -t, --template string   Set a template image file. (default "template.png")
```
