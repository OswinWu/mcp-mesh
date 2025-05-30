package main

import (
	"flag"
	"fmt"
	"mcp-mesh/config"
)

func main() {
	configPath := flag.String("config", `config.yaml`, "filePath of config")
	flag.Parse()
	config.Init(*configPath)
	fmt.Println(config.Get())
}
