package main

import (
	"github.com/hashicorp/packer/packer/plugin"
	"github.com/saymedia/packer-post-processor-template/template"
)

func main() {
	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}
	server.RegisterPostProcessor(new(template.PostProcessor))
	server.Serve()
}
