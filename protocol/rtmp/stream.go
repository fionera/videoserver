package rtmp

import (
	"github.com/fionera/VideoServer/core"
	"github.com/fionera/VideoServer/core/av"
	"github.com/fionera/VideoServer/protocol/rtmp/cache"
	"github.com/orcaman/concurrent-map"
)

var (
	EmptyID = ""
)

type RtmpStream struct {
	stream *Stream
	info   *av.Info
	hub    *core.Hub
}

func NewRtmpStream(hub *core.Hub) *RtmpStream {
	ret := &RtmpStream{
		hub: hub,
	}
	return ret
}

type Stream struct {
	isStart bool
	cache   *cache.Cache
	r       av.ReadCloser
	ws      cmap.ConcurrentMap
	info    av.Info
}

type PackWriterCloser struct {
	init bool
	w    av.WriteCloser
}

func (p *PackWriterCloser) GetWriter() av.WriteCloser {
	return p.w
}

func NewStream() *Stream {
	return &Stream{
		cache: cache.NewCache(),
		ws:    cmap.New(),
	}
}

func (s *Stream) GetReader() av.ReadCloser {
	return s.r
}

func (s *Stream) GetWs() cmap.ConcurrentMap {
	return s.ws
}

func (s *Stream) Copy(dst *Stream) {
	for item := range s.ws.IterBuffered() {
		v := item.Val.(*PackWriterCloser)
		s.ws.Remove(item.Key)
		v.w.CalcBaseTimestamp()
		dst.AddWriter(v.w)
	}
}

func (s *Stream) AddReader(r av.ReadCloser) {
	s.r = r
	go s.TransStart()
}

func (s *Stream) AddWriter(w av.WriteCloser) {
	pw := &PackWriterCloser{w: w}
	s.ws.Set(s.info.UID, pw)
}

func (s *Stream) TransStart() {
	s.isStart = true
	var p av.Packet

	log.Printf("TransStart:%v", s.info)

	for {
		if !s.isStart {
			s.closeInter()
			return
		}
		err := s.r.Read(&p)
		if err != nil {
			s.closeInter()
			s.isStart = false
			return
		}

		s.cache.Write(p)

		for item := range s.ws.IterBuffered() {
			v := item.Val.(*PackWriterCloser)
			if !v.init {
				//log.Printf("cache.send: %v", v.w.Info())
				if err = s.cache.Send(v.w); err != nil {
					//log.Printf("[%s] send cache packet error: %v, remove", v.w.Info(), err)
					s.ws.Remove(item.Key)
					continue
				}
				v.init = true
			} else {
				new_packet := p
				//writeType := reflect.TypeOf(v.w)
				//log.Printf("w.Write: type=%v, %v", writeType, v.w.Info())
				if err = v.w.Write(&new_packet); err != nil {
					//log.Printf("[%s] write packet error: %v, remove", v.w.Info(), err)
					s.ws.Remove(item.Key)
				}
			}
		}
	}
}

func (s *Stream) TransStop() {
	log.Printf("TransStop: %s", s.info.UID)

	if s.isStart && s.r != nil {
		s.r.Close()//errors.New("stop old"))
	}

	s.isStart = false
}

func (s *Stream) CheckAlive() (n int) {
	if s.r != nil && s.isStart {
		if s.r.Alive() {
			n++
		} else {
			s.r.Close()//errors.New("read timeout"))
		}
	}
	for item := range s.ws.IterBuffered() {
		v := item.Val.(*PackWriterCloser)
		if v.w != nil {
			if !v.w.Alive() && s.isStart {
				s.ws.Remove(item.Key)
				v.w.Close()//errors.New("write timeout"))
				continue
			}
			n++
		}

	}
	return
}

func (s *Stream) closeInter() {
	for item := range s.ws.IterBuffered() {
		v := item.Val.(*PackWriterCloser)
		if v.w != nil {
			//if v.w.Info().IsInterval() {
			//	v.w.Close(errors.New("closed"))
			//	s.ws.Remove(item.Key)
			//	log.Printf("[%v] player closed and remove\n", v.w.Info())
			//}
		}

	}
}
