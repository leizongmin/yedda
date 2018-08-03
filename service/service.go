package service

import (
	"github.com/leizongmin/simple-limiter-service/core"
	"time"
)

const DefaultTickerDuration = 100 * time.Millisecond
const DefaultDatabaseSize = 16

type Service struct {
	database *core.Database
	ticker   *time.Ticker
}

type Options struct {
	DatabaseSize   uint32
	TickerDuration time.Duration
}

func NewService(options Options) *Service {
	db := uint32(DefaultDatabaseSize)
	if options.DatabaseSize > 0 {
		db = options.DatabaseSize
	}
	td := DefaultTickerDuration
	if options.TickerDuration > 0 {
		td = options.TickerDuration
	}
	return &Service{
		database: core.NewDataBase(db),
		ticker:   time.NewTicker(td),
	}
}

func (s *Service) Start() {
	go func() {
		for {
			t := <-s.ticker.C
			//fmt.Printf("ticker %s\n", t)
			s.database.DeleteExpired(t)
		}
	}()
}

func (s *Service) Stop() {
	s.ticker.Stop()
}

func (s *Service) Destroy() {
	s.ticker.Stop()
	s.database.Destroy()
}

func (s *Service) CmdIncr(db uint32, ns string, milliseconds uint32, key []byte, count uint32) uint32 {
	return s.database.Get(db).Get(ns, time.Duration(milliseconds)*time.Millisecond).Incr(key, count)
}

func (s *Service) CmdGet(db uint32, ns string, milliseconds uint32, key []byte) uint32 {
	return s.database.Get(db).Get(ns, time.Duration(milliseconds)*time.Millisecond).Get(key)
}
