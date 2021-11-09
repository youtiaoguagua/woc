package client

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"time"
	containerModel "woc/container"
	"woc/container/util"
	"woc/filters"
)

type Client interface {
	GetAllContainers(filters.Filter) ([]containerModel.TypeContainer, error)
	GetContainerInfo(*types.Container) (containerModel.TypeContainer, error)
	StopContainer(container containerModel.TypeContainer) error
}

type dockerClient struct {
	cli *client.Client
}

func (client dockerClient) StopContainer(container containerModel.TypeContainer) error {
	id := container.ContainerInfo.ID
	shortID := util.ShortId(id)

	if container.ContainerInfo.State.Running {
		logrus.Infof("Stopping %s (%s) with SIGTERM", container.Name(), shortID)
		if err := client.cli.ContainerKill(context.Background(), container.ContainerInfo.ID, "SIGTERM"); err != nil {
			return err
		}
	}

	_ = client.SleepCheck(container, 10*time.Minute)

	if container.ContainerInfo.HostConfig.AutoRemove {
		logrus.Debugf("AutoRemove container %s, skipping ContainerRemove call.", shortID)
	} else {
		logrus.Debugf("Removing container %s", shortID)
		if err := client.cli.ContainerRemove(context.Background(), container.ContainerInfo.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
			return err
		}
	}

	if err := client.SleepCheck(container, 10*time.Minute); err == nil {
		return fmt.Errorf("container %s (%s) could not be removed", container.Name(), shortID)
	}
	return nil
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

//GetAllContainers 获取所有需要过滤的容器
func (client dockerClient) GetAllContainers(filter filters.Filter) ([]containerModel.TypeContainer, error) {
	containerList, err := client.cli.ContainerList(context.Background(), types.ContainerListOptions{Filters: filters.ArgsFilter()})
	if err != nil {
		return nil, err
	}

	containerListJson, _ := jsoniter.MarshalToString(containerList)
	logrus.Info("all container is :", containerListJson)

	var resultContainer []containerModel.TypeContainer
	for _, container := range containerList {
		containerInfo, err := client.GetContainerInfo(&container)
		if err != nil {
			return nil, err
		}
		if !filter(containerInfo) {
			continue
		}

		resultContainer = append(resultContainer, containerInfo)
	}
	return resultContainer, nil

}

func (client dockerClient) SleepCheck(container containerModel.TypeContainer, duration time.Duration) error {
	timeout := time.After(duration)

	for {
		select {
		case <-timeout:
			return nil
		default:
			if ci, err := client.cli.ContainerInspect(context.Background(), container.ContainerInfo.ID); err != nil {
				return err
			} else if !ci.State.Running {
				return nil
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func NewClient() Client {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	return dockerClient{
		cli: cli,
	}
}
