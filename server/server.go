package server

import "net"

const DefaultListenNetwork = "tcp"
const DefaultListenAddress = "127.0.0.1:16789"

type Server struct {
	listener net.Listener
	closed   bool
}

type Options struct {
	Network string
	Address string
}

func NewServer(options Options) (s *Server, err error) {
	listener, err := net.Listen(options.Network, options.Address)
	if err != nil {
		return s, err
	}
	if len(options.Address) < 1 {
		options.Address = DefaultListenAddress
	}
	if len(options.Network) < 1 {
		options.Network = DefaultListenNetwork
	}
	s = &Server{
		listener: listener,
		closed:   false,
	}
	return s, err
}

func (s *Server) Addr() net.Addr {
	return s.listener.Addr()
}

func (s *Server) Close() error {
	return s.listener.Close()
}

func (s *Server) Loop() error {
	for !s.closed {
		conn, err := s.listener.Accept()
		if err != nil {
			return err
		}
		s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	//conn.
}
