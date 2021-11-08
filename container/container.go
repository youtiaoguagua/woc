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
	ContainerInfo *types.ContainerJSON
	ImageInfo     *types.ImageInspect
}

type Container interface {
	Name() string
	IsWoc() bool
	Enabled() (bool, bool)
}

func (c TypeContainer) Name() string {
	return c.ContainerInfo.Name
}

func (c TypeContainer) IsWoc() bool {
	s, ok := c.ContainerInfo.Config.Labels[wocLabel]
	return ok && s == "true"
}

func (c TypeContainer) Enabled() (bool, bool) {
	val, ok := c.ContainerInfo.Config.Labels[enableLabel]
	if !ok {
		return false, false
	}

	parsedBool, err := strconv.ParseBool(val)
	if err != nil {
		return false, false
	}

	return parsedBool, true
}
