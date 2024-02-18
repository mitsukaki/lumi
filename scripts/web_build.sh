#!/bin/bash

# Change directory to the "web" directory
cd /web

# Run the build command using "vite"
npx vite build . --outDir ../build/public --emptyOutDir

