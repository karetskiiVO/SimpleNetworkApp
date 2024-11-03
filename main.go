package main

import (
	"log"
	"net"

	"github.com/jessevdk/go-flags"
	netapp "github.com/karetskiiVO/SimpleNetworkApp/netapplication"
)

func main() {
	var options struct {
		Args struct {
			Mode string
			Addr string
			Port uint16
		} `positional-args:"yes" required:"yes"`
	}

	_, err := flags.NewParser(&options, flags.Default).Parse()
	if err != nil {
		log.Fatal(err)
	}
	ip := net.ParseIP(options.Args.Addr)
	if ip == nil {
		log.Fatalf("%v - incorrect ip address", options.Args.Addr)
	}

	var app netapp.Application
	switch options.Args.Mode {
	case "tcpserver":
		app, err = netapp.NewTCPserver(ip, options.Args.Port)
		if err != nil {
			log.Fatal(err)
		}
	case "tcpclient":
		app, err = netapp.NewTCPclient(ip, options.Args.Port)
		if err != nil {
			log.Fatal(err)
		}
	case "udpserver":
		app, err = netapp.NewUDPserver(ip, options.Args.Port)
		if err != nil {
			log.Fatal(err)
		}
	case "udpclient":
		app, err = netapp.NewUDPclient(ip, options.Args.Port)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer app.Close()
	app.Run()
}
