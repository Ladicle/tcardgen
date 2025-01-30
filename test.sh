#!/usr/bin/env bash

set -xeo pipefail

FNAME=KintoSans.zip
VERSION=${VERSION:-v1.0.1}

# Download Fonts to `font` directory
if [ ! -f font/KintoSans-Regular.ttf ]; then
    mkdir -p font
    wget https://github.com/ookamiinc/kinto/releases/download/${VERSION}/${FNAME}
    unzip -j $FNAME "*.ttf" -d font
    rm $FNAME
fi

# Generate an image and check diff
[ -d test ] && rm -r test
mkdir -p test
for name in "blog-post" "blog-post2"; do
    echo "Test $name"
    go run main.go \
       -c example/default.config.yaml \
       -f font \
       -o test/ \
       -t example/template.png \
       example/$name.md
    diff test/$name.png example/$name.png
done
