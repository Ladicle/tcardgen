# Twitter Card Image Generator

Generate Twitter card (OGP) images for your blog posts.
Supported front-matters are title, author, categories, tags, and date.
Also, both toml and yaml formats are supported.

![sample](./example/blog-post2.png)

## Installation

```
go get github.com/Ladicle/tcardgen
```

## Getting Started

1. Install `tcardgen` command
2. Download your favorite TrueType fonts (the above sample use [KintoSans](https://github.com/ookamiinc/kinto))
3. Create template image (The easyest way is to replace the author image of the template in the [example](./example) directory.)
4. Run the following command

```
$ tcardgen -f path/to/fontDir \
           -o path/to/hugo/static/imgDir \
           -t path/to/templateFile \
           path/to/hugo/content/posts/*.md
```

After successfully executing the command, a PNG image with the same name as the specified content name is generated in the output directory.

## Usage

```
$ tcardgen -h
Generate TwitterCard(OGP) images for your Hugo posts.
Supported front-matters are title, author, categories, tags, and date.

Usage:
  tcardgen [-f <FONTDIR>] [-o <OUTDIR>] [-t <TEMPLATE>] <FILE>...

Examples:
# Generate a image and output to the example directory.
tcardgen --fontDir=font --outDir=example --template=example/template.png example/blog-post.md

# Generate multiple images.
tcardgen --template=example/template.png example/*.md

Flags:
  -f, --fontDir string    Set a font directory. (default "font")
  -h, --help              help for tcardgen
  -o, --outDir string     Set an output directory. (default "out")
  -t, --template string   Set a template image file. (default "template.png")
```
