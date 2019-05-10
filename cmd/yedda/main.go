package main

import (
	"flag"
	"log"
	"time"

	"github.com/leizongmin/yedda/server"
)

func main() {

	options := server.Options{}
	flag.StringVar(&options.Network, "listen-type", server.DefaultListenNetwork, "listen type, 'tcp' OR 'unix'")
	flag.StringVar(&options.Address, "listen", server.DefaultListenAddress, "listen address")
	flag.BoolVar(&options.EnableLog, "log", true, "enable log output")
	var dbSize, timeAccuracy uint64
	flag.Uint64Var(&dbSize, "size", 256, "how many database")
	flag.Uint64Var(&timeAccuracy, "accuracy", 100, "time accuracy (ms)")
	flag.Parse()

	options.DatabaseSize = uint32(dbSize)
	options.TimeAccuracy = time.Duration(timeAccuracy) * time.Millisecond

	s, err := server.NewServer(options)
	if err != nil {
		log.Fatalln(err)
	}
	if options.EnableLog {
		go func() {
			lastC := s.GetConnectionCount()
			for {
				c := s.GetConnectionCount()
				if c != lastC {
					s.Log("[main]", "There are currently %d connections", c)
				}
				lastC = c
				time.Sleep(time.Second)
			}
		}()
	}
	err = s.Loop()
	if err != nil {
		log.Fatalln(err)
	}
}
