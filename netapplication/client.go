package netapplication

import (
	"fmt"
	"net"
	"syscall"
	"time"
)

// TCPclient implements tcp client
type TCPclient struct {
	fd int
}

// NewTCPclient constructs TCPclient
func NewTCPclient(ip net.IP, port uint16) (client *TCPclient, err error) {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return
	}

	err = syscall.Connect(fd, &syscall.SockaddrInet4{
		Port: int(htons(port)),
		Addr: [4]byte(ip.To4()),
	})
	if err != nil {
		return
	}

	client = &TCPclient{fd: fd}
	return
}

// ListenAndServe implement standart application loop
func (tcp TCPclient) ListenAndServe() error {
	cnt := 0
	for {
		sendMsg(tcp.fd, []byte("hello " + fmt.Sprint(cnt)))
		cnt++
		time.Sleep(time.Millisecond * 1000)
	}
}
