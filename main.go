package main

import (
	"flag"
	"fmt"
	kitlog "github.com/go-kit/kit/log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var level, listen, backend string
	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flags.StringVar(&listen, "listen", "", "Listen for connections on this address.")
	flags.StringVar(&backend, "backend", "", "The address of the backend to forward to.")
	flags.StringVar(&level, "level", "info", "The logging level.")
	flags.Parse(os.Args[1:])

	w := kitlog.NewSyncWriter(os.Stderr)
	logger := kitlog.NewLogfmtLogger(w)

	if listen == "" || backend == "" {
		fmt.Fprintln(os.Stderr, "listen and backend options required")
		os.Exit(1)
	}

	p := Proxy{Logger: logger, Listen: listen, Backend: backend}

	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		if err := p.Close(); err != nil {
			logger.Log("err", err)
		}
	}()

	if err := p.Run(); err != nil {
		logger.Log("err", err)
	}
}
