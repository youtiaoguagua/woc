package action

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"sort"
	cli "woc/client"
	"woc/container"
	sorter "woc/container/sort"
	"woc/filters"
)

func CheckMultiWoc(client cli.Client) error {
	containers, err := client.GetAllContainers(filters.FilterByWoc(filters.FinalFilter))
	if err != nil {
		return err
	}
	if len(containers) <= 0 {
		return nil
	}
	logrus.Debug("find multiply woc ,them will be clean next!")
	return cleanUpWoc(containers, client)
}

func cleanUpWoc(containers []container.TypeContainer, client cli.Client) error {
	var stopErrors int

	sort.Sort(sorter.ByCreated(containers))

	needCleanUpContainers := containers[0 : len(containers)-1]
	for _, container := range needCleanUpContainers {
		err := client.StopContainer(container)
		logrus.WithError(err).Error("Could not stop a previous woc instance.")
		stopErrors++
	}

	if stopErrors > 0 {
		return fmt.Errorf("%d errors while stopping woc containers", stopErrors)
	}

	return nil
}
