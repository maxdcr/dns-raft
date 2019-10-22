package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/maxdrc/dns-raft/dns"
	"github.com/maxdrc/dns-raft/store"
)

var (
	dnsaddr  string
	raftaddr string
	raftjoin string
	raftid   string
	zonefile string
)

func init() {
	flag.StringVar(&dnsaddr, "dns.addr", ":5350", "DNS listen address")
	flag.StringVar(&raftaddr, "raft.addr", ":15370", "Raft bus transport bind address")
	flag.StringVar(&raftjoin, "raft.join", "", "Join to already exist cluster")
	flag.StringVar(&raftid, "id", "", "node id")
	flag.StringVar(&zonefile, "zone.file", "", "Zone file containing resource records")
}

func main() {
	flag.Parse()

	quitCh := make(chan int)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	kvs := store.InitStore(raftaddr, raftjoin, raftid)
	dns := dns.NewDNS(kvs, dnsaddr)
	dns.Start()
	dns.LoadZone(zonefile)
	go func() {
		for {
			s := <-sigCh
			switch s {
			case syscall.SIGHUP:
				dns.LoadZone(zonefile)
			case syscall.SIGINT:
				fmt.Println("leaving")
				kvs.Stop()
				fmt.Println("DONE")
				quitCh <- 0
			default:
				fmt.Println("shutdown")
				quitCh <- 0
			}
		}
	}()
	code := <-quitCh
	os.Exit(code)
}
