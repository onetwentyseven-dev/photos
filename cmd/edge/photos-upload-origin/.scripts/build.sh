#!/bin/bash

set -e
rm -f bundle.zip bundle.js
esbuild index.js --bundle --platform=node --target=node18 --external:sharp --outfile=bundle.js
# esbuild index.js --bundle --minify --platform=node --target=node18 --external:sharp --outfile=bundle.js
mv package.json package.json.bck && mv package-lock.json package-lock.json.bck && rm -rf node_modules
SHARP_IGNORE_GLOBAL_LIBVIPS=1 npm install --arch=x64 --platform=linux --libc=glibc sharp
zip -qr photos-upload-origin.zip bundle.js node_modules/ && ls -lha photos-upload-origin.zip
mv package.json.bck package.json && mv package-lock.json.bck package-lock.json && rm -rf node_modules && npm install
