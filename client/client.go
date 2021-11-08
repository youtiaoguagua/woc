package client

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
	containerModel "woc/container"
	"woc/filters"
)

type Client interface {
	GetAllContainers()
	GetContainerInfo(types.Container)
}

type dockerClient struct {
	cli *client.Client
}

func (client dockerClient) GetContainerInfo(container *types.Container) (containerModel.TypeContainer, error) {
	bg := context.Background()
	inspect, err := client.cli.ContainerInspect(bg, container.ID)
	if err != nil {
		logrus.Errorf("fail to get container info from %v,cause:%v", container, err)
		return containerModel.TypeContainer{}, err
	}
	raw, _, err := client.cli.ImageInspectWithRaw(bg, inspect.Image)
	if err != nil {
		logrus.Errorf("fail to get image info from %v,cause:%v", container, err)
		return containerModel.TypeContainer{}, err
	}
	return containerModel.TypeContainer{Old: false, ContainerInfo: &inspect, ImageInfo: &raw}, nil
}

func (client dockerClient) GetAllContainers(func(container containerModel.TypeContainer) bool) ([]containerModel.TypeContainer, error) {
	containerList, err := client.cli.ContainerList(context.Background(), types.ContainerListOptions{Filters: filters.ArgsFilter()})
	logrus.Infof("all container is %v", containerList)
	if err != nil {
		return nil, err
	}
	resultContainer := []containerModel.TypeContainer{}
	//for _, container := range containerList {
	//	containerInfo, err := client.GetContainerInfo(&container)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//}
	return resultContainer, nil

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
