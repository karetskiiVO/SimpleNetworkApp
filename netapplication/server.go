package netapplication

import (
	"log"
	"net"
	"syscall"
)

const (
	batchSize = 1024
)

// Application interface
type Application interface {
	Close() error
	ListenAndServe() error
}

// TCPserver implements tcp server
type TCPserver struct {
	fd int
}

func htons(i uint16) uint16 {
	return (i<<8)&0xff00 | i>>8
}

// NewTCPserver constructs TCPserver
func NewTCPserver(ip net.IP, port uint16) (server *TCPserver, err error) {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return
	}

	err = syscall.Bind(fd, &syscall.SockaddrInet4{
		Port: int(htons(port)),
		Addr: [4]byte(ip.To4()),
	})
	if err != nil {
		return
	}

	err = syscall.Listen(fd, 1)
	if err != nil {
		return
	}

	server = &TCPserver{fd: fd}
	return
}

// Close tcp server session
func (tcp TCPserver) Close() (err error) {
	err = syscall.Close(tcp.fd)
	return
}

// ListenAndServe implement standart application loop
func (tcp TCPserver) ListenAndServe() error {
	for {
		conn, addr, err := syscall.Accept(tcp.fd)
		if err != nil {
			return err
		}

		msg, err := recieveFullMsg(conn)
		if err != nil {
			return err
		}

		log.Printf("Recieved from %v:%v\n", addr, msg)
	}
}
