package container

import (
	"github.com/docker/docker/api/types"
	"strconv"
)

const (
	enableLabel     = "com.woc.enableLabel.enable"
	enableHttpLabel = "com.woc.enableLabel.enable"
	wocLabel        = "com.woc.wocLabel.enable"
)

type TypeContainer struct {
	Old           bool
	ContainerInfo *types.ContainerJSON
	ImageInfo     *types.ImageInspect
}

type Container interface {
	Name() string
	IsWoc() bool
	Enabled() bool
	HttpEnable() bool
}

func (c TypeContainer) HttpEnable() bool {
	val, ok := c.ContainerInfo.Config.Labels[enableHttpLabel]
	if !ok {
		return false
	}

	parsedBool, err := strconv.ParseBool(val)
	if err != nil {
		return false
	}

	return parsedBool
}

func (c TypeContainer) Name() string {
	return c.ContainerInfo.Name
}

func (c TypeContainer) IsWoc() bool {
	s, ok := c.ContainerInfo.Config.Labels[wocLabel]
	return ok && s == "true"
}

func (c TypeContainer) Enabled() bool {
	val, ok := c.ContainerInfo.Config.Labels[enableLabel]
	if !ok {
		return false
	}

	parsedBool, err := strconv.ParseBool(val)
	if err != nil {
		return false
	}

	return parsedBool
}
