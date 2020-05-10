package api

import (
	"encoding/json"
	"github.com/fionera/VideoServer/core/av"
	"github.com/fionera/VideoServer/protocol/rtmp"
	"net/http"
)

type WebApi struct {
	handler av.Handler
}

func NewWebApi(h av.Handler) *WebApi {
	return &WebApi{
		handler: h,
	}
}

type Stream struct {
	Key string `json:"key"`
	Id  string `json:"id"`
}

type StreamList struct {
	Publishers []Stream `json:"publishers"`
	Players    []Stream `json:"players"`
}

func (w *WebApi) ListStreams(writer http.ResponseWriter, request *http.Request) {
	streams := w.handler.GetStreams()

	streamList := StreamList{}
	for item := range rtmpStream.GetStreams().IterBuffered() {
		if s, ok := item.Val.(*rtmp.Stream); ok {
			if s.GetReader() != nil {
				msg := Stream{item.Key, s.GetReader().Info().UID}
				streamList.Publishers = append(streamList.Publishers, msg)
			}
		}
	}

	for item := range rtmpStream.GetStreams().IterBuffered() {
		ws := item.Val.(*rtmp.Stream).GetWs()
		for s := range ws.IterBuffered() {
			if pw, ok := s.Val.(*rtmp.PackWriterCloser); ok {
				if pw.GetWriter() != nil {
					msg := Stream{item.Key, pw.GetWriter().Info().UID}
					streamList.Players = append(streamList.Players, msg)
				}
			}
		}
	}

	resp, _ := json.Marshal(streamList)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}