package container

import (
	"github.com/docker/docker/api/types"
	"strconv"
)

const (
	enableLabel = "com.woc.enableLabel.enable"
	wocLabel    = "com.woc.wocLabel.enable"
)

type TypeContainer struct {
	Old           bool
	containerInfo *types.ContainerJSON
	imageInfo     *types.ImageInspect
}

type Container interface {
	Name() string
	IsWoc() bool
	Enabled() (bool, bool)
}

func (c TypeContainer) Name() string {
	return c.containerInfo.Name
}

func (c TypeContainer) IsWoc() bool {
	s, ok := c.containerInfo.Config.Labels[wocLabel]
	return ok && s == "true"
}

func (c TypeContainer) Enabled() (bool, bool) {
	val, ok := c.containerInfo.Config.Labels[enableLabel]
	if !ok {
		return false, false
	}

	parsedBool, err := strconv.ParseBool(val)
	if err != nil {
		return false, false
	}

	return parsedBool, true
}
