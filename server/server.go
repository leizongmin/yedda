package server

import (
	"github.com/leizongmin/simple-limiter-service/protocol"
	"github.com/leizongmin/simple-limiter-service/service"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

const DefaultListenNetwork = "tcp"
const DefaultListenAddress = "127.0.0.1:16789"
const Version = "1.0.0-alpha"

type Server struct {
	listener  net.Listener
	closed    bool
	service   *service.Service
	options   Options
	enableLog bool
	connMap   *sync.Map
}

type Options struct {
	Network      string
	Address      string
	TimeAccuracy time.Duration
	DatabaseSize uint32
	EnableLog    bool
}

func NewServer(options Options) (s *Server, err error) {
	if err != nil {
		return s, err
	}
	if len(options.Address) < 1 {
		options.Address = DefaultListenAddress
	}
	if len(options.Network) < 1 {
		options.Network = DefaultListenNetwork
	}
	listener, err := net.Listen(options.Network, options.Address)
	if err != nil {
		return s, err
	}
	s = &Server{
		listener: listener,
		closed:   false,
		service: service.NewService(service.Options{
			TimeAccuracy: options.TimeAccuracy,
			DatabaseSize: options.DatabaseSize,
		}),
		options:   options,
		enableLog: options.EnableLog,
		connMap:   &sync.Map{},
	}
	s.Log("[main]", "server version: %s, protocol version: %d", Version, protocol.CurrentVersion)
	s.Log("[main]", "dbsize: %d, ", options.DatabaseSize)
	s.Log("[main]", "server listen on %s", listener.Addr())
	return s, err
}

func (s *Server) Addr() net.Addr {
	return s.listener.Addr()
}

func (s *Server) Close() error {
	err := s.listener.Close()
	if err == nil {
		s.closed = true
		s.service.Destroy()
		s.Log("[main]", "server closed")
	}
	return err
}

func (s *Server) Loop() error {
	s.service.Start()
	defer s.service.Stop()
	for !s.closed {
		conn, err := s.listener.Accept()
		if err != nil {
			return err
		}
		s.handleConnection(conn)
	}
	return nil
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	addr := conn.RemoteAddr()
	s.connMap.Store(addr, conn)
	defer func() {
		s.connMap.Delete(addr)
	}()
	prefix := "[" + addr.String() + "]"
	s.Log(prefix, "Connected. There are currently %d connections", syncMapLen(s.connMap))
	for {
		var err error
		p, err := protocol.NewPackageFromReader(conn)
		if err != nil {
			if err == io.EOF {
				s.Log(prefix, "Connection closed")
			} else {
				s.Log(prefix, "Read error: %s", err)
			}
			break
		}
		if p.Version != protocol.CurrentVersion {
			s.Log(prefix, "Ignore protocol version %d", p.Version)
			continue
		}
		switch p.Op {
		case protocol.OpPing:
			err = s.Reply(conn, p.ID, protocol.OpPong, p.Data)
		case protocol.OpPong:
			// do nothing
		case protocol.OpGet:
			a, err := service.NewCmdArgFromBytes(p.Data)
			if err == nil {
				c := s.service.Get(a)
				err = s.Reply(conn, p.ID, protocol.OpGetResult, protocol.Uint32ToBytes(c))
			}
		case protocol.OpIncr:
			a, err := service.NewCmdArgFromBytes(p.Data)
			if err == nil {
				c := s.service.Incr(a)
				err = s.Reply(conn, p.ID, protocol.OpIncrResult, protocol.Uint32ToBytes(c))
			}
		default:
			s.Log(prefix, "Unknown op type %+v", p)
		}
		if err != nil {
			s.Log(prefix, "Unexpected error: %s", err)
		}
	}
}

func (s *Server) Reply(conn net.Conn, id uint32, op protocol.OpType, data []byte) error {
	return protocol.PackToWriter(conn, protocol.CurrentVersion, id, op, data)
}

func (s *Server) Log(prefix string, format string, v ...interface{}) {
	if s.enableLog {
		if len(prefix) > 0 {
			format = prefix + " " + format
		}
		log.Printf(format, v...)
	}
}

func (s *Server) GetConnectionCount() int {
	return syncMapLen(s.connMap)
}

func syncMapLen(m *sync.Map) int {
	length := 0
	m.Range(func(_, _ interface{}) bool {
		length++
		return true
	})
	return length
}
