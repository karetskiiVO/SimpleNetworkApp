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

// Close tcp server session
func (tcp TCPclient) Close() error {
	return syscall.Close(tcp.fd)
}

// Run implement standart application loop
func (tcp TCPclient) Run() error {
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
		for {
			msg, err := tcpRecieveMsg(tcp.fd)
			if err != nil {
				return
			}
			go func() {
				reply <- fmt.Sprintf("Recieved from server: %v", string(msg))
			}()
		}
	}()

	fmt.Print("You: ")
	for msg := range input {
		err := tcpSendMsg(tcp.fd, []byte(msg))
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
