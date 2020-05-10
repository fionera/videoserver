package rtmp

import (
	"github.com/fionera/VideoServer/core"
	rtmpCore "github.com/fionera/VideoServer/protocol/rtmp/protocol"
	"net"
)

type RtmpServer struct {
	options RtmpOptions
	hub     * core.Hub
}

type RtmpOptions struct {
	enableInput  bool
	enableOutput bool
}

func NewRtmpIO(hub *core.Hub, opt RtmpOptions) *RtmpServer {
	return &RtmpServer{
		options: opt,
		hub: hub,
	}
}

func (s *RtmpServer) Serve(listener net.Listener) (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("rtmp serve panic: %s", r)
		}
	}()

	for {
		var netconn net.Conn
		netconn, err = listener.Accept()
		if err != nil {
			return
		}
		conn := rtmpCore.NewConn(netconn, 4*1024)
		log.Printf("new client, connect remote: %s - local: %s", conn.RemoteAddr().String(), conn.LocalAddr().String())
		go s.handleConn(conn)
	}
}

func (s *RtmpServer) handleConn(conn *rtmpCore.Conn) error {
	if err := conn.HandshakeServer(); err != nil {
		conn.Close()
		log.Printf("handleConn HandshakeServer err: %s", err)
		return err
	}
	connServer := rtmpCore.NewConnServer(conn)

	if err := connServer.ReadMsg(); err != nil {
		conn.Close()
		log.Printf("handleConn read msg err: %s", err)
		return err
	}

	appname, _, _ := connServer.GetInfo()
	log.Printf("Got connection on app %s", appname)

	//if ret := configure.CheckAppName(appname); !ret {
	//	err := errors.New(fmt.Sprintf("application name=%s is not configured", appname))
	//	conn.Close()
	//	log.Println("CheckAppName err:", err)
	//	return err
	//}

	log.Printf("handleConn: IsPublisher=%v", connServer.IsPublisher())
	if connServer.IsPublisher() {
		if !s.options.enableInput {
			return connServer.Close()
		}

		reader := NewVirReader(connServer)
		log.Printf("new publisher: %+v", reader.Info())

		s.hub.NewInputStream(reader)
	} else {
		if !s.options.enableOutput {
			return connServer.Close()
		}

		writer := NewVirWriter(connServer)
		log.Printf("new player: %+v", writer.Info())
		s.hub.NewOutputStream(writer)
	}

	return nil
}
