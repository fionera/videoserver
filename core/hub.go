package core

import (
	"github.com/fionera/VideoServer/core/av"
	"github.com/fionera/VideoServer/utils"
)

type Hub struct {
	streams utils.ConcurrentStreamList
}

func NewHub() *Hub {
	return &Hub{
		streams: utils.ConcurrentStreamList{},
	}
}

func (h *Hub) GetStreams() []*av.Stream {
	return h.streams.GetCopy()
}

func (h *Hub) NewInputStream(input av.ReadCloser) *av.Stream {
	return nil
}

func (h *Hub) NewOutputStream(output av.WriteCloser) {

}
