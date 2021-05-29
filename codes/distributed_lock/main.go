package main

import (
	"example.com/distributed_lock/etcd"
	"flag"
	"log"
	"strings"
)

var (
	addr = flag.String("addr", "192.168.64.4:30453", "etcd addresses")
	name = flag.String("name", "", "give a name")
)

func main() {
	flag.Parse()
	endpoints := strings.Split(*addr, ",")
	cli, err := etcd.NewClient(endpoints)
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	etcd.UseLock(cli, *name)
}
