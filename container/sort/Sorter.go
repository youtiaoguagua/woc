package sorter

import (
	"time"
	"woc/container"
)

type ByCreated []container.TypeContainer

func (c ByCreated) Len() int { return len(c) }
func (c ByCreated) Less(i, j int) bool {
	t1, err := time.Parse(time.RFC3339Nano, c[i].ContainerInfo.Created)
	if err != nil {
		t1 = time.Now()
	}

	t2, _ := time.Parse(time.RFC3339Nano, c[j].ContainerInfo.Created)
	if err != nil {
		t1 = time.Now()
	}

	return t1.Before(t2)
}
func (c ByCreated) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
