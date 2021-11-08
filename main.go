package main

import (
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
	client3 "woc/client"
	client2 "woc/filters"
)

var (
	name *client.Client
)

func init() {

}

func main() {
	containerNames := viper.GetString("container.names")
	names := strings.Split(containerNames, ",")
	client := client3.NewClient()
	filter := client2.BuildFilter(names)
	containers, err := client.GetAllContainers(filter)
	if err != nil {
		panic(err)
	}
	logrus.Debugf("all filter contrainers is %v", containers)
}
