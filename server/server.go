package server

import (
	"github.com/leizongmin/simple-limiter-service/protocol"
	"github.com/leizongmin/simple-limiter-service/service"
	"io"
	"log"
	"net"
	"os"
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
		options: options,
	}
	if options.EnableLog {
		log.Printf("server listen on %s", listener.Addr())
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
			log.Println("server closed")
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
	logger := log.New(os.Stdout, "["+addr.String()+"] ", log.Ltime)
	enableLog := s.options.EnableLog
	if enableLog {
		logger.Println("connected")
	}
	for {
		var err error
		p, err := protocol.NewPackageFromReader(conn)
		if err != nil {
			if enableLog {
				if err == io.EOF {
					logger.Println("connection closed")
				} else {
					logger.Printf("read error: %s", err)
				}
			}
			break
		}
		if p.Version != protocol.CurrentVersion {
			if enableLog {
				logger.Printf("ignore protocol version %d", p.Version)
			}
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
			if enableLog {
				logger.Printf("unknown op type %+v", p)
			}
		}
		if enableLog {
			if err != nil {
				logger.Printf("unexpected error: %s", err)
			}
		}
	}
}

func (s *Server) reply(conn net.Conn, id uint32, op protocol.OpType, data []byte) error {
	return protocol.PackToWriter(conn, protocol.CurrentVersion, id, op, data)
}
