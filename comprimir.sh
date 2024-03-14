#!/bin/bash
# se setean las variables para que el compilador
# este preparado para amazon
#GOOS=linux GOARCH=amd64
# se compila go y se comprime en un archivo main.zip
#go build -o main && zip -r main.zip main
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main
#GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o main
zip -r main.zip main