package main

import (
	"github.com/mitchellh/packer/packer/plugin"
	"github.com/saymedia/packer-post-processor-terraform/terraform"
)

func main() {
	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}
	server.RegisterPostProcessor(new(terraform.PostProcessor))
	server.Serve()
}
