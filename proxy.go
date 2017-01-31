package main

import (
	kitlog "github.com/go-kit/kit/log"
	"net"
	"sync"
)

// Proxy connections from Listen to Backend.
type Proxy struct {
	Logger   kitlog.Logger
	Listen   string
	Backend  string
	listener net.Listener
}

func (p *Proxy) Run() error {
	var err error
	if p.listener, err = net.Listen("tcp", p.Listen); err != nil {
		return err
	}

	wg := &sync.WaitGroup{}
	for {
		if conn, err := p.listener.Accept(); err == nil {
			wg.Add(1)
			go func() {
				defer wg.Done()
				p.handle(conn)
			}()
		} else {
			return nil
		}
	}
	wg.Wait()
	return nil
}

func (p *Proxy) Close() error {
	return p.listener.Close()
}

func (p *Proxy) handle(upConn net.Conn) {
	defer upConn.Close()
	p.Logger.Log("connection", upConn.RemoteAddr())
	downConn, err := net.Dial("tcp", p.Backend)
	if err != nil {
		p.Logger.Log("msg", "unable to connect", "backend", p.Backend, "err", err)
		return
	}
	defer downConn.Close()
	if err := Pipe(upConn, downConn); err != nil {
		p.Logger.Log("msg", "pipe failed", "err", err)
	} else {
		p.Logger.Log("msg", "disconnected", "client", upConn.RemoteAddr())
	}
}
