package netapplication

import (
	"log"
	"net"
	"syscall"
)

// UDPserver implements  server
type UDPserver struct {
	fd int
}

// NewUDPserver constructs UDPserver
func NewUDPserver(ip net.IP, port uint16) (server *UDPserver, err error) {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
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

	server = &UDPserver{fd: fd}
	return
}

// Close tcp server session
func (udp UDPserver) Close() error {
	return syscall.Close(udp.fd)
}

// Run implement standart application loop
func (udp UDPserver) Run() error {
	buffer := make([]byte, batchSize)
	for {
		size, addr, err := syscall.Recvfrom(udp.fd, buffer, 0)
		if err != nil {
			return err
		}

		ipAdrr := addr.(*syscall.SockaddrInet4).Addr
		log.Printf("Recieved from %d:%d:%d:%d: %v", ipAdrr[0], ipAdrr[1], ipAdrr[2], ipAdrr[3], string(buffer[:size]))

		
		err = syscall.Sendto(udp.fd, buffer[:size], 0, addr)
		if err != nil {
			return err
		}
	}
}
