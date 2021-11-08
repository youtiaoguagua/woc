package filters

import (
	f "github.com/docker/docker/api/types/filters"
	"strings"
	"woc/container"
)

type Filter func(container.TypeContainer) bool

func ArgsFilter() f.Args {
	args := f.NewArgs(
		f.Arg("status", "running"),
	)
	return args
}

func FinalFilter(container container.TypeContainer) bool {
	return true
}

//filterByName 过滤带有名字的容器
func filterByName(filter Filter, names []string) Filter {
	if len(names) == 0 {
		return filter
	}
	return func(container container.TypeContainer) bool {
		for _, name := range names {
			if name == container.Name() || strings.TrimSpace(container.Name()) == name {
				return filter(container)
			}
		}
		return false
	}
}

//filterLabel 过滤带有标签的容器
func filterByLabel(filter Filter) Filter {
	return func(container container.TypeContainer) bool {
		enabled, b := container.Enabled()
		if enabled && b {
			return true
		}
		return filter(container)
	}
}

//filterByWoc 过滤本容器
func filterByWoc(filter Filter) Filter {
	return func(container container.TypeContainer) bool {
		if container.IsWoc() {
			return false
		}
		return filter(container)
	}
}

//BuildFilter 构建过滤器
func BuildFilter(name []string) func(container container.TypeContainer) bool {
	filter := FinalFilter
	filter = filterByLabel(filter)
	filter = filterByName(filter, name)
	filter = filterByWoc(filter)
	return filter
}
