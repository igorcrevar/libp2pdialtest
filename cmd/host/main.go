package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/igorcrevar/libp2pdialproblem/library"
)

type Configuration struct {
	port             int
	dialAfter        time.Duration
	timeToLive       time.Duration
	closeConnTime    time.Duration
	printAddress     bool
	dialMultiaddress string
	privateKeyString string
	closeInbound     bool
	closeOutbound    bool
}

func main() {
	var config Configuration
	flag.IntVar(&config.port, "port", 0, "port")
	flag.StringVar(&config.privateKeyString, "pk", "", "private key")
	flag.DurationVar(&config.dialAfter, "dial-after", 0, "dial after")
	flag.DurationVar(&config.timeToLive, "time-to-live", 0, "time to live")
	flag.DurationVar(&config.closeConnTime, "close-conn-time", 0, "close connection sleep time")
	flag.BoolVar(&config.printAddress, "print-address", false, "print address")
	flag.BoolVar(&config.closeInbound, "close-inbound", false, "close inbound connection")
	flag.BoolVar(&config.closeOutbound, "close-outbound", false, "close outbound connection")
	flag.StringVar(&config.dialMultiaddress, "dial", "", "dial address")

	flag.Parse()
	if config.port == 0 {
		fmt.Println("port is not defined")
		return
	}
	if config.privateKeyString == "" {
		fmt.Println("private key is not specified")
		return
	}
	if config.dialAfter <= 0 && !config.printAddress {
		fmt.Println("dial after not specified")
		return
	}
	if config.timeToLive <= 0 && !config.printAddress {
		fmt.Println("time to live not specified")
		return
	}

	pk, err := library.DecodeLibP2PKey(config.privateKeyString)
	if err != nil {
		fmt.Println(err)
		return
	}

	node, err := library.NewNode(library.NodeConfig{
		PrivateKey:    pk,
		Port:          config.port,
		CloseInbound:  config.closeInbound,
		CloseOutbound: config.closeOutbound,
		CloseConnTime: config.closeConnTime,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	defer node.Stop()

	if config.printAddress {
		fmt.Print(node.Address())
	} else {
		library.Log("Started")
		time.Sleep(config.dialAfter)
		library.Log("Dialing %s", config.dialMultiaddress)
		err = node.Dial(config.dialMultiaddress)
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(config.timeToLive)
		library.Log("Finished")
	}
}
