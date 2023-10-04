# Twitter Card Image Generator

Generate Twitter card (OGP) images for your blog posts.
Supported front-matters are title, author, categories, tags, and date.
Also, both toml and yaml formats are supported.

![sample](./example/blog-post2.png)

## Installation

### Go version < 1.16

```bash
go get github.com/mathieu-gilloots/tcardgen@latest
```

### Go 1.16+

```bash
go install github.com/mathieu-gilloots/tcardgen@latest
```

## Getting Started

1. Install `tcardgen` command
2. Download your favorite TrueType fonts (the above sample use [KintoSans](https://github.com/ookamiinc/kinto))
3. Create template image (The easyest way is to replace the author image of the template in the [example](./example) directory.)
4. Run the following command

> **NOTE**: `tcardgen` parses a font style from the file name,
> so the font file must follow the naming rule (`<name>-<style>.ttf`), and arrange font files as follows:

```bash
$ tree font/
font/
├── KintoSans-Bold.ttf
├── KintoSans-Medium.ttf
└── KintoSans-Regular.ttf

0 directories, 3 files

$ tcardgen -f path/to/fontDir \
           -o path/to/hugo/static/imgDir \
           -t path/to/templateFile \
           path/to/hugo/content/posts/*.md
```

After successfully executing the command, a PNG image with the same name as the specified content name is generated in the output directory.

## Advanced Generation

If you want to change the color, style, or position of text, you can pass a configuration file with the `--config(-c)` option.
Refer to the [example/template3.config.yaml](example/template3.config.yaml) to see how to configure it.

```bash
$ tcardgen -c example/template3.config.yaml example/blog-post2.md
Load fonts from "font" directory
Load template from "example/template3.png"
Success to generate twitter card into out/blog-post2.png
```

### Result
<img src="./example/template3-config-output.png" width="300">

## OGP setting for Hugo Theme

On my blog, I place the generated images in the `static/tcard` directory. In order to load this image, I set the following OGP information for my blog theme.
If the thumbnail is defined in the post, it is used first. Otherwise, the generated Twitter Card is used. If the page is not blog post, to set the default image.

```html
<!-- General -->
<meta property="og:url" content="{{ .Permalink }}" />
<meta property="og:type" content="{{ if .IsHome }}website{{ else }}article{{ end }}" />
<meta property="og:site_name" content="{{ .Site.Title }}" />
<meta property="og:title" content="{{ .Title }}" />
<meta property="og:description" content="{{ with .Description -}}{{ . }}{{ else -}}{{ if .IsPage }}{{ substr .Summary 0 300 }}{{ else }}{{ with .Site.Params.description }}{{ . }}{{ end }}{{ end }}{{ end }}" />
<meta property="og:image" content="{{ if .Params.thumbnail -}}{{ .Params.thumbnail|absURL }}{{ else if hasPrefix .File.Path "post" -}}{{ path.Join "tcard" (print .File.BaseFileName ".png") | absURL }}{{ else -}}{{ "img/default.png" | absURL }}{{ end -}}" />
<!-- Twitter -->
<meta name="twitter:card" content="summary_large_image" />
<meta name="twitter:site" content="@{{ .Site.Params.twitterName }}" />
```

### Generate images of updated articles

You can generate only the image of the updated article by using `git diff` and `tcardgen`.

```bash
$ git diff --name-only HEAD content/post |\
    xargs tcardgen -o static/tcard -f assets/fonts/kinto-sans -t assets/template.png
```


## Usage

```bash
$ tcardgen -h
Generate TwitterCard(OGP) images for your Hugo posts.
Supported front-matters are title, author, categories, tags, and date.

Usage:
  tcardgen [-f <FONTDIR>] [-o <OUTPUT>] [-t <TEMPLATE>] [-c <CONFIG>] <FILE>...

Examples:
# Generate a image and output to the example directory.
tcardgen --fontDir=font --output=example --template=example/template.png example/blog-post.md

# Generate a image and output to the example directory as "featured.png".
tcardgen --fontDir=font --output=example/featured.png --template=example/template.png example/blog-post.md

# Generate multiple images.
tcardgen --template=example/template.png example/*.md

# Genrate an image based on the drawing configuration.
tcardgen --config=config.yaml example/*.md

Flags:
  -c, --config string     Set a drawing configuration file.
  -f, --fontDir string    Set a font directory. (default "font")
  -h, --help              help for tcardgen
      --outDir string     (DEPRECATED) Set an output directory.
  -o, --output string     Set an output directory or filename (only png format). (default "out")
  -t, --template string   Set a template image file. (default example/template.png)
```
