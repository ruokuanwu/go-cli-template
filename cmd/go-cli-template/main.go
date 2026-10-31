package main

import (
	"os"

	"example.com/go-cli-template/internal/app"
)

var version = "dev"

func main() {
	os.Exit(app.Run(os.Args[1:], os.Stdout, os.Stderr, version))
}
