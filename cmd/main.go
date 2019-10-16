package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/maxdcr/pagesjaunes/server"
	"github.com/maxdcr/pagesjaunes/store"
)

var (
	tcpaddr  string
	dnsaddr  string
	raftaddr string
	raftjoin string
	raftdir  string
	nodeID   string
	zonefile string
)

func init() {
	flag.StringVar(&tcpaddr, "tcp.addr", ":5370", "TCP listen address")
	flag.StringVar(&dnsaddr, "dns.addr", ":5354", "DNS listen address")
	flag.StringVar(&raftaddr, "raft.addr", ":15370", "Raft bus transport bind address")
	flag.StringVar(&raftjoin, "raft.join", "", "Join to already exist cluster")
	flag.StringVar(&raftdir, "raft.dir", "./raft", "Raft data directory")
	flag.StringVar(&nodeID, "id", "", "node id")
	flag.StringVar(&zonefile, "zone.file", "", "Zone file containing resource records")
}

func main() {
	flag.Parse()

	kvs := store.InitStore(raftdir, raftaddr, raftjoin, nodeID)
	server.InitTCP(kvs, tcpaddr)
	server.InitDNS(kvs, dnsaddr, zonefile)

	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, os.Kill, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-quitCh
}
