package server

import (
	"github.com/leizongmin/simple-limiter-service/protocol"
	"github.com/leizongmin/simple-limiter-service/service"
	"io"
	"log"
	"net"
	"time"
)

const DefaultListenNetwork = "tcp"
const DefaultListenAddress = "127.0.0.1:16789"

type Server struct {
	listener net.Listener
	closed   bool
	service  *service.Service
	options  Options
}

type Options struct {
	Network      string
	Address      string
	TimeAccuracy time.Duration
	DatabaseSize uint32
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
		options: options,
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
	log.Printf("Accept new connection: %s", addr)
	for {
		var err error
		p, err := protocol.NewPackageFromReader(conn)
		if err != nil {
			if err == io.EOF {
				log.Printf("Connection %s closed", addr)
			} else {
				log.Printf("Read from connection %s error: %s", addr, err)
			}
			break
		}
		if p.Version != protocol.CurrentVersion {
			log.Printf("Ignore protocol version %d from %s", p.Version, addr)
			continue
		}
		switch p.Op {
		case protocol.OpPing:
			err = protocol.PackToWriter(conn, protocol.CurrentVersion, protocol.OpPong, p.Data)
		case protocol.OpPong:
			// do nothing
		case protocol.OpGet:
			a, err := service.NewCmdArgFromBytes(p.Data)
			if err == nil {
				c := s.service.Get(a)
				err = protocol.PackToWriter(conn, protocol.CurrentVersion, protocol.OpGetResult, protocol.Uint32ToBytes(c))
			}
		case protocol.OpIncr:
			a, err := service.NewCmdArgFromBytes(p.Data)
			if err == nil {
				c := s.service.Incr(a)
				err = protocol.PackToWriter(conn, protocol.CurrentVersion, protocol.OpIncrResult, protocol.Uint32ToBytes(c))
			}
		default:
			log.Printf("Unknown OpType %+v from connection %s", p, addr)
		}
		if err != nil {
			log.Printf("Unexpected error from connection %s: %s", addr, err)
		}
	}
}
