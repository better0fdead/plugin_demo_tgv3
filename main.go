package main

import (
	_ "embed"
	"log"

	"github.com/better0fdead/plugin_demo_tgv3/plugin"
)

//go:embed about.md
var about []byte

// Plugin version
const Version = "v0.0.1"

// Plugin description
const description = "plugin for tgv3 demo"

// Plugin source url
const source = "github.com/better0fdead/plugin_demo_tgv3"

// Plugin help message
func help() string {
	helpMsg := description + "\n"

	return helpMsg
}

func main() {
	err := plugin.Start(about, source, description, Version, help())
	if err != nil {
		log.Fatalln(err)
	}
}
