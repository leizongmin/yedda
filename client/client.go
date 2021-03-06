package client

import (
	"github.com/leizongmin/yedda/protocol"
	"github.com/leizongmin/yedda/service"
	"log"
	"net"
	"strings"
	"time"
)

const DefaultDialNetwork = "tcp"
const DefaultDialAddress = "127.0.0.1:16789"

type Client struct {
	options          Options
	closed           bool
	conn             net.Conn
	currentID        uint32
	resultGet        chan uint32
	resultIncr       chan uint32
	pingMilliseconds uint64
	resultPing       chan uint64
}

type Options struct {
	Network string
	Address string
	Db      uint32
}

func NewClient(options Options) (*Client, error) {
	if len(options.Network) < 1 {
		options.Network = DefaultDialNetwork
	}
	if len(options.Address) < 1 {
		options.Address = DefaultDialAddress
	}
	conn, err := net.Dial(options.Network, options.Address)
	if err != nil {
		return nil, err
	}
	c := Client{
		options:          options,
		closed:           false,
		conn:             conn,
		resultGet:        make(chan uint32),
		resultIncr:       make(chan uint32),
		pingMilliseconds: 0,
		resultPing:       make(chan uint64),
	}
	go c.loop()
	return &c, nil
}

func (c *Client) Close() {
	c.closed = true
	c.conn.Close()
}

func (c *Client) loop() {
	conn := c.conn
	addr := conn.RemoteAddr()
	for !c.closed {
		var err error
		p, err := protocol.NewPackageFromReader(conn)
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				c.closed = true
				break
			}
			log.Printf("Read from remote %s error: %s", addr, err)
			break
		}
		switch p.Op {
		case protocol.OpPing:
			err = c.send(protocol.OpPong, p.Data)
		case protocol.OpPong:
			t := protocol.BytesToUint64(p.Data)
			c.pingMilliseconds = getMillisecondsTimestamp() - t
			c.resultPing <- c.pingMilliseconds
		case protocol.OpGetResult:
			c.resultGet <- protocol.BytesToUint32(p.Data)
		case protocol.OpIncrResult:
			c.resultIncr <- protocol.BytesToUint32(p.Data)
		default:
			log.Printf("Unknown OpType %+v from remote %s", p, addr)
		}
		if err != nil {
			log.Printf("Unexpected error from %s: %s", addr, err)
		}
	}
}

func (c *Client) send(op protocol.OpType, data []byte) error {
	c.currentID++
	return protocol.PackToWriter(c.conn, protocol.CurrentVersion, c.currentID, op, data)
}

func (c *Client) Ping() (r uint64, err error) {
	err = c.send(protocol.OpPing, protocol.Uint64ToBytes(getMillisecondsTimestamp()))
	if err != nil {
		return r, err
	}
	r = <-c.resultPing
	return r, err
}

func (c *Client) Get(ns string, key string, milliseconds uint32) (r uint32, err error) {
	a := service.NewCmdArg(c.options.Db, ns, milliseconds, []byte(key), 0)
	b, err := a.Bytes()
	if err != nil {
		return r, err
	}
	err = c.send(protocol.OpGet, b)
	if err != nil {
		return r, err
	}
	r = <-c.resultGet
	return r, err
}

func (c *Client) Incr(ns string, key string, milliseconds uint32) (r uint32, err error) {
	return c.IncrN(ns, key, milliseconds, 1)
}

func (c *Client) IncrN(ns string, key string, milliseconds uint32, count uint32) (r uint32, err error) {
	a := service.NewCmdArg(c.options.Db, ns, milliseconds, []byte(key), count)
	b, err := a.Bytes()
	if err != nil {
		return r, err
	}
	err = c.send(protocol.OpIncr, b)
	if err != nil {
		return r, err
	}
	r = <-c.resultIncr
	return r, err
}

func getMillisecondsTimestamp() uint64 {
	return uint64(time.Now().UnixNano() / int64(time.Millisecond))
}
