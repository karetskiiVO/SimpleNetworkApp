package netapplication

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"syscall"
	"time"
)

// UDPclient implements udp
type UDPclient struct {
	fd   int
	dest syscall.Sockaddr
}

// NewUDPclient constructs UDPclient
func NewUDPclient(ip net.IP, port uint16) (server *UDPclient, err error) {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err != nil {
		return
	}

	server = &UDPclient{
		fd: fd,
		dest: &syscall.SockaddrInet4{
			Addr: [4]byte(ip.To4()),
			Port: int(htons(port)),
		},
	}

	return
}

// Close гвз client session
func (udp UDPclient) Close() error {
	return syscall.Close(udp.fd)
}

// Run implement standart application loop
func (udp UDPclient) Run() error {
	input := make(chan string)
	reply := make(chan string)

	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			msg, _ := reader.ReadString('\n')
			go func() {
				input <- msg
			}()
		}
	}()

	go func() {
		buffer := make([]byte, batchSize)
		for {
			size, addr, err := syscall.Recvfrom(udp.fd, buffer, 0)
			if err != nil {
				return
			}

			ipAdrr := addr.(*syscall.SockaddrInet4).Addr
			go func() {
				reply <- fmt.Sprintf("Recieved from server %d:%d:%d:%d: %v", ipAdrr[0], ipAdrr[1], ipAdrr[2], ipAdrr[3], string(buffer[:size]))
			}()
		}
	}()

	fmt.Print("You: ")
	for msg := range input {
		err := syscall.Sendto(udp.fd, []byte(msg), 0, udp.dest)
		if err != nil {
			return err
		}

		time.Sleep(20 * time.Millisecond)
	loop:
		for {
			select {
			case msg := <-reply:
				log.Print(msg)
			default:
				break loop
			}
		}
		fmt.Print("You: ")
	}
	return nil
}
