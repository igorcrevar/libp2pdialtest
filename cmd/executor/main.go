package main

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/igorcrevar/libp2pdialproblem/library"
)

const (
	executable = "libp2pdialtest"
	port1      = 10000
	port2      = 10001
)

type dialParams struct {
	port          int
	privateKey    string
	dialAfter     time.Duration
	timeToLive    time.Duration
	dialAddress   string
	closeInbound  bool
	closeOutbound bool
}

func executeCommand(cmd *exec.Cmd) (string, error) {
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	err := cmd.Run()
	if err != nil {
		return "", err
	}
	if errb.Len() > 0 {
		return "", errors.New(errb.String())
	}
	return outb.String(), nil
}

func getAddress(port int, privateKey string) (string, error) {
	cmd := exec.Command(executable,
		"--port="+strconv.Itoa(port), "--pk="+privateKey, "--print-address")
	return executeCommand(cmd)
}

func executeServer(address string, params dialParams) {
	cmd := exec.Command(executable,
		"--port="+strconv.Itoa(params.port),
		"--pk="+params.privateKey,
		"--dial="+params.dialAddress,
		"--dial-after="+params.dialAfter.String(),
		"--time-to-live="+params.timeToLive.String(),
		"--close-inbound="+strconv.FormatBool(params.closeInbound),
		"--close-outbound="+strconv.FormatBool(params.closeOutbound))
	output, err := executeCommand(cmd)
	peerId := address[strings.LastIndex(address, "/")+1:]
	if err == nil {
		fmt.Printf("\nNode = %s output:\n%s", peerId, string(output))
	} else {
		fmt.Printf("\nNode = %s output:\nerror = %v\n", peerId, err)
	}
}

func main() {
	privateKey1, err := library.GenerateLibP2PKey()
	if err != nil {
		fmt.Println(err)
		return
	}
	privateKey2, err := library.GenerateLibP2PKey()
	if err != nil {
		fmt.Println(err)
		return
	}

	maAddress1, err := getAddress(port1, privateKey1)
	if err != nil {
		fmt.Println(err)
		return
	}
	maAddress2, err := getAddress(port2, privateKey2)
	if err != nil {
		fmt.Println(err)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		executeServer(maAddress1, dialParams{
			port:          port1,
			privateKey:    privateKey1,
			dialAddress:   maAddress2,
			timeToLive:    time.Second * 10,
			dialAfter:     time.Millisecond * 1000,
			closeInbound:  false,
			closeOutbound: true,
		})
		wg.Done()
	}()
	go func() {
		executeServer(maAddress2, dialParams{
			port:          port2,
			privateKey:    privateKey2,
			dialAddress:   maAddress1,
			timeToLive:    time.Second * 10,
			dialAfter:     time.Millisecond * 1000,
			closeInbound:  false,
			closeOutbound: false,
		})
		wg.Done()
	}()

	wg.Wait()
}
