package container

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
	"woc/filters"
)

type Client interface {
	GetAllContainers()
	GetContainerInfo(types.Container)
}

type dockerClient struct {
	cli *client.Client
}

func (client dockerClient) GetContainerInfo(container *types.Container) (TypeContainer, error) {
	bg := context.Background()
	inspect, err := client.cli.ContainerInspect(bg, container.ID)
	if err != nil {
		logrus.Errorf("fail to get container info from %v,cause:%v", container, err)
		return TypeContainer{}, err
	}
	raw, _, err := client.cli.ImageInspectWithRaw(bg, inspect.Image)
	if err != nil {
		logrus.Errorf("fail to get image info from %v,cause:%v", container, err)
		return TypeContainer{}, err
	}
	return TypeContainer{Old: false, containerInfo: &inspect, imageInfo: &raw}, nil
}

func (client dockerClient) GetAllContainers() ([]TypeContainer, error) {
	containerList, err := client.cli.ContainerList(context.Background(), types.ContainerListOptions{Filters: filters.ArgsFilter()})
	if err != nil {
		return nil, err
	}
	resultContainer := []TypeContainer{}
	for _, container := range containerList {
		containerInfo, err := client.GetContainerInfo(&container)
		if err != nil {
			return nil, err
		}

	}

}

func NewClient() *dockerClient {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	return &dockerClient{
		cli: cli,
	}
}
