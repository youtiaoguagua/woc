package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
	"woc/action"
	cli "woc/client"
	filter "woc/filters"
)

var (
	client cli.Client
)

func init() {

}

func main() {
	containerNames := viper.GetString("container.names")
	names := strings.Split(containerNames, ",")
	client := cli.NewClient()
	filter := filter.BuildFilter(names)
	logrus.Debugf("filter list is %v", filter)
	action.CheckMultiWoc(client)
}
