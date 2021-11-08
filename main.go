package main

import (
	"github.com/docker/docker/client"
	"woc/container"
)

var (
	name *client.Client
)

func main() {
	client := container.NewClient()
	client.GetAllContainers()
}
