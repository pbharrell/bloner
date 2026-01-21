#!/bin/bash

env GOOS=js GOARCH=wasm go build -o bloner.wasm github.com/pbharrell/bloner
