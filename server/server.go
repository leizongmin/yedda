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
	if options.EnableLog {
		s.log("[main]", "server listen on %s", listener.Addr())
	}
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
		if s.options.EnableLog {
			s.log("[main]", "server closed")
		}
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
	s.log(prefix, "Connected. There are currently %d connections", syncMapLen(s.connMap))
	for {
		var err error
		p, err := protocol.NewPackageFromReader(conn)
		if err != nil {
			if err == io.EOF {
				s.log(prefix, "Connection closed")
			} else {
				s.log(prefix, "Read error: %s", err)
			}
			break
		}
		if p.Version != protocol.CurrentVersion {
			s.log(prefix, "Ignore protocol version %d", p.Version)
			continue
		}
		switch p.Op {
		case protocol.OpPing:
			err = s.reply(conn, p.ID, protocol.OpPong, p.Data)
		case protocol.OpPong:
			// do nothing
		case protocol.OpGet:
			a, err := service.NewCmdArgFromBytes(p.Data)
			if err == nil {
				c := s.service.Get(a)
				err = s.reply(conn, p.ID, protocol.OpGetResult, protocol.Uint32ToBytes(c))
			}
		case protocol.OpIncr:
			a, err := service.NewCmdArgFromBytes(p.Data)
			if err == nil {
				c := s.service.Incr(a)
				err = s.reply(conn, p.ID, protocol.OpIncrResult, protocol.Uint32ToBytes(c))
			}
		default:
			s.log(prefix, "Unknown op type %+v", p)
		}
		if err != nil {
			s.log(prefix, "Unexpected error: %s", err)
		}
	}
}

func (s *Server) reply(conn net.Conn, id uint32, op protocol.OpType, data []byte) error {
	return protocol.PackToWriter(conn, protocol.CurrentVersion, id, op, data)
}

func (s *Server) log(prefix string, format string, v ...interface{}) {
	if s.enableLog {
		if len(prefix) > 0 {
			format = prefix + " " + format
		}
		log.Printf(format, v...)
	}
}

func syncMapLen(m *sync.Map) int {
	length := 0
	m.Range(func(_, _ interface{}) bool {
		length++
		return true
	})
	return length
}
