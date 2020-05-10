package av

import (
	"fmt"
	"net/url"
)

type Stream struct {
	info *Info
}

func NewStream(info *Info) *Stream {
	return &Stream{
		info: info,
	}
}

// Info is the Information about a Stream, ignoring if its an Input or an Output.
// Key is a user definable route/key for identifying the stream.
// URL is the full URL from which the Stream comes / goes to
// UID is a unique string identifying this unique Stream
// Interval TODO: Find out what this is
type Info struct {
	URL      *url.URL
	UID      string
	Interval bool
}

func (info Info) IsInterval() bool {
	return info.Interval
}

func (info Info) String() string {
	return fmt.Sprintf("<URL: %s, UID: %s, Interval: %v>",
		info.URL, info.UID, info.Interval)
}
