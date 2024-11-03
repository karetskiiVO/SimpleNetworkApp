package netapplication

import (
	"log"
	"net"
	"syscall"
	"time"
)

const (
	ttl time.Duration = time.Minute
)

// TCPserver implements tcp server
type TCPserver struct {
	fd int
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

// Run implement standart application loop
func (tcp TCPserver) Run() error {
	for {
		conn, addr, err := syscall.Accept(tcp.fd)
		if err != nil {
			continue
		}

		go func() {
			ipAdrr := addr.(*syscall.SockaddrInet4).Addr
			var err error
			
			for {
				var msg []byte
				msg, err = tcpRecieveMsg(conn)
				if err != nil {
					break
				}

				log.Printf("Recieved from %d:%d:%d:%d: %v", ipAdrr[0], ipAdrr[1], ipAdrr[2], ipAdrr[3], string(msg))

				err = tcpSendMsg(conn, msg)
				if err != nil {
					break
				}
			}

			log.Printf("Connection with d:%d:%d:%d lost with %v", ipAdrr[0], ipAdrr[1], ipAdrr[2], ipAdrr[3], err)
			syscall.Close(conn)
		}()
	}
}
