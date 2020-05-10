package utils

import (
	"github.com/fionera/VideoServer/core/av"
	"sync"
)

type ConcurrentStreamList struct {
	mtx sync.RWMutex
	data []*av.Stream
}

func (c *ConcurrentStreamList) GetCopy() []*av.Stream {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	dataCopy := make([]*av.Stream, len(c.data))
	copy(dataCopy, c.data)

	return dataCopy
}
