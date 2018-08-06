package service

import (
	"github.com/leizongmin/simple-limiter-service/core"
	"time"
)

// 默认时间精度，100ms
const DefaultTimeAccuracy = 100 * time.Millisecond

// 默认数据库数量，1
const DefaultDatabaseSize = 1

type Service struct {
	database *core.Database
	ticker   *time.Ticker
	cmdReq   chan *cmdReq
	cmdRes   chan uint32
	stopChan chan bool
	isStop   bool
}

type Options struct {
	// 数据库数量，>= 1
	DatabaseSize uint32
	// 时间精度
	TimeAccuracy time.Duration
}

// 创建新服务实例
func NewService(options Options) *Service {
	db := uint32(DefaultDatabaseSize)
	if options.DatabaseSize > 0 {
		db = options.DatabaseSize
	}
	td := DefaultTimeAccuracy
	if options.TimeAccuracy > 0 {
		td = options.TimeAccuracy
	}
	return &Service{
		database: core.NewDataBase(db),
		ticker:   time.NewTicker(td),
		cmdReq:   make(chan *cmdReq),
		cmdRes:   make(chan uint32),
		stopChan: make(chan bool),
		isStop:   true,
	}
}

// 开始服务
func (s *Service) Start() {
	s.isStop = false
	go func() {
		for {
			select {
			case <-s.stopChan:
				break
			case t := <-s.ticker.C:
				//fmt.Println("ticker", t)
				s.database.DeleteExpired(t)
			case req := <-s.cmdReq:
				if req != nil {
					//fmt.Println("req", req)
					a := req.Arg
					switch req.Cmd {
					case cmdGet:
						s.cmdRes <- s.database.Get(a.Db).Get(a.Ns, time.Duration(a.Milliseconds)*time.Millisecond).Incr(a.Key, a.Count)
					case cmdIncr:
						s.cmdRes <- s.database.Get(a.Db).Get(a.Ns, time.Duration(a.Milliseconds)*time.Millisecond).Incr(a.Key, a.Count)
					}
				}
			}
		}
	}()
}

// 停止服务
func (s *Service) Stop() {
	s.isStop = true
	s.stopChan <- true
	s.ticker.Stop()
}

// 销毁服务
func (s *Service) Destroy() {
	s.Stop()
	close(s.cmdReq)
	s.database.Destroy()
}

type cmdReq struct {
	Cmd cmdType
	Arg *CmdArg
}

type cmdType uint8

const (
	_ cmdType = iota
	cmdIncr
	cmdGet
)

// 执行 INCR 命令
func (s *Service) Incr(a *CmdArg) uint32 {
	if s.isStop {
		return 1
	}
	s.cmdReq <- &cmdReq{Cmd: cmdIncr, Arg: a}
	return <-s.cmdRes
}

// 执行 GET 命令
func (s *Service) Get(a *CmdArg) uint32 {
	if s.isStop {
		return 1
	}
	s.cmdReq <- &cmdReq{Cmd: cmdGet, Arg: a}
	return <-s.cmdRes
}
